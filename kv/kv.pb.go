// Code generated by protoc-gen-go. DO NOT EDIT.
// source: kv.proto

/*
Package kv is a generated protocol buffer package.

It is generated from these files:
	kv.proto

It has these top-level messages:
	Data
	Result
	Binlog
*/
package kv

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ERRNO int32

const (
	ERRNO_SUCCESS ERRNO = 0
	ERRNO_VERSION ERRNO = 1
	ERRNO_UNKNOW  ERRNO = -1
)

var ERRNO_name = map[int32]string{
	0:  "SUCCESS",
	1:  "VERSION",
	-1: "UNKNOW",
}
var ERRNO_value = map[string]int32{
	"SUCCESS": 0,
	"VERSION": 1,
	"UNKNOW":  -1,
}

func (x ERRNO) String() string {
	return proto.EnumName(ERRNO_name, int32(x))
}
func (ERRNO) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The request message containing the user's name.
type Data struct {
	Key     []byte `protobuf:"bytes,2,opt,name=Key,proto3" json:"Key,omitempty"`
	Value   []byte `protobuf:"bytes,3,opt,name=Value,proto3" json:"Value,omitempty"`
	Version int64  `protobuf:"varint,4,opt,name=Version" json:"Version,omitempty"`
}

func (m *Data) Reset()                    { *m = Data{} }
func (m *Data) String() string            { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()               {}
func (*Data) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Data) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Data) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Data) GetVersion() int64 {
	if m != nil {
		return m.Version
	}
	return 0
}

