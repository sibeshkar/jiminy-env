// Code generated by protoc-gen-go. DO NOT EDIT.
// source: env.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Request struct {
	EnvId                string   `protobuf:"bytes,1,opt,name=env_id,json=envId,proto3" json:"env_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbfacb289c786e17, []int{0}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetEnvId() string {
	if m != nil {
		return m.EnvId
	}
	return ""
}

type Response struct {
	Status               string   `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbfacb289c786e17, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type Reward struct {
	Reward               float32  `protobuf:"fixed32,1,opt,name=reward,proto3" json:"reward,omitempty"`
	Done                 bool     `protobuf:"varint,2,opt,name=done,proto3" json:"done,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Reward) Reset()         { *m = Reward{} }
func (m *Reward) String() string { return proto.CompactTextString(m) }
func (*Reward) ProtoMessage()    {}
func (*Reward) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbfacb289c786e17, []int{2}
}

func (m *Reward) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Reward.Unmarshal(m, b)
}
func (m *Reward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Reward.Marshal(b, m, deterministic)
}
func (m *Reward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Reward.Merge(m, src)
}
func (m *Reward) XXX_Size() int {
	return xxx_messageInfo_Reward.Size(m)
}
func (m *Reward) XXX_DiscardUnknown() {
	xxx_messageInfo_Reward.DiscardUnknown(m)
}

var xxx_messageInfo_Reward proto.InternalMessageInfo

func (m *Reward) GetReward() float32 {
	if m != nil {
		return m.Reward
	}
	return 0
}

func (m *Reward) GetDone() bool {
	if m != nil {
		return m.Done
	}
	return false
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_bbfacb289c786e17, []int{3}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Request)(nil), "proto.Request")
	proto.RegisterType((*Response)(nil), "proto.Response")
	proto.RegisterType((*Reward)(nil), "proto.Reward")
	proto.RegisterType((*Empty)(nil), "proto.Empty")
}

func init() { proto.RegisterFile("env.proto", fileDescriptor_bbfacb289c786e17) }

var fileDescriptor_bbfacb289c786e17 = []byte{
	// 220 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x8f, 0xcf, 0x4a, 0xc4, 0x30,
	0x10, 0xc6, 0xc9, 0xba, 0xc9, 0x6e, 0x07, 0xff, 0x40, 0x40, 0x59, 0xf6, 0x54, 0x72, 0xb1, 0x5e,
	0x7a, 0x50, 0xdf, 0x40, 0x8a, 0x14, 0x3c, 0xe5, 0x05, 0xa4, 0x9a, 0x01, 0x0b, 0x3a, 0xa9, 0xcd,
	0x34, 0xe2, 0x93, 0xfa, 0x3a, 0xd2, 0x18, 0x7a, 0xb5, 0xa7, 0x99, 0x1f, 0xdf, 0xef, 0xf0, 0x7d,
	0x50, 0x20, 0xc5, 0x7a, 0x18, 0x3d, 0x7b, 0x2d, 0xd3, 0x31, 0x25, 0xec, 0x2c, 0x7e, 0x4e, 0x18,
	0x58, 0x5f, 0x82, 0x42, 0x8a, 0xcf, 0xbd, 0x3b, 0x88, 0x52, 0x54, 0x85, 0x95, 0x48, 0xb1, 0x75,
	0xc6, 0xc0, 0xde, 0x62, 0x18, 0x3c, 0x05, 0xd4, 0x57, 0xa0, 0x02, 0x77, 0x3c, 0x85, 0xac, 0x64,
	0x32, 0xf7, 0xa0, 0x2c, 0x7e, 0x75, 0xa3, 0x9b, 0x8d, 0x31, 0x7d, 0xc9, 0xd8, 0xd8, 0x4c, 0x5a,
	0xc3, 0xd6, 0x79, 0xc2, 0xc3, 0xa6, 0x14, 0xd5, 0xde, 0xa6, 0xdf, 0xec, 0x40, 0x36, 0x1f, 0x03,
	0x7f, 0xdf, 0xfe, 0x08, 0x38, 0x69, 0x28, 0xea, 0x6b, 0xd8, 0xb6, 0xd4, 0xb3, 0x3e, 0xff, 0xeb,
	0x58, 0xe7, 0x66, 0xc7, 0x8b, 0x85, 0x73, 0x8f, 0x1b, 0x50, 0x4f, 0xdd, 0x44, 0xaf, 0x6f, 0xff,
	0xab, 0x15, 0x48, 0x8b, 0x01, 0x79, 0x95, 0xf9, 0xf0, 0xee, 0x03, 0xae, 0x31, 0x8b, 0x47, 0xe4,
	0xbc, 0xf8, 0x34, 0xa7, 0x69, 0xca, 0xf1, 0x6c, 0x71, 0xe7, 0xf0, 0x45, 0x25, 0xba, 0xfb, 0x0d,
	0x00, 0x00, 0xff, 0xff, 0xae, 0x1e, 0xbd, 0xc1, 0x79, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EnvClient is the client API for Env service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EnvClient interface {
	Init(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	Launch(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	Reset(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	Close(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	GetReward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Reward, error)
}

type envClient struct {
	cc *grpc.ClientConn
}

func NewEnvClient(cc *grpc.ClientConn) EnvClient {
	return &envClient{cc}
}

func (c *envClient) Init(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Env/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envClient) Launch(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Env/Launch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envClient) Reset(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Env/Reset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envClient) Close(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.Env/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *envClient) GetReward(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Reward, error) {
	out := new(Reward)
	err := c.cc.Invoke(ctx, "/proto.Env/GetReward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnvServer is the server API for Env service.
type EnvServer interface {
	Init(context.Context, *Request) (*Response, error)
	Launch(context.Context, *Request) (*Response, error)
	Reset(context.Context, *Request) (*Response, error)
	Close(context.Context, *Request) (*Response, error)
	GetReward(context.Context, *Empty) (*Reward, error)
}

// UnimplementedEnvServer can be embedded to have forward compatible implementations.
type UnimplementedEnvServer struct {
}

func (*UnimplementedEnvServer) Init(ctx context.Context, req *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (*UnimplementedEnvServer) Launch(ctx context.Context, req *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Launch not implemented")
}
func (*UnimplementedEnvServer) Reset(ctx context.Context, req *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}
func (*UnimplementedEnvServer) Close(ctx context.Context, req *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (*UnimplementedEnvServer) GetReward(ctx context.Context, req *Empty) (*Reward, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReward not implemented")
}

func RegisterEnvServer(s *grpc.Server, srv EnvServer) {
	s.RegisterService(&_Env_serviceDesc, srv)
}

func _Env_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Env/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvServer).Init(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Env_Launch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvServer).Launch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Env/Launch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvServer).Launch(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Env_Reset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvServer).Reset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Env/Reset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvServer).Reset(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Env_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Env/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvServer).Close(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Env_GetReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnvServer).GetReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Env/GetReward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnvServer).GetReward(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _Env_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Env",
	HandlerType: (*EnvServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _Env_Init_Handler,
		},
		{
			MethodName: "Launch",
			Handler:    _Env_Launch_Handler,
		},
		{
			MethodName: "Reset",
			Handler:    _Env_Reset_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Env_Close_Handler,
		},
		{
			MethodName: "GetReward",
			Handler:    _Env_GetReward_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "env.proto",
}
