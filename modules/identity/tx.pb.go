// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: identity/tx.proto

package identity

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgCreateIdentity defines a message to create an identity
type MsgCreateIdentity struct {
	Id          string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PubKey      *PubKeyInfo `protobuf:"bytes,2,opt,name=pub_key,json=pubKey,proto3" json:"pubkey" yaml:"pubkey"`
	Certificate string      `protobuf:"bytes,3,opt,name=certificate,proto3" json:"certificate,omitempty"`
	Credentials string      `protobuf:"bytes,4,opt,name=credentials,proto3" json:"credentials,omitempty"`
	Owner       string      `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (m *MsgCreateIdentity) Reset()         { *m = MsgCreateIdentity{} }
func (m *MsgCreateIdentity) String() string { return proto.CompactTextString(m) }
func (*MsgCreateIdentity) ProtoMessage()    {}
func (*MsgCreateIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a49ec0beed01e79, []int{0}
}
func (m *MsgCreateIdentity) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateIdentity.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateIdentity.Merge(m, src)
}
func (m *MsgCreateIdentity) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateIdentity proto.InternalMessageInfo

// MsgCreateIdentityResponse defines the Msg/Create response type.
type MsgCreateIdentityResponse struct {
}

func (m *MsgCreateIdentityResponse) Reset()         { *m = MsgCreateIdentityResponse{} }
func (m *MsgCreateIdentityResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreateIdentityResponse) ProtoMessage()    {}
func (*MsgCreateIdentityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a49ec0beed01e79, []int{1}
}
func (m *MsgCreateIdentityResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateIdentityResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateIdentityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateIdentityResponse.Merge(m, src)
}
func (m *MsgCreateIdentityResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateIdentityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateIdentityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateIdentityResponse proto.InternalMessageInfo

// MsgUpdateIdentity defines a message to update an identity
type MsgUpdateIdentity struct {
	Id          string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PubKey      *PubKeyInfo `protobuf:"bytes,2,opt,name=pub_key,json=pubKey,proto3" json:"pubkey" yaml:"pubkey"`
	Certificate string      `protobuf:"bytes,3,opt,name=certificate,proto3" json:"certificate,omitempty"`
	Credentials string      `protobuf:"bytes,4,opt,name=credentials,proto3" json:"credentials,omitempty"`
	Owner       string      `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (m *MsgUpdateIdentity) Reset()         { *m = MsgUpdateIdentity{} }
func (m *MsgUpdateIdentity) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateIdentity) ProtoMessage()    {}
func (*MsgUpdateIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a49ec0beed01e79, []int{2}
}
func (m *MsgUpdateIdentity) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateIdentity.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateIdentity.Merge(m, src)
}
func (m *MsgUpdateIdentity) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateIdentity proto.InternalMessageInfo

// MsgUpdateIdentityResponse defines the Msg/Update response type.
type MsgUpdateIdentityResponse struct {
}

func (m *MsgUpdateIdentityResponse) Reset()         { *m = MsgUpdateIdentityResponse{} }
func (m *MsgUpdateIdentityResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUpdateIdentityResponse) ProtoMessage()    {}
func (*MsgUpdateIdentityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a49ec0beed01e79, []int{3}
}
func (m *MsgUpdateIdentityResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUpdateIdentityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUpdateIdentityResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUpdateIdentityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUpdateIdentityResponse.Merge(m, src)
}
func (m *MsgUpdateIdentityResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUpdateIdentityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUpdateIdentityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUpdateIdentityResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgCreateIdentity)(nil), "iritamod.identity.MsgCreateIdentity")
	proto.RegisterType((*MsgCreateIdentityResponse)(nil), "iritamod.identity.MsgCreateIdentityResponse")
	proto.RegisterType((*MsgUpdateIdentity)(nil), "iritamod.identity.MsgUpdateIdentity")
	proto.RegisterType((*MsgUpdateIdentityResponse)(nil), "iritamod.identity.MsgUpdateIdentityResponse")
}

func init() { proto.RegisterFile("identity/tx.proto", fileDescriptor_4a49ec0beed01e79) }

var fileDescriptor_4a49ec0beed01e79 = []byte{
	// 378 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x93, 0xcb, 0x4e, 0xf2, 0x40,
	0x14, 0xc7, 0x3b, 0xdc, 0xbe, 0x7c, 0x43, 0x24, 0xa1, 0x21, 0xb1, 0x42, 0x2c, 0xa4, 0x71, 0xc1,
	0x42, 0xda, 0x04, 0x77, 0x2c, 0x71, 0x45, 0x08, 0x09, 0xc1, 0xb8, 0x71, 0x63, 0x5a, 0xe6, 0x50,
	0x47, 0x68, 0xa7, 0x69, 0xa7, 0xd1, 0xbe, 0x85, 0x8f, 0xe0, 0xe3, 0xb0, 0xc4, 0x1d, 0x2b, 0xa2,
	0xb0, 0x31, 0x2e, 0x7d, 0x02, 0x43, 0x4b, 0x51, 0x6e, 0x89, 0x6b, 0x77, 0xe7, 0xf2, 0x9f, 0xf3,
	0xcf, 0xef, 0x4c, 0x0e, 0xce, 0x53, 0x02, 0x36, 0xa7, 0x3c, 0xd0, 0xf8, 0xa3, 0xea, 0xb8, 0x8c,
	0x33, 0x31, 0x4f, 0x5d, 0xca, 0x75, 0x8b, 0x11, 0x35, 0xee, 0x15, 0x8f, 0xd7, 0xaa, 0x38, 0x88,
	0xb4, 0xc5, 0x82, 0xc9, 0x4c, 0x16, 0x86, 0xda, 0x32, 0x8a, 0xaa, 0xca, 0x14, 0xe1, 0x7c, 0xc7,
	0x33, 0x2f, 0x5d, 0xd0, 0x39, 0xb4, 0x56, 0x2f, 0xc4, 0x1c, 0x4e, 0x50, 0x22, 0xa1, 0x0a, 0xaa,
	0xfe, 0xef, 0x25, 0x28, 0x11, 0xaf, 0xf0, 0x3f, 0xc7, 0x37, 0x6e, 0x87, 0x10, 0x48, 0x89, 0x0a,
	0xaa, 0x66, 0xeb, 0xa7, 0xea, 0x8e, 0xb3, 0xda, 0xf5, 0x8d, 0x36, 0x04, 0x2d, 0x7b, 0xc0, 0x9a,
	0xa5, 0x8f, 0x59, 0x39, 0xe3, 0xf8, 0xc6, 0x10, 0x82, 0xcf, 0x59, 0xf9, 0x28, 0xd0, 0xad, 0x51,
	0x43, 0x89, 0x72, 0xa5, 0xb7, 0x6c, 0xb4, 0x21, 0x10, 0x2b, 0x38, 0xdb, 0x07, 0x97, 0xd3, 0x01,
	0xed, 0xeb, 0x1c, 0xa4, 0x64, 0xe8, 0xf6, 0xb3, 0x14, 0x2a, 0x5c, 0x08, 0xe7, 0xeb, 0x23, 0x4f,
	0x4a, 0xad, 0x14, 0xdf, 0x25, 0xb1, 0x80, 0xd3, 0xec, 0xc1, 0x06, 0x57, 0x4a, 0x87, 0xbd, 0x28,
	0x69, 0xa4, 0xde, 0x9f, 0xcb, 0x48, 0x29, 0xe1, 0x93, 0x1d, 0xb2, 0x1e, 0x78, 0x0e, 0xb3, 0x3d,
	0x88, 0xb9, 0xaf, 0x1d, 0xf2, 0x47, 0xb9, 0x37, 0xc9, 0x62, 0xee, 0xfa, 0x0b, 0xc2, 0xc9, 0x8e,
	0x67, 0x8a, 0x04, 0xe7, 0xb6, 0xfe, 0xfc, 0x6c, 0x0f, 0xda, 0xce, 0xfe, 0x8a, 0xe7, 0xbf, 0x51,
	0xc5, 0x6e, 0x4b, 0x97, 0xad, 0x0d, 0x1f, 0x70, 0xd9, 0x54, 0x1d, 0x72, 0xd9, 0xcf, 0xd4, 0xec,
	0x8e, 0xdf, 0x64, 0x61, 0x3c, 0x97, 0xd1, 0x64, 0x2e, 0xa3, 0xd7, 0xb9, 0x8c, 0x9e, 0x16, 0xb2,
	0x30, 0x59, 0xc8, 0xc2, 0x74, 0x21, 0x0b, 0x37, 0x75, 0x93, 0xf2, 0x3b, 0xdf, 0x50, 0xfb, 0xcc,
	0xd2, 0x0c, 0xaa, 0xdb, 0xf7, 0x14, 0x74, 0xaa, 0x85, 0xf3, 0x6b, 0x1e, 0x19, 0xd6, 0x4c, 0xa6,
	0x59, 0x8c, 0xf8, 0x23, 0xf0, 0xd6, 0x17, 0x63, 0x64, 0xc2, 0xe3, 0xb8, 0xf8, 0x0a, 0x00, 0x00,
	0xff, 0xff, 0x69, 0x6b, 0xec, 0x77, 0x73, 0x03, 0x00, 0x00,
}

func (this *MsgCreateIdentity) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MsgCreateIdentity)
	if !ok {
		that2, ok := that.(MsgCreateIdentity)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	if !this.PubKey.Equal(that1.PubKey) {
		return false
	}
	if this.Certificate != that1.Certificate {
		return false
	}
	if this.Credentials != that1.Credentials {
		return false
	}
	if this.Owner != that1.Owner {
		return false
	}
	return true
}
func (this *MsgUpdateIdentity) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MsgUpdateIdentity)
	if !ok {
		that2, ok := that.(MsgUpdateIdentity)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Id != that1.Id {
		return false
	}
	if !this.PubKey.Equal(that1.PubKey) {
		return false
	}
	if this.Certificate != that1.Certificate {
		return false
	}
	if this.Credentials != that1.Credentials {
		return false
	}
	if this.Owner != that1.Owner {
		return false
	}
	return true
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// CreateIdentity defines a method for creating a new identity.
	CreateIdentity(ctx context.Context, in *MsgCreateIdentity, opts ...grpc.CallOption) (*MsgCreateIdentityResponse, error)
	// UpdateIdentity defines a method for Updating a identity.
	UpdateIdentity(ctx context.Context, in *MsgUpdateIdentity, opts ...grpc.CallOption) (*MsgUpdateIdentityResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateIdentity(ctx context.Context, in *MsgCreateIdentity, opts ...grpc.CallOption) (*MsgCreateIdentityResponse, error) {
	out := new(MsgCreateIdentityResponse)
	err := c.cc.Invoke(ctx, "/iritamod.identity.Msg/CreateIdentity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateIdentity(ctx context.Context, in *MsgUpdateIdentity, opts ...grpc.CallOption) (*MsgUpdateIdentityResponse, error) {
	out := new(MsgUpdateIdentityResponse)
	err := c.cc.Invoke(ctx, "/iritamod.identity.Msg/UpdateIdentity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// CreateIdentity defines a method for creating a new identity.
	CreateIdentity(context.Context, *MsgCreateIdentity) (*MsgCreateIdentityResponse, error)
	// UpdateIdentity defines a method for Updating a identity.
	UpdateIdentity(context.Context, *MsgUpdateIdentity) (*MsgUpdateIdentityResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) CreateIdentity(ctx context.Context, req *MsgCreateIdentity) (*MsgCreateIdentityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIdentity not implemented")
}
func (*UnimplementedMsgServer) UpdateIdentity(ctx context.Context, req *MsgUpdateIdentity) (*MsgUpdateIdentityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateIdentity not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iritamod.identity.Msg/CreateIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateIdentity(ctx, req.(*MsgCreateIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateIdentity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateIdentity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iritamod.identity.Msg/UpdateIdentity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateIdentity(ctx, req.(*MsgUpdateIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "iritamod.identity.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIdentity",
			Handler:    _Msg_CreateIdentity_Handler,
		},
		{
			MethodName: "UpdateIdentity",
			Handler:    _Msg_UpdateIdentity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "identity/tx.proto",
}

func (m *MsgCreateIdentity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateIdentity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Credentials) > 0 {
		i -= len(m.Credentials)
		copy(dAtA[i:], m.Credentials)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Credentials)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Certificate) > 0 {
		i -= len(m.Certificate)
		copy(dAtA[i:], m.Certificate)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Certificate)))
		i--
		dAtA[i] = 0x1a
	}
	if m.PubKey != nil {
		{
			size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateIdentityResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateIdentityResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgUpdateIdentity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateIdentity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateIdentity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Credentials) > 0 {
		i -= len(m.Credentials)
		copy(dAtA[i:], m.Credentials)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Credentials)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Certificate) > 0 {
		i -= len(m.Certificate)
		copy(dAtA[i:], m.Certificate)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Certificate)))
		i--
		dAtA[i] = 0x1a
	}
	if m.PubKey != nil {
		{
			size, err := m.PubKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUpdateIdentityResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUpdateIdentityResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUpdateIdentityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgCreateIdentity) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.PubKey != nil {
		l = m.PubKey.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Certificate)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Credentials)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgCreateIdentityResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgUpdateIdentity) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.PubKey != nil {
		l = m.PubKey.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Certificate)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Credentials)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgUpdateIdentityResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgCreateIdentity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreateIdentity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateIdentity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PubKey == nil {
				m.PubKey = &PubKeyInfo{}
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Certificate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Certificate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Credentials", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Credentials = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgCreateIdentityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgCreateIdentityResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateIdentityResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateIdentity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpdateIdentity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateIdentity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PubKey == nil {
				m.PubKey = &PubKeyInfo{}
			}
			if err := m.PubKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Certificate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Certificate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Credentials", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Credentials = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgUpdateIdentityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgUpdateIdentityResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUpdateIdentityResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTx
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
