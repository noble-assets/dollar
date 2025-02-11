// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/v1/tx.proto

package types

import (
	context "context"
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

// MsgClaimYield is a message holders of the Noble Dollar can use to claim their yield.
type MsgClaimYield struct {
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgClaimYield) Reset()         { *m = MsgClaimYield{} }
func (m *MsgClaimYield) String() string { return proto.CompactTextString(m) }
func (*MsgClaimYield) ProtoMessage()    {}
func (*MsgClaimYield) Descriptor() ([]byte, []int) {
	return fileDescriptor_cda63894d37623a5, []int{0}
}
func (m *MsgClaimYield) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimYield) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimYield.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimYield) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimYield.Merge(m, src)
}
func (m *MsgClaimYield) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimYield) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimYield.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimYield proto.InternalMessageInfo

// MsgClaimYieldResponse is the response of the ClaimYield message.
type MsgClaimYieldResponse struct {
}

func (m *MsgClaimYieldResponse) Reset()         { *m = MsgClaimYieldResponse{} }
func (m *MsgClaimYieldResponse) String() string { return proto.CompactTextString(m) }
func (*MsgClaimYieldResponse) ProtoMessage()    {}
func (*MsgClaimYieldResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cda63894d37623a5, []int{1}
}
func (m *MsgClaimYieldResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimYieldResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimYieldResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimYieldResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimYieldResponse.Merge(m, src)
}
func (m *MsgClaimYieldResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimYieldResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimYieldResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimYieldResponse proto.InternalMessageInfo

// MsgSetPausedState allows the authority to configure the Noble Dollar Portal paused state.
type MsgSetPausedState struct {
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Paused bool   `protobuf:"varint,2,opt,name=paused,proto3" json:"paused,omitempty"`
}

func (m *MsgSetPausedState) Reset()         { *m = MsgSetPausedState{} }
func (m *MsgSetPausedState) String() string { return proto.CompactTextString(m) }
func (*MsgSetPausedState) ProtoMessage()    {}
func (*MsgSetPausedState) Descriptor() ([]byte, []int) {
	return fileDescriptor_cda63894d37623a5, []int{2}
}
func (m *MsgSetPausedState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetPausedState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetPausedState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetPausedState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetPausedState.Merge(m, src)
}
func (m *MsgSetPausedState) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetPausedState) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetPausedState.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetPausedState proto.InternalMessageInfo

// MsgSetPausedStateResponse is the response of the SetPausedState message.
type MsgSetPausedStateResponse struct {
}

func (m *MsgSetPausedStateResponse) Reset()         { *m = MsgSetPausedStateResponse{} }
func (m *MsgSetPausedStateResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetPausedStateResponse) ProtoMessage()    {}
func (*MsgSetPausedStateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cda63894d37623a5, []int{3}
}
func (m *MsgSetPausedStateResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetPausedStateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetPausedStateResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetPausedStateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetPausedStateResponse.Merge(m, src)
}
func (m *MsgSetPausedStateResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetPausedStateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetPausedStateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetPausedStateResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgClaimYield)(nil), "noble.dollar.v1.MsgClaimYield")
	proto.RegisterType((*MsgClaimYieldResponse)(nil), "noble.dollar.v1.MsgClaimYieldResponse")
	proto.RegisterType((*MsgSetPausedState)(nil), "noble.dollar.v1.MsgSetPausedState")
	proto.RegisterType((*MsgSetPausedStateResponse)(nil), "noble.dollar.v1.MsgSetPausedStateResponse")
}

func init() { proto.RegisterFile("noble/dollar/v1/tx.proto", fileDescriptor_cda63894d37623a5) }

