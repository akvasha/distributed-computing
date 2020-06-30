// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/validate.proto

package pbauth

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

type ValidateRequest struct {
	AccessToken          string   `protobuf:"bytes,1,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidateRequest) Reset()         { *m = ValidateRequest{} }
func (m *ValidateRequest) String() string { return proto.CompactTextString(m) }
func (*ValidateRequest) ProtoMessage()    {}
func (*ValidateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4f819e4945ce700c, []int{0}
}

func (m *ValidateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateRequest.Unmarshal(m, b)
}
func (m *ValidateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateRequest.Marshal(b, m, deterministic)
}
func (m *ValidateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateRequest.Merge(m, src)
}
func (m *ValidateRequest) XXX_Size() int {
	return xxx_messageInfo_ValidateRequest.Size(m)
}
func (m *ValidateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateRequest proto.InternalMessageInfo

func (m *ValidateRequest) GetAccessToken() string {
	if m != nil {
		return m.AccessToken
	}
	return ""
}

type ValidateResponse struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Admin                bool     `protobuf:"varint,2,opt,name=admin,proto3" json:"admin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ValidateResponse) Reset()         { *m = ValidateResponse{} }
func (m *ValidateResponse) String() string { return proto.CompactTextString(m) }
func (*ValidateResponse) ProtoMessage()    {}
func (*ValidateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4f819e4945ce700c, []int{1}
}

func (m *ValidateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ValidateResponse.Unmarshal(m, b)
}
func (m *ValidateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ValidateResponse.Marshal(b, m, deterministic)
}
func (m *ValidateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ValidateResponse.Merge(m, src)
}
func (m *ValidateResponse) XXX_Size() int {
	return xxx_messageInfo_ValidateResponse.Size(m)
}
func (m *ValidateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ValidateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ValidateResponse proto.InternalMessageInfo

func (m *ValidateResponse) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *ValidateResponse) GetAdmin() bool {
	if m != nil {
		return m.Admin
	}
	return false
}

func init() {
	proto.RegisterType((*ValidateRequest)(nil), "pbauth.ValidateRequest")
	proto.RegisterType((*ValidateResponse)(nil), "pbauth.ValidateResponse")
}

func init() {
	proto.RegisterFile("proto/validate.proto", fileDescriptor_4f819e4945ce700c)
}

var fileDescriptor_4f819e4945ce700c = []byte{
	// 171 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x2f, 0x4b, 0xcc, 0xc9, 0x4c, 0x49, 0x2c, 0x49, 0xd5, 0x03, 0x73, 0x85, 0xd8, 0x0a,
	0x92, 0x12, 0x4b, 0x4b, 0x32, 0x94, 0x8c, 0xb9, 0xf8, 0xc3, 0xa0, 0x32, 0x41, 0xa9, 0x85, 0xa5,
	0xa9, 0xc5, 0x25, 0x42, 0x0a, 0x5c, 0xdc, 0x89, 0xc9, 0xc9, 0xa9, 0xc5, 0xc5, 0x21, 0xf9, 0xd9,
	0xa9, 0x79, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0xc8, 0x42, 0x4a, 0x2e, 0x5c, 0x02, 0x08,
	0x4d, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x52, 0x5c, 0x1c, 0xa5, 0xc5, 0xa9, 0x45, 0x79,
	0x89, 0xb9, 0xa9, 0x50, 0x2d, 0x70, 0xbe, 0x90, 0x08, 0x17, 0x6b, 0x62, 0x4a, 0x6e, 0x66, 0x9e,
	0x04, 0x93, 0x02, 0xa3, 0x06, 0x47, 0x10, 0x84, 0x63, 0xe4, 0xc1, 0xc5, 0xee, 0x58, 0x5a, 0x92,
	0x11, 0x54, 0x90, 0x2c, 0x64, 0xcb, 0xc5, 0x01, 0x33, 0x50, 0x48, 0x5c, 0x0f, 0xe2, 0x34, 0x3d,
	0x34, 0x77, 0x49, 0x49, 0x60, 0x4a, 0x40, 0xec, 0x4e, 0x62, 0x03, 0xfb, 0xc9, 0x18, 0x10, 0x00,
	0x00, 0xff, 0xff, 0xdb, 0xeb, 0x7d, 0xf0, 0xeb, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthRpcClient is the client API for AuthRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthRpcClient interface {
	Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error)
}

type authRpcClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthRpcClient(cc grpc.ClientConnInterface) AuthRpcClient {
	return &authRpcClient{cc}
}

func (c *authRpcClient) Validate(ctx context.Context, in *ValidateRequest, opts ...grpc.CallOption) (*ValidateResponse, error) {
	out := new(ValidateResponse)
	err := c.cc.Invoke(ctx, "/pbauth.AuthRpc/Validate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthRpcServer is the server API for AuthRpc service.
type AuthRpcServer interface {
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
}

// UnimplementedAuthRpcServer can be embedded to have forward compatible implementations.
type UnimplementedAuthRpcServer struct {
}

func (*UnimplementedAuthRpcServer) Validate(ctx context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Validate not implemented")
}

func RegisterAuthRpcServer(s *grpc.Server, srv AuthRpcServer) {
	s.RegisterService(&_AuthRpc_serviceDesc, srv)
}

func _AuthRpc_Validate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthRpcServer).Validate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pbauth.AuthRpc/Validate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthRpcServer).Validate(ctx, req.(*ValidateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthRpc_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pbauth.AuthRpc",
	HandlerType: (*AuthRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Validate",
			Handler:    _AuthRpc_Validate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/validate.proto",
}
