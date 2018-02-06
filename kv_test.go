package main

import (
	"fmt"
	kv "github.com/humboldt-xie/mykv/kv"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	testCase := []struct {
		Key        string
		Value      string
		Version    int64
		Errno      kv.ERRNO
		ResVersion int64
	}{
		{"a", "b", 0, kv.ERRNO_SUCCESS, 1},
		{"a", "b", 0, kv.ERRNO_SUCCESS, 1}, //no change version no change
		{"a", "c", 0, kv.ERRNO_SUCCESS, 2}, //no change version no change
	}
	file := "/tmp/" + fmt.Sprintf("%d", time.Now().Unix())
	mkv := &KvServer{Filename: file}
	mkv.Init()
	mkv.Remove()

	for index, v := range testCase {
		res, err := mkv.Set(nil, &kv.Data{Key: []byte(v.Key), Value: []byte(v.Value), Version: v.Version})
		if err != nil {
			t.Errorf("%d %v error %s", index, v, err)
			continue
		}
		if res == nil {
			t.Errorf("%d %v res is nil", index, v)
			continue
		}
		if res.Errno != v.Errno {
			t.Errorf("%d %v errno:%s", index, v, res.Errno)
			continue
		}
		if res.Data.Version != v.ResVersion {
			t.Errorf("%d %v version no match %d \n", index, v, res.Data.Version)
			continue
		}
	}
}