var fileDescriptor_cda63894d37623a5 = []byte{
	// 370 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xc8, 0xcb, 0x4f, 0xca,
	0x49, 0xd5, 0x4f, 0xc9, 0xcf, 0xc9, 0x49, 0x2c, 0xd2, 0x2f, 0x33, 0xd4, 0x2f, 0xa9, 0xd0, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x07, 0xcb, 0xe8, 0x41, 0x64, 0xf4, 0xca, 0x0c, 0xa5, 0x04,
	0x13, 0x73, 0x33, 0xf3, 0xf2, 0xf5, 0xc1, 0x24, 0x44, 0x8d, 0x94, 0x78, 0x72, 0x7e, 0x71, 0x6e,
	0x7e, 0xb1, 0x7e, 0x6e, 0x71, 0x3a, 0x48, 0x6f, 0x6e, 0x71, 0x3a, 0x54, 0x42, 0x12, 0x22, 0x11,
	0x0f, 0xe6, 0xe9, 0x43, 0x38, 0x50, 0x29, 0x91, 0xf4, 0xfc, 0xf4, 0x7c, 0x88, 0x38, 0x88, 0x05,
	0x11, 0x55, 0xca, 0xe1, 0xe2, 0xf5, 0x2d, 0x4e, 0x77, 0xce, 0x49, 0xcc, 0xcc, 0x8d, 0xcc, 0x4c,
	0xcd, 0x49, 0x11, 0x32, 0xe0, 0x62, 0x2b, 0xce, 0x4c, 0xcf, 0x4b, 0x2d, 0x92, 0x60, 0x54, 0x60,
	0xd4, 0xe0, 0x74, 0x92, 0xb8, 0xb4, 0x45, 0x57, 0x04, 0x6a, 0x90, 0x63, 0x4a, 0x4a, 0x51, 0x6a,
	0x71, 0x71, 0x70, 0x49, 0x51, 0x66, 0x5e, 0x7a, 0x10, 0x54, 0x9d, 0x95, 0x66, 0xc7, 0x02, 0x79,
	0x86, 0x17, 0x0b, 0xe4, 0x19, 0x9a, 0x9e, 0x6f, 0xd0, 0x82, 0x0a, 0x76, 0x3d, 0xdf, 0xa0, 0x25,
	0x08, 0xf5, 0x1c, 0xc2, 0x70, 0x25, 0x71, 0x2e, 0x51, 0x14, 0xdb, 0x82, 0x52, 0x8b, 0x0b, 0xf2,
	0xf3, 0x8a, 0x53, 0x95, 0x7a, 0x18, 0xb9, 0x04, 0x7d, 0x8b, 0xd3, 0x83, 0x53, 0x4b, 0x02, 0x12,
	0x4b, 0x8b, 0x53, 0x53, 0x82, 0x4b, 0x12, 0x4b, 0x52, 0x49, 0x77, 0x8b, 0x90, 0x18, 0x17, 0x5b,
	0x01, 0xd8, 0x00, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x8e, 0x20, 0x28, 0xcf, 0x4a, 0x17, 0x87, 0x1b,
	0x45, 0xa1, 0x6e, 0x44, 0xb5, 0x58, 0x49, 0x9a, 0x4b, 0x12, 0xc3, 0x35, 0x30, 0xb7, 0x1a, 0x1d,
	0x61, 0xe4, 0x62, 0xf6, 0x2d, 0x4e, 0x17, 0x0a, 0xe1, 0xe2, 0x42, 0x0a, 0x37, 0x39, 0x3d, 0xb4,
	0x78, 0xd3, 0x43, 0xf1, 0xa9, 0x94, 0x1a, 0x7e, 0x79, 0x98, 0xe9, 0x42, 0x09, 0x5c, 0x7c, 0x68,
	0xa1, 0xa0, 0x84, 0x4d, 0x27, 0xaa, 0x1a, 0x29, 0x2d, 0xc2, 0x6a, 0x60, 0x36, 0x48, 0xb1, 0x36,
	0x3c, 0xdf, 0xa0, 0xc5, 0xe8, 0x64, 0x70, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f,
	0x1e, 0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c,
	0x51, 0x62, 0x50, 0x53, 0x20, 0x46, 0x56, 0x54, 0x56, 0xe9, 0x97, 0x54, 0x16, 0xa4, 0x16, 0x27,
	0xb1, 0x81, 0x93, 0x8c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x61, 0x7d, 0x1c, 0x66, 0xbc, 0x02,
	0x00, 0x00,
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
	ClaimYield(ctx context.Context, in *MsgClaimYield, opts ...grpc.CallOption) (*MsgClaimYieldResponse, error)
	SetPausedState(ctx context.Context, in *MsgSetPausedState, opts ...grpc.CallOption) (*MsgSetPausedStateResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) ClaimYield(ctx context.Context, in *MsgClaimYield, opts ...grpc.CallOption) (*MsgClaimYieldResponse, error) {
	out := new(MsgClaimYieldResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.v1.Msg/ClaimYield", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetPausedState(ctx context.Context, in *MsgSetPausedState, opts ...grpc.CallOption) (*MsgSetPausedStateResponse, error) {
	out := new(MsgSetPausedStateResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.v1.Msg/SetPausedState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	ClaimYield(context.Context, *MsgClaimYield) (*MsgClaimYieldResponse, error)
	SetPausedState(context.Context, *MsgSetPausedState) (*MsgSetPausedStateResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) ClaimYield(ctx context.Context, req *MsgClaimYield) (*MsgClaimYieldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimYield not implemented")
}
func (*UnimplementedMsgServer) SetPausedState(ctx context.Context, req *MsgSetPausedState) (*MsgSetPausedStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPausedState not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_ClaimYield_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimYield)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimYield(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.v1.Msg/ClaimYield",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimYield(ctx, req.(*MsgClaimYield))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetPausedState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPausedState)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPausedState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.v1.Msg/SetPausedState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPausedState(ctx, req.(*MsgSetPausedState))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_serviceDesc = _Msg_serviceDesc
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "noble.dollar.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClaimYield",
			Handler:    _Msg_ClaimYield_Handler,
		},
		{
			MethodName: "SetPausedState",
			Handler:    _Msg_SetPausedState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/dollar/v1/tx.proto",
}

func (m *MsgClaimYield) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimYield) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimYield) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgClaimYieldResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimYieldResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimYieldResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgSetPausedState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetPausedState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetPausedState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Paused {
		i--
		if m.Paused {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
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

func (m *MsgSetPausedStateResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetPausedStateResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetPausedStateResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgClaimYield) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgClaimYieldResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgSetPausedState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Paused {
		n += 2
	}
	return n
}

func (m *MsgSetPausedStateResponse) Size() (n int) {
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
func (m *MsgClaimYield) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgClaimYield: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimYield: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgClaimYieldResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgClaimYieldResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimYieldResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgSetPausedState) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetPausedState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetPausedState: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field Paused", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Paused = bool(v != 0)
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
func (m *MsgSetPausedStateResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetPausedStateResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetPausedStateResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
