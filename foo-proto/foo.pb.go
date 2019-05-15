// Code generated by protoc-gen-go. DO NOT EDIT.
// source: foo/foo.proto

package foo

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type BarRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BarRequest) Reset()         { *m = BarRequest{} }
func (m *BarRequest) String() string { return proto.CompactTextString(m) }
func (*BarRequest) ProtoMessage()    {}
func (*BarRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a61802764a7d7e5a, []int{0}
}

func (m *BarRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BarRequest.Unmarshal(m, b)
}
func (m *BarRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BarRequest.Marshal(b, m, deterministic)
}
func (m *BarRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BarRequest.Merge(m, src)
}
func (m *BarRequest) XXX_Size() int {
	return xxx_messageInfo_BarRequest.Size(m)
}
func (m *BarRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BarRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BarRequest proto.InternalMessageInfo

type BarResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BarResponse) Reset()         { *m = BarResponse{} }
func (m *BarResponse) String() string { return proto.CompactTextString(m) }
func (*BarResponse) ProtoMessage()    {}
func (*BarResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a61802764a7d7e5a, []int{1}
}

func (m *BarResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BarResponse.Unmarshal(m, b)
}
func (m *BarResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BarResponse.Marshal(b, m, deterministic)
}
func (m *BarResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BarResponse.Merge(m, src)
}
func (m *BarResponse) XXX_Size() int {
	return xxx_messageInfo_BarResponse.Size(m)
}
func (m *BarResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BarResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BarResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*BarRequest)(nil), "geometa.BarRequest")
	proto.RegisterType((*BarResponse)(nil), "geometa.BarResponse")
}

func init() { proto.RegisterFile("foo/foo.proto", fileDescriptor_a61802764a7d7e5a) }

var fileDescriptor_a61802764a7d7e5a = []byte{
	// 148 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0xcb, 0xcf, 0xd7,
	0x4f, 0xcb, 0xcf, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0x4f, 0xcd, 0xcf, 0x4d,
	0x2d, 0x49, 0x54, 0xe2, 0xe1, 0xe2, 0x72, 0x4a, 0x2c, 0x0a, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e,
	0x51, 0xe2, 0xe5, 0xe2, 0x06, 0xf3, 0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x8d, 0x2c, 0xb9, 0x98,
	0xdd, 0xf2, 0xf3, 0x85, 0x8c, 0xb8, 0x98, 0x9d, 0x12, 0x8b, 0x84, 0x84, 0xf5, 0xa0, 0x9a, 0xf4,
	0x10, 0x3a, 0xa4, 0x44, 0x50, 0x05, 0x21, 0x1a, 0x95, 0x18, 0x9c, 0xb4, 0xa3, 0x34, 0xd3, 0x33,
	0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0xcb, 0x53, 0x13, 0x4b, 0x32, 0x52, 0x8b,
	0x8a, 0xf3, 0x4b, 0x8b, 0x92, 0x53, 0xf5, 0xc1, 0x4e, 0x28, 0x4a, 0x2d, 0xc8, 0xd7, 0x4f, 0x07,
	0xbb, 0x29, 0x89, 0x0d, 0x2c, 0x62, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x0d, 0x05, 0x8c, 0xc1,
	0xa5, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FooClient is the client API for Foo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FooClient interface {
	Bar(ctx context.Context, in *BarRequest, opts ...grpc.CallOption) (*BarResponse, error)
}

type fooClient struct {
	cc *grpc.ClientConn
}

func NewFooClient(cc *grpc.ClientConn) FooClient {
	return &fooClient{cc}
}

func (c *fooClient) Bar(ctx context.Context, in *BarRequest, opts ...grpc.CallOption) (*BarResponse, error) {
	out := new(BarResponse)
	err := c.cc.Invoke(ctx, "/geometa.Foo/Bar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FooServer is the server API for Foo service.
type FooServer interface {
	Bar(context.Context, *BarRequest) (*BarResponse, error)
}

func RegisterFooServer(s *grpc.Server, srv FooServer) {
	s.RegisterService(&_Foo_serviceDesc, srv)
}

func _Foo_Bar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FooServer).Bar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/geometa.Foo/Bar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FooServer).Bar(ctx, req.(*BarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Foo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "geometa.Foo",
	HandlerType: (*FooServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bar",
			Handler:    _Foo_Bar_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "foo/foo.proto",
}