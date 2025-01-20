// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/portal/v1/tx.proto

package portal

import (
	context "context"
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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

// MsgDeliver is the entrypoint for delivering Noble Dollar Portal messages.
// This will primarily be used by validators in a vote extension, however is
// left public to enable permissionless manual relaying.
type MsgDeliver struct {
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Vaa    []byte `protobuf:"bytes,2,opt,name=vaa,proto3" json:"vaa,omitempty"`
}

func (m *MsgDeliver) Reset()         { *m = MsgDeliver{} }
func (m *MsgDeliver) String() string { return proto.CompactTextString(m) }
func (*MsgDeliver) ProtoMessage()    {}
func (*MsgDeliver) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5414e5ec63723f0, []int{0}
}
func (m *MsgDeliver) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeliver) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeliver.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeliver) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeliver.Merge(m, src)
}
func (m *MsgDeliver) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeliver) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeliver.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeliver proto.InternalMessageInfo

// MsgDeliverResponse is the response of the Deliver message.
type MsgDeliverResponse struct {
}

func (m *MsgDeliverResponse) Reset()         { *m = MsgDeliverResponse{} }
func (m *MsgDeliverResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDeliverResponse) ProtoMessage()    {}
func (*MsgDeliverResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5414e5ec63723f0, []int{1}
}
func (m *MsgDeliverResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeliverResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeliverResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeliverResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeliverResponse.Merge(m, src)
}
func (m *MsgDeliverResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeliverResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeliverResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeliverResponse proto.InternalMessageInfo

// MsgSetPeer allows the Noble Dollar Portal owner to set external peers.
type MsgSetPeer struct {
	Signer      string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Chain       uint16 `protobuf:"varint,2,opt,name=chain,proto3,customtype=uint16" json:"chain"`
	Transceiver []byte `protobuf:"bytes,3,opt,name=transceiver,proto3" json:"transceiver,omitempty"`
	Manager     []byte `protobuf:"bytes,4,opt,name=manager,proto3" json:"manager,omitempty"`
}

func (m *MsgSetPeer) Reset()         { *m = MsgSetPeer{} }
func (m *MsgSetPeer) String() string { return proto.CompactTextString(m) }
func (*MsgSetPeer) ProtoMessage()    {}
func (*MsgSetPeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5414e5ec63723f0, []int{2}
}
func (m *MsgSetPeer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetPeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetPeer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetPeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetPeer.Merge(m, src)
}
func (m *MsgSetPeer) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetPeer) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetPeer.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetPeer proto.InternalMessageInfo

// MsgSetPeerResponse is the response of the SetPeer message.
type MsgSetPeerResponse struct {
}

func (m *MsgSetPeerResponse) Reset()         { *m = MsgSetPeerResponse{} }
func (m *MsgSetPeerResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetPeerResponse) ProtoMessage()    {}
func (*MsgSetPeerResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5414e5ec63723f0, []int{3}
}
func (m *MsgSetPeerResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetPeerResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetPeerResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetPeerResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetPeerResponse.Merge(m, src)
}
func (m *MsgSetPeerResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetPeerResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetPeerResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetPeerResponse proto.InternalMessageInfo

// MsgTransfer is the entrypoint for transferring the Noble Dollar cross-chain.
type MsgTransfer struct {
	Signer    string                `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Chain     uint16                `protobuf:"varint,2,opt,name=chain,proto3,customtype=uint16" json:"chain"`
	Recipient []byte                `protobuf:"bytes,3,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount    cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
}

func (m *MsgTransfer) Reset()         { *m = MsgTransfer{} }
func (m *MsgTransfer) String() string { return proto.CompactTextString(m) }
func (*MsgTransfer) ProtoMessage()    {}
func (*MsgTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5414e5ec63723f0, []int{4}
}
func (m *MsgTransfer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgTransfer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgTransfer.Merge(m, src)
}
func (m *MsgTransfer) XXX_Size() int {
	return m.Size()
}
func (m *MsgTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_MsgTransfer proto.InternalMessageInfo

// MsgTransferResponse is the response of the Transfer message.
type MsgTransferResponse struct {
}

func (m *MsgTransferResponse) Reset()         { *m = MsgTransferResponse{} }
func (m *MsgTransferResponse) String() string { return proto.CompactTextString(m) }
func (*MsgTransferResponse) ProtoMessage()    {}
func (*MsgTransferResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f5414e5ec63723f0, []int{5}
}
func (m *MsgTransferResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgTransferResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgTransferResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgTransferResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgTransferResponse.Merge(m, src)
}
func (m *MsgTransferResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgTransferResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgTransferResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgTransferResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgDeliver)(nil), "noble.dollar.portal.v1.MsgDeliver")
	proto.RegisterType((*MsgDeliverResponse)(nil), "noble.dollar.portal.v1.MsgDeliverResponse")
	proto.RegisterType((*MsgSetPeer)(nil), "noble.dollar.portal.v1.MsgSetPeer")
	proto.RegisterType((*MsgSetPeerResponse)(nil), "noble.dollar.portal.v1.MsgSetPeerResponse")
	proto.RegisterType((*MsgTransfer)(nil), "noble.dollar.portal.v1.MsgTransfer")
	proto.RegisterType((*MsgTransferResponse)(nil), "noble.dollar.portal.v1.MsgTransferResponse")
}

