package main

import (
	"bytes"
	"context"
	"fmt"
	"hash/crc64"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	kv "github.com/humboldt-xie/mykv/kv"
	ldb "github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
)

func I2S(seq int64) []byte {
	return []byte(fmt.Sprintf("%16x", seq))
}
func S2I(seq []byte) int64 {
	sseq := strings.Trim(string(seq), " ")
	useq, _ := strconv.ParseInt(sseq, 16, 64)
	return useq
}

func Hash(b []byte) uint64 {
	tab := crc64.MakeTable(crc64.ISO)
	crc := uint64(0)
	crc = crc64.Update(crc, tab, b)
	return crc
}

type Watchers struct {
	m        sync.Mutex
	watchers []chan int64
	closer   []chan int64
}

func (w *Watchers) Fire(v int64) {
	// put event
	for _, ch := range w.watchers {
		select {
		case ch <- v:
		default:
		}
	}

}
func (w *Watchers) Release(ch chan int64) {
	w.m.Lock()
	defer w.m.Unlock()
	index := -1
	for i, v := range w.watchers {
		if v == ch {
			index = i
			break
		}
	}
	if index < 0 {
		return
	}
	w.watchers = append(w.watchers[:index], w.watchers[index+1:]...)
	w.closer = append(w.closer, ch)
}

func (w *Watchers) New() chan int64 {
	w.m.Lock()
	defer w.m.Unlock()
	ch := make(chan int64, 1)
	w.watchers = append(w.watchers, ch)
	return ch
}

type KvServer struct {
	Listen       string
	Filename     string
	DB           *ldb.DB
	locks        []sync.Mutex
	wl           sync.Mutex
	Watcher      Watchers
	LastSequence int64
}

func (s *KvServer) Init() (err error) {
	if len(s.Filename) < 3 {
		panic("filename too short")
	}
	s.DB, err = ldb.OpenFile(s.Filename, nil)
	if err != nil {
		return err
	}
	s.locks = make([]sync.Mutex, 1024)
	return nil
}
func (s *KvServer) get(key []byte) (old_data *kv.Data, err error) {
	old_data = &kv.Data{Key: key, Value: []byte{}, Version: 0}
	ob, err := s.DB.Get(s.Key("d", key), nil)
	if err != nil && err != ldb.ErrNotFound {
		return nil, err
	}
	if err == ldb.ErrNotFound {
		return old_data, nil
	}
	err = proto.Unmarshal(ob, old_data)
	if err != nil {
		return nil, err
	}
	return old_data, nil
}

func (s *KvServer) slaveOf(address string) {

}

func (s *KvServer) Key(prefix string, key []byte) []byte {
	return append([]byte(prefix+"-"), key...)
}

func (s *KvServer) put(data *kv.Data, withBindlog bool) error {
	b, err := proto.Marshal(data)
	if err != nil {
		return err
	}
	s.wl.Lock()
	defer s.wl.Unlock()
	seq := s.LastSequence + 1
	batch := new(ldb.Batch)
	s.DB.Put(s.Key("b", I2S(seq)), b, nil)
	s.DB.Put(s.Key("d", data.Key), b, nil)
	err = s.DB.Write(batch, nil)
	if err != nil {
		return err
	}
	s.LastSequence = seq
	s.Watcher.Fire(seq)
	return nil
}

func (s *KvServer) Set(ctx context.Context, d *kv.Data) (rs *kv.Result, err error) {
	index := Hash(d.Key) % uint64(len(s.locks))
	locker := &s.locks[index]
	locker.Lock()
	defer locker.Unlock()

	old_data, err := s.get(d.Key)
	if err != nil {
		return nil, err
	}

	if d.Version != 0 && old_data.Version != 0 && d.Version != old_data.Version {
		return &kv.Result{
			Errno: kv.ERRNO_VERSION,
			Data:  old_data,
		}, nil
	}

	if bytes.Compare(d.Value, old_data.Value) == 0 {
		return &kv.Result{
			Errno: kv.ERRNO_SUCCESS,
			Data:  &kv.Data{Version: old_data.Version},
		}, nil
	}
	d.Version = old_data.Version + 1
	log.Printf("old %#v new %#v\n", old_data, d)
	err = s.put(d, true)
	if err != nil {
		return nil, err
	}
	return &kv.Result{
		Errno: kv.ERRNO_SUCCESS,
		Data:  &kv.Data{Version: d.Version},
	}, err
}

func (s *KvServer) Get(ctx context.Context, d *kv.Data) (rs *kv.Result, err error) {
	old_data, err := s.get(d.Key)
	if err != nil {
		return nil, err
	}
	return &kv.Result{
		Errno: kv.ERRNO_SUCCESS,
		Data:  old_data,
	}, nil
}
func (s *KvServer) Sync(last *kv.Binlog, rep kv.Replica_SyncServer) (err error) {
	if last.Sequence == 0 {
		//copy
		last.Sequence = s.LastSequence
		iter := s.DB.NewIterator(nil, nil)
		for ok := iter.Seek(s.Key("d", []byte{})); ok; iter.Next() {

			key := iter.Key()
			if len(key) < 2 || string(key[:2]) != "d-" {
				break
			}
			value := iter.Value()
			data := &kv.Data{}
			err := proto.Unmarshal(value, data)
			if err != nil {
				log.Debugf("data error:%s %s", string(key), err)
			}
			err = rep.Send(&kv.Binlog{Sequence: 0, Data: data})
			if err != nil {
				return err
			}
		}
		iter.Release()
	}
	ch := s.Watcher.New()
	for {
		if last.Sequence >= s.LastSequence {
			<-ch
		}
		iter := s.DB.NewIterator(nil, nil)
		for ok := iter.Seek(s.Key("b", I2S(last.Sequence))); ok; iter.Next() {
			key := iter.Key()
			if len(key) < 2 || string(key[:2]) != "b-" {
				break
			}
			seq := S2I(key[2:])
			last.Sequence = seq
			value := iter.Value()
			data := &kv.Data{}
			err := proto.Unmarshal(value, data)
			if err != nil {
				log.Debugf("data error:%s %s", string(key), err)
				continue
			}
			err = rep.Send(&kv.Binlog{Sequence: last.Sequence, Data: data})
			if err != nil {
				return err
			}
		}
	}

	//will not reach
	return nil
}

func (s *KvServer) Remove() {
	if len(s.Filename) < 3 {
		panic("file name too short")
	}
	defer os.RemoveAll(s.Filename)
}

func (s *KvServer) Serve() {
	lis, err := net.Listen("tcp", s.Listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	kv.RegisterKvServer(server, s)
	kv.RegisterReplicaServer(server, s)
	server.Serve(lis)
}

func main() {
	s := &KvServer{Listen: ":8101", Filename: "/tmp/a"}
	s.Serve()
}
