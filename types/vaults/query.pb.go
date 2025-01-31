// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/vaults/v1/query.proto

package vaults

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type QueryPositionsByProvider struct {
	Provider string `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
}

func (m *QueryPositionsByProvider) Reset()         { *m = QueryPositionsByProvider{} }
func (m *QueryPositionsByProvider) String() string { return proto.CompactTextString(m) }
func (*QueryPositionsByProvider) ProtoMessage()    {}
func (*QueryPositionsByProvider) Descriptor() ([]byte, []int) {
	return fileDescriptor_958128dd0b4264ed, []int{0}
}
func (m *QueryPositionsByProvider) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPositionsByProvider) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPositionsByProvider.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPositionsByProvider) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPositionsByProvider.Merge(m, src)
}
func (m *QueryPositionsByProvider) XXX_Size() int {
	return m.Size()
}
func (m *QueryPositionsByProvider) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPositionsByProvider.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPositionsByProvider proto.InternalMessageInfo

func (m *QueryPositionsByProvider) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

type QueryPositionsByProviderResponse struct {
	Positions []PositionEntry `protobuf:"bytes,1,rep,name=positions,proto3" json:"positions"`
}

func (m *QueryPositionsByProviderResponse) Reset()         { *m = QueryPositionsByProviderResponse{} }
func (m *QueryPositionsByProviderResponse) String() string { return proto.CompactTextString(m) }
func (*QueryPositionsByProviderResponse) ProtoMessage()    {}
func (*QueryPositionsByProviderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_958128dd0b4264ed, []int{1}
}
func (m *QueryPositionsByProviderResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPositionsByProviderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPositionsByProviderResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPositionsByProviderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPositionsByProviderResponse.Merge(m, src)
}
func (m *QueryPositionsByProviderResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryPositionsByProviderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPositionsByProviderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPositionsByProviderResponse proto.InternalMessageInfo

func (m *QueryPositionsByProviderResponse) GetPositions() []PositionEntry {
	if m != nil {
		return m.Positions
	}
	return nil
}

type QueryPaused struct {
}

func (m *QueryPaused) Reset()         { *m = QueryPaused{} }
func (m *QueryPaused) String() string { return proto.CompactTextString(m) }
func (*QueryPaused) ProtoMessage()    {}
func (*QueryPaused) Descriptor() ([]byte, []int) {
	return fileDescriptor_958128dd0b4264ed, []int{2}
}
func (m *QueryPaused) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPaused) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPaused.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPaused) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPaused.Merge(m, src)
}
func (m *QueryPaused) XXX_Size() int {
	return m.Size()
}
func (m *QueryPaused) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPaused.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPaused proto.InternalMessageInfo

type QueryPausedResponse struct {
	Paused PausedType `protobuf:"varint,1,opt,name=paused,proto3,enum=noble.dollar.vaults.v1.PausedType" json:"paused,omitempty"`
}

func (m *QueryPausedResponse) Reset()         { *m = QueryPausedResponse{} }
func (m *QueryPausedResponse) String() string { return proto.CompactTextString(m) }
func (*QueryPausedResponse) ProtoMessage()    {}
func (*QueryPausedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_958128dd0b4264ed, []int{3}
}
func (m *QueryPausedResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryPausedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryPausedResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryPausedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryPausedResponse.Merge(m, src)
}
func (m *QueryPausedResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryPausedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryPausedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryPausedResponse proto.InternalMessageInfo

func (m *QueryPausedResponse) GetPaused() PausedType {
	if m != nil {
		return m.Paused
	}
	return NONE
}

func init() {
	proto.RegisterType((*QueryPositionsByProvider)(nil), "noble.dollar.vaults.v1.QueryPositionsByProvider")
	proto.RegisterType((*QueryPositionsByProviderResponse)(nil), "noble.dollar.vaults.v1.QueryPositionsByProviderResponse")
	proto.RegisterType((*QueryPaused)(nil), "noble.dollar.vaults.v1.QueryPaused")
	proto.RegisterType((*QueryPausedResponse)(nil), "noble.dollar.vaults.v1.QueryPausedResponse")
}

func init() {
	proto.RegisterFile("noble/dollar/vaults/v1/query.proto", fileDescriptor_958128dd0b4264ed)
}

var fileDescriptor_958128dd0b4264ed = []byte{
	// 396 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xcb, 0x4a, 0xf3, 0x40,
	0x14, 0xce, 0xf4, 0xff, 0x2d, 0x76, 0x8a, 0x2e, 0xa6, 0x22, 0x25, 0x6a, 0x0c, 0x29, 0x42, 0xb1,
	0x92, 0xb1, 0x15, 0xbc, 0x2d, 0x0b, 0x2e, 0xdc, 0xb5, 0xc1, 0x95, 0xbb, 0xd4, 0x0e, 0x21, 0x90,
	0x66, 0xc6, 0x4c, 0x1a, 0x8c, 0xe2, 0xc6, 0x8d, 0x2e, 0x05, 0x5f, 0xc4, 0x07, 0xf0, 0x01, 0xba,
	0x2c, 0xb8, 0x71, 0x25, 0xd2, 0x0a, 0xbe, 0x86, 0x74, 0x72, 0xa9, 0x48, 0x23, 0xba, 0x3b, 0x33,
	0xf3, 0xdd, 0xce, 0x39, 0x03, 0x35, 0x97, 0x76, 0x1c, 0x82, 0xbb, 0xd4, 0x71, 0x4c, 0x0f, 0x07,
	0x66, 0xdf, 0xf1, 0x39, 0x0e, 0xea, 0xf8, 0xbc, 0x4f, 0xbc, 0x50, 0x67, 0x1e, 0xf5, 0x29, 0x5a,
	0x16, 0x18, 0x3d, 0xc2, 0xe8, 0x11, 0x46, 0x0f, 0xea, 0xf2, 0xca, 0x19, 0xe5, 0x3d, 0xca, 0x23,
	0xec, 0x37, 0x92, 0xbc, 0x64, 0x51, 0x8b, 0x8a, 0x12, 0x4f, 0xaa, 0xf8, 0x76, 0xd5, 0xa2, 0xd4,
	0x72, 0x08, 0x36, 0x99, 0x8d, 0x4d, 0xd7, 0xa5, 0xbe, 0xe9, 0xdb, 0xd4, 0xe5, 0xf1, 0x6b, 0x25,
	0x23, 0x4c, 0x6c, 0x29, 0x40, 0xda, 0x2e, 0x2c, 0xb7, 0x27, 0x3e, 0x2d, 0xca, 0x6d, 0x41, 0x6e,
	0x86, 0x2d, 0x8f, 0x06, 0x76, 0x97, 0x78, 0x48, 0x86, 0xf3, 0x2c, 0xae, 0xcb, 0x40, 0x05, 0xd5,
	0x82, 0x91, 0x9e, 0xb5, 0x1e, 0x54, 0xb3, 0x78, 0x06, 0xe1, 0x8c, 0xba, 0x9c, 0xa0, 0x63, 0x58,
	0x60, 0xc9, 0x73, 0x19, 0xa8, 0xff, 0xaa, 0xc5, 0xc6, 0x86, 0x3e, 0xbb, 0x7b, 0x3d, 0xd1, 0x39,
	0x72, 0x7d, 0x2f, 0x6c, 0xfe, 0x1f, 0xbc, 0xae, 0x4b, 0xc6, 0x94, 0xad, 0x2d, 0xc0, 0x62, 0x64,
	0x67, 0xf6, 0x39, 0xe9, 0x6a, 0x6d, 0x58, 0xfa, 0x72, 0x4c, 0x0d, 0x0f, 0x61, 0x9e, 0x89, 0x1b,
	0x11, 0x77, 0xb1, 0xa1, 0x65, 0xba, 0x09, 0xd4, 0x49, 0xc8, 0x88, 0x11, 0x33, 0x1a, 0xc3, 0x1c,
	0x9c, 0x13, 0x9a, 0xe8, 0x09, 0xc0, 0xd2, 0xac, 0x71, 0x6c, 0x67, 0xa9, 0x65, 0x0d, 0x42, 0xde,
	0xff, 0x2b, 0x23, 0xe9, 0x44, 0x3b, 0xb8, 0xfb, 0x78, 0xdc, 0x04, 0x37, 0xcf, 0xef, 0x0f, 0x39,
	0x1d, 0x6d, 0xe1, 0x8c, 0x4d, 0xa6, 0xf3, 0xc1, 0x57, 0xc9, 0x62, 0xae, 0xd1, 0x2d, 0x80, 0xf9,
	0xa8, 0x3f, 0x54, 0xf9, 0xd9, 0x5f, 0x80, 0xe4, 0xda, 0x2f, 0x40, 0x69, 0xae, 0xda, 0x34, 0x97,
	0x8a, 0x94, 0xcc, 0x5c, 0x82, 0xd4, 0xdc, 0x1b, 0x8c, 0x14, 0x30, 0x1c, 0x29, 0xe0, 0x6d, 0xa4,
	0x80, 0xfb, 0xb1, 0x22, 0x0d, 0xc7, 0x8a, 0xf4, 0x32, 0x56, 0xa4, 0xd3, 0xb5, 0xd8, 0x2c, 0x72,
	0xbe, 0x08, 0x2f, 0xb1, 0x1f, 0x32, 0xc2, 0x63, 0x89, 0x4e, 0x5e, 0xfc, 0xcd, 0x9d, 0xcf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x3c, 0x29, 0x91, 0xe6, 0x4f, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	PositionsByProvider(ctx context.Context, in *QueryPositionsByProvider, opts ...grpc.CallOption) (*QueryPositionsByProviderResponse, error)
	Paused(ctx context.Context, in *QueryPaused, opts ...grpc.CallOption) (*QueryPausedResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) PositionsByProvider(ctx context.Context, in *QueryPositionsByProvider, opts ...grpc.CallOption) (*QueryPositionsByProviderResponse, error) {
	out := new(QueryPositionsByProviderResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.vaults.v1.Query/PositionsByProvider", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Paused(ctx context.Context, in *QueryPaused, opts ...grpc.CallOption) (*QueryPausedResponse, error) {
	out := new(QueryPausedResponse)
	err := c.cc.Invoke(ctx, "/noble.dollar.vaults.v1.Query/Paused", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	PositionsByProvider(context.Context, *QueryPositionsByProvider) (*QueryPositionsByProviderResponse, error)
	Paused(context.Context, *QueryPaused) (*QueryPausedResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) PositionsByProvider(ctx context.Context, req *QueryPositionsByProvider) (*QueryPositionsByProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PositionsByProvider not implemented")
}
func (*UnimplementedQueryServer) Paused(ctx context.Context, req *QueryPaused) (*QueryPausedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Paused not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_PositionsByProvider_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPositionsByProvider)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PositionsByProvider(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.vaults.v1.Query/PositionsByProvider",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PositionsByProvider(ctx, req.(*QueryPositionsByProvider))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Paused_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPaused)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Paused(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/noble.dollar.vaults.v1.Query/Paused",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Paused(ctx, req.(*QueryPaused))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_serviceDesc = _Query_serviceDesc
var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "noble.dollar.vaults.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PositionsByProvider",
			Handler:    _Query_PositionsByProvider_Handler,
		},
		{
			MethodName: "Paused",
			Handler:    _Query_Paused_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "noble/dollar/vaults/v1/query.proto",
}

func (m *QueryPositionsByProvider) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPositionsByProvider) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPositionsByProvider) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Provider) > 0 {
		i -= len(m.Provider)
		copy(dAtA[i:], m.Provider)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Provider)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryPositionsByProviderResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPositionsByProviderResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPositionsByProviderResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryPaused) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPaused) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPaused) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryPausedResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryPausedResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryPausedResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Paused != 0 {
		i = encodeVarintQuery(dAtA, i, uint64(m.Paused))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryPositionsByProvider) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Provider)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryPositionsByProviderResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	return n
}

func (m *QueryPaused) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryPausedResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Paused != 0 {
		n += 1 + sovQuery(uint64(m.Paused))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryPositionsByProvider) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPositionsByProvider: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPositionsByProvider: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Provider", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Provider = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPositionsByProviderResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPositionsByProviderResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPositionsByProviderResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Positions = append(m.Positions, PositionEntry{})
			if err := m.Positions[len(m.Positions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPaused) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPaused: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPaused: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryPausedResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryPausedResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryPausedResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Paused", wireType)
			}
			m.Paused = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Paused |= PausedType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