func init() { proto.RegisterFile("noble/dollar/portal/v1/tx.proto", fileDescriptor_f5414e5ec63723f0) }

var fileDescriptor_f5414e5ec63723f0 = []byte{
	// 536 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x53, 0x31, 0x6f, 0xd3, 0x40,
	0x14, 0xb6, 0x1b, 0x9a, 0x92, 0x2b, 0x20, 0x30, 0x49, 0x31, 0x16, 0xd8, 0x91, 0x61, 0xa8, 0x82,
	0x62, 0x37, 0x20, 0x81, 0xd4, 0x8d, 0x88, 0x81, 0x0e, 0x91, 0x90, 0xcb, 0x02, 0x42, 0x42, 0xd7,
	0xe4, 0xb8, 0x5a, 0xd8, 0x77, 0xd6, 0xdd, 0x35, 0x6a, 0x59, 0x40, 0x4c, 0x88, 0x89, 0x9f, 0x90,
	0x91, 0x31, 0x43, 0x7f, 0x44, 0x37, 0xaa, 0x4e, 0x88, 0xa1, 0x42, 0xc9, 0x10, 0xfe, 0x04, 0x12,
	0xf2, 0xdd, 0xb9, 0x69, 0x11, 0x84, 0x0a, 0xa9, 0x8b, 0xe5, 0xf7, 0xde, 0x77, 0xdf, 0x7b, 0xdf,
	0x77, 0xef, 0x80, 0x47, 0xe8, 0x46, 0x82, 0xc2, 0x1e, 0x4d, 0x12, 0xc8, 0xc2, 0x8c, 0x32, 0x01,
	0x93, 0xb0, 0xdf, 0x0a, 0xc5, 0x76, 0x90, 0x31, 0x2a, 0xa8, 0xb5, 0x24, 0x01, 0x81, 0x02, 0x04,
	0x0a, 0x10, 0xf4, 0x5b, 0xce, 0x15, 0x98, 0xc6, 0x84, 0x86, 0xf2, 0xab, 0xa0, 0xce, 0xb5, 0x2e,
	0xe5, 0x29, 0xe5, 0x61, 0xca, 0x71, 0x4e, 0x91, 0x72, 0xac, 0x0b, 0xd7, 0x55, 0xe1, 0xa5, 0x8c,
	0x42, 0x15, 0xe8, 0x52, 0x15, 0x53, 0x4c, 0x55, 0x3e, 0xff, 0x53, 0x59, 0xff, 0x2d, 0x00, 0x1d,
	0x8e, 0x1f, 0xa1, 0x24, 0xee, 0x23, 0x66, 0xad, 0x80, 0x32, 0x8f, 0x31, 0x41, 0xcc, 0x36, 0xeb,
	0xe6, 0x72, 0xa5, 0x6d, 0x1f, 0xec, 0x36, 0xab, 0x9a, 0xe5, 0x61, 0xaf, 0xc7, 0x10, 0xe7, 0xeb,
	0x82, 0xc5, 0x04, 0x47, 0x1a, 0x67, 0x5d, 0x06, 0xa5, 0x3e, 0x84, 0xf6, 0x5c, 0xdd, 0x5c, 0xbe,
	0x10, 0xe5, 0xbf, 0xab, 0xcd, 0x0f, 0x03, 0xcf, 0xf8, 0x31, 0xf0, 0x8c, 0xf7, 0x93, 0x61, 0x43,
	0xc3, 0x3e, 0x4e, 0x86, 0x8d, 0xda, 0x49, 0xe5, 0xba, 0xa5, 0x5f, 0x05, 0xd6, 0x74, 0x80, 0x08,
	0xf1, 0x8c, 0x12, 0x8e, 0xfc, 0x2f, 0xa6, 0x9c, 0x6b, 0x1d, 0x89, 0x27, 0xe8, 0xbf, 0xe6, 0xba,
	0x0d, 0xe6, 0xbb, 0x9b, 0x30, 0x26, 0x72, 0xb2, 0x8b, 0xed, 0x4b, 0x7b, 0x87, 0x9e, 0xf1, 0xed,
	0xd0, 0x2b, 0x6f, 0xc5, 0x44, 0xb4, 0xee, 0x47, 0xaa, 0x68, 0xd5, 0xc1, 0xa2, 0x60, 0x90, 0xf0,
	0x2e, 0xca, 0xbb, 0xdb, 0x25, 0xa9, 0xe2, 0x78, 0xca, 0xb2, 0xc1, 0x42, 0x0a, 0x09, 0xc4, 0x88,
	0xd9, 0xe7, 0x64, 0xb5, 0x08, 0x4f, 0xab, 0x53, 0x4b, 0xd0, 0x3a, 0x75, 0x74, 0xa4, 0xf3, 0xa7,
	0x09, 0x16, 0x3b, 0x1c, 0x3f, 0xcd, 0x3b, 0xbe, 0x3a, 0x43, 0xa1, 0x37, 0x40, 0x85, 0xa1, 0x6e,
	0x9c, 0xc5, 0x88, 0x08, 0x2d, 0x73, 0x9a, 0xb0, 0x1e, 0x83, 0x32, 0x4c, 0xe9, 0x16, 0x11, 0x52,
	0x63, 0xa5, 0xbd, 0xa2, 0x49, 0x6a, 0xaa, 0x33, 0xef, 0xbd, 0x0e, 0x62, 0x1a, 0xa6, 0x50, 0x6c,
	0x06, 0x6b, 0x44, 0x1c, 0xec, 0x36, 0x81, 0x1e, 0x69, 0x8d, 0x88, 0xcf, 0x93, 0x61, 0xc3, 0x8c,
	0xf4, 0xf9, 0xd5, 0xe0, 0x2f, 0xa6, 0x2c, 0x9d, 0x34, 0xa5, 0xd0, 0xeb, 0xd7, 0xc0, 0xd5, 0x63,
	0xf2, 0x0b, 0x5b, 0xee, 0x0e, 0xe6, 0x40, 0xa9, 0xc3, 0xb1, 0xf5, 0x0c, 0x2c, 0x14, 0xab, 0xe9,
	0x07, 0x7f, 0x7e, 0x1e, 0xc1, 0x74, 0x7b, 0x9c, 0xc6, 0xbf, 0x31, 0x45, 0x8b, 0x9c, 0xba, 0xd8,
	0xae, 0x59, 0xd4, 0x1a, 0x33, 0x93, 0xfa, 0xb7, 0x4b, 0xb5, 0x5e, 0x80, 0xf3, 0x47, 0x17, 0x7a,
	0x6b, 0xc6, 0xb9, 0x02, 0xe4, 0xdc, 0x39, 0x05, 0xa8, 0x60, 0x77, 0xe6, 0xdf, 0xe5, 0x8e, 0xb7,
	0x1f, 0xec, 0x8d, 0x5c, 0x73, 0x7f, 0xe4, 0x9a, 0xdf, 0x47, 0xae, 0xf9, 0x69, 0xec, 0x1a, 0xfb,
	0x63, 0xd7, 0xf8, 0x3a, 0x76, 0x8d, 0xe7, 0x37, 0x35, 0x8d, 0xe2, 0xdc, 0xde, 0x79, 0x13, 0x8a,
	0x9d, 0x0c, 0x71, 0xed, 0xfd, 0x46, 0x59, 0x3e, 0xfc, 0x7b, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff,
	0xae, 0xdb, 0x9e, 0xf6, 0x90, 0x04, 0x00, 0x00,
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
	Deliver(ctx context.Context, in *MsgDeliver, opts ...grpc.CallOption) (*MsgDeliverResponse, error)
	SetPeer(ctx context.Context, in *MsgSetPeer, opts ...grpc.CallOption) (*MsgSetPeerResponse, error)
	Transfer(ctx context.Context, in *MsgTransfer, opts ...grpc.CallOption) (*MsgTransferResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) Deliver(ctx context.Context, in *MsgDeliver, opts ...grpc.CallOption) (*MsgDeliverResponse, error) {
	out := new(MsgDeliverResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.portal.v1.Msg/Deliver", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetPeer(ctx context.Context, in *MsgSetPeer, opts ...grpc.CallOption) (*MsgSetPeerResponse, error) {
	out := new(MsgSetPeerResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.portal.v1.Msg/SetPeer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Transfer(ctx context.Context, in *MsgTransfer, opts ...grpc.CallOption) (*MsgTransferResponse, error) {
	out := new(MsgTransferResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.portal.v1.Msg/Transfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	Deliver(context.Context, *MsgDeliver) (*MsgDeliverResponse, error)
	SetPeer(context.Context, *MsgSetPeer) (*MsgSetPeerResponse, error)
	Transfer(context.Context, *MsgTransfer) (*MsgTransferResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) Deliver(ctx context.Context, req *MsgDeliver) (*MsgDeliverResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deliver not implemented")
}
func (*UnimplementedMsgServer) SetPeer(ctx context.Context, req *MsgSetPeer) (*MsgSetPeerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPeer not implemented")
}
func (*UnimplementedMsgServer) Transfer(ctx context.Context, req *MsgTransfer) (*MsgTransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Transfer not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_Deliver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeliver)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Deliver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.portal.v1.Msg/Deliver",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Deliver(ctx, req.(*MsgDeliver))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetPeer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPeer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPeer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.portal.v1.Msg/SetPeer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPeer(ctx, req.(*MsgSetPeer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgTransfer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.portal.v1.Msg/Transfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Transfer(ctx, req.(*MsgTransfer))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_serviceDesc = _Msg_serviceDesc
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "noble.dollar.portal.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deliver",
			Handler:    _Msg_Deliver_Handler,
		},
		{
			MethodName: "SetPeer",
			Handler:    _Msg_SetPeer_Handler,
		},
		{
			MethodName: "Transfer",
			Handler:    _Msg_Transfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/dollar/portal/v1/tx.proto",
}

func (m *MsgDeliver) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeliver) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeliver) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Vaa) > 0 {
		i -= len(m.Vaa)
		copy(dAtA[i:], m.Vaa)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Vaa)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDeliverResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeliverResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeliverResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgSetPeer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetPeer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetPeer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Manager) > 0 {
		i -= len(m.Manager)
		copy(dAtA[i:], m.Manager)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Manager)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Transceiver) > 0 {
		i -= len(m.Transceiver)
		copy(dAtA[i:], m.Transceiver)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Transceiver)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Chain != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Chain))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSetPeerResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetPeerResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetPeerResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgTransfer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgTransfer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgTransfer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Recipient) > 0 {
		i -= len(m.Recipient)
		copy(dAtA[i:], m.Recipient)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Recipient)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Chain != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Chain))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgTransferResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgTransferResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgTransferResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgDeliver) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Vaa)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgDeliverResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgSetPeer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Chain != 0 {
		n += 1 + sovTx(uint64(m.Chain))
	}
	l = len(m.Transceiver)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Manager)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSetPeerResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgTransfer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Chain != 0 {
		n += 1 + sovTx(uint64(m.Chain))
	}
	l = len(m.Recipient)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgTransferResponse) Size() (n int) {
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
func (m *MsgDeliver) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDeliver: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeliver: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
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
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vaa", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Vaa = append(m.Vaa[:0], dAtA[iNdEx:postIndex]...)
			if m.Vaa == nil {
				m.Vaa = []byte{}
			}
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
func (m *MsgDeliverResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDeliverResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeliverResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgSetPeer) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetPeer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetPeer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
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
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			m.Chain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Chain |= uint16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transceiver", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Transceiver = append(m.Transceiver[:0], dAtA[iNdEx:postIndex]...)
			if m.Transceiver == nil {
				m.Transceiver = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Manager", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Manager = append(m.Manager[:0], dAtA[iNdEx:postIndex]...)
			if m.Manager == nil {
				m.Manager = []byte{}
			}
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
func (m *MsgSetPeerResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetPeerResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetPeerResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgTransfer) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgTransfer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgTransfer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
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
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			m.Chain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Chain |= uint16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Recipient", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Recipient = append(m.Recipient[:0], dAtA[iNdEx:postIndex]...)
			if m.Recipient == nil {
				m.Recipient = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgTransferResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgTransferResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgTransferResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