// The response message containing the greetings
type Result struct {
	Errno ERRNO `protobuf:"varint,1,opt,name=Errno,enum=kv.ERRNO" json:"Errno,omitempty"`
	Data  *Data `protobuf:"bytes,2,opt,name=Data" json:"Data,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Result) GetErrno() ERRNO {
	if m != nil {
		return m.Errno
	}
	return ERRNO_SUCCESS
}

func (m *Result) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

type Binlog struct {
	// if Sequence == 0 is copy mode
	// on copy,reset status
	Sequence int64 `protobuf:"varint,1,opt,name=Sequence" json:"Sequence,omitempty"`
	Data     *Data `protobuf:"bytes,2,opt,name=Data" json:"Data,omitempty"`
}

func (m *Binlog) Reset()                    { *m = Binlog{} }
func (m *Binlog) String() string            { return proto.CompactTextString(m) }
func (*Binlog) ProtoMessage()               {}
func (*Binlog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Binlog) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Binlog) GetData() *Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Data)(nil), "kv.Data")
	proto.RegisterType((*Result)(nil), "kv.Result")
	proto.RegisterType((*Binlog)(nil), "kv.Binlog")
	proto.RegisterEnum("kv.ERRNO", ERRNO_name, ERRNO_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Kv service

type KvClient interface {
	// return Version only
	// if Errno is VERSION_ERROR ,it will return the newer value
	Set(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error)
	Get(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error)
}

type kvClient struct {
	cc *grpc.ClientConn
}

func NewKvClient(cc *grpc.ClientConn) KvClient {
	return &kvClient{cc}
}

func (c *kvClient) Set(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/kv.kv/Set", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvClient) Get(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/kv.kv/Get", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Kv service

type KvServer interface {
	// return Version only
	// if Errno is VERSION_ERROR ,it will return the newer value
	Set(context.Context, *Data) (*Result, error)
	Get(context.Context, *Data) (*Result, error)
}

func RegisterKvServer(s *grpc.Server, srv KvServer) {
	s.RegisterService(&_Kv_serviceDesc, srv)
}

func _Kv_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.kv/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServer).Set(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kv_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kv.kv/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServer).Get(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

var _Kv_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kv.kv",
	HandlerType: (*KvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _Kv_Set_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Kv_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kv.proto",
}

// Client API for Replica service

type ReplicaClient interface {
	// request send last binlog
	Sync(ctx context.Context, in *Binlog, opts ...grpc.CallOption) (Replica_SyncClient, error)
}

type replicaClient struct {
	cc *grpc.ClientConn
}

func NewReplicaClient(cc *grpc.ClientConn) ReplicaClient {
	return &replicaClient{cc}
}

func (c *replicaClient) Sync(ctx context.Context, in *Binlog, opts ...grpc.CallOption) (Replica_SyncClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Replica_serviceDesc.Streams[0], c.cc, "/kv.Replica/Sync", opts...)
	if err != nil {
		return nil, err
	}
	x := &replicaSyncClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Replica_SyncClient interface {
	Recv() (*Binlog, error)
	grpc.ClientStream
}

type replicaSyncClient struct {
	grpc.ClientStream
}

func (x *replicaSyncClient) Recv() (*Binlog, error) {
	m := new(Binlog)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Replica service

type ReplicaServer interface {
	// request send last binlog
	Sync(*Binlog, Replica_SyncServer) error
}

func RegisterReplicaServer(s *grpc.Server, srv ReplicaServer) {
	s.RegisterService(&_Replica_serviceDesc, srv)
}

func _Replica_Sync_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Binlog)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ReplicaServer).Sync(m, &replicaSyncServer{stream})
}

type Replica_SyncServer interface {
	Send(*Binlog) error
	grpc.ServerStream
}

type replicaSyncServer struct {
	grpc.ServerStream
}

func (x *replicaSyncServer) Send(m *Binlog) error {
	return x.ServerStream.SendMsg(m)
}

var _Replica_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kv.Replica",
	HandlerType: (*ReplicaServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Sync",
			Handler:       _Replica_Sync_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "kv.proto",
}

func init() { proto.RegisterFile("kv.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 279 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x5f, 0x4b, 0xb4, 0x40,
	0x18, 0xc5, 0x77, 0xd6, 0xbf, 0xef, 0xb3, 0x2f, 0x21, 0x4f, 0x5d, 0x88, 0x14, 0xc9, 0x5c, 0x49,
	0x90, 0x84, 0x41, 0x1f, 0xc0, 0x4d, 0xb6, 0x58, 0x50, 0x98, 0x61, 0xed, 0xda, 0x64, 0x08, 0x51,
	0x74, 0x73, 0x55, 0xd8, 0x4f, 0x5f, 0xcc, 0xc8, 0x2e, 0xdd, 0x94, 0x57, 0xe7, 0x77, 0x8e, 0x9c,
	0xc3, 0xc3, 0x80, 0x5d, 0x4f, 0xe1, 0xbe, 0xef, 0x86, 0x0e, 0x97, 0xf5, 0x44, 0x5f, 0x40, 0x7f,
	0x2e, 0x86, 0x02, 0x1d, 0xd0, 0xb6, 0xe2, 0xe8, 0x2e, 0x7d, 0x12, 0xfc, 0x67, 0x52, 0xe2, 0x15,
	0x18, 0x79, 0xd1, 0x8c, 0xc2, 0xd5, 0x94, 0x37, 0x03, 0xba, 0x60, 0xe5, 0xa2, 0x3f, 0x54, 0x5d,
	0xeb, 0xea, 0x3e, 0x09, 0x34, 0x76, 0x42, 0xba, 0x01, 0x93, 0x89, 0xc3, 0xd8, 0x0c, 0x78, 0x0b,
	0x46, 0xd2, 0xf7, 0x6d, 0xe7, 0x12, 0x9f, 0x04, 0x17, 0xd1, 0xbf, 0xb0, 0x9e, 0xc2, 0x84, 0xb1,
	0x34, 0x63, 0xb3, 0x8f, 0xd7, 0xf3, 0xa8, 0x5a, 0x5b, 0x45, 0xb6, 0xcc, 0x25, 0x33, 0xe5, 0xd2,
	0x18, 0xcc, 0xb8, 0x6a, 0x9b, 0xee, 0x03, 0x3d, 0xb0, 0xb9, 0xf8, 0x1c, 0x45, 0x5b, 0x0a, 0xd5,
	0xa5, 0xb1, 0x33, 0xff, 0xdd, 0x71, 0xf7, 0x04, 0x86, 0x5a, 0xc4, 0x15, 0x58, 0x7c, 0xb7, 0x5e,
	0x27, 0x9c, 0x3b, 0x0b, 0x09, 0x79, 0xc2, 0xf8, 0x6b, 0x96, 0x3a, 0x04, 0x2f, 0xc1, 0xdc, 0xa5,
	0xdb, 0x34, 0x7b, 0x73, 0xbe, 0x4e, 0x1f, 0x89, 0x62, 0x58, 0xd6, 0x13, 0xde, 0x80, 0xc6, 0xc5,
	0x80, 0xe7, 0x52, 0x0f, 0xa4, 0x9a, 0xaf, 0xa3, 0x0b, 0x19, 0x6f, 0x7e, 0x8f, 0xa3, 0x7b, 0xb0,
	0x98, 0xd8, 0x37, 0x55, 0x59, 0x20, 0x05, 0x9d, 0x1f, 0xdb, 0x12, 0xd5, 0x0f, 0xf3, 0x51, 0xde,
	0x0f, 0x4d, 0x17, 0x0f, 0xe4, 0xdd, 0x54, 0x8f, 0xf1, 0xf8, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x79,
	0x91, 0x8e, 0x7e, 0x98, 0x01, 0x00, 0x00,
}
