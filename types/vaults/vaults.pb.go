// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/vaults/v1/vaults.proto

package vaults

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// buf:lint:ignore ENUM_VALUE_PREFIX
type VaultType int32

const (
	// buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
	UNSPECIFIED VaultType = 0
	STAKED      VaultType = 1
	FLEXIBLE    VaultType = 2
)

var VaultType_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "STAKED",
	2: "FLEXIBLE",
}

var VaultType_value = map[string]int32{
	"UNSPECIFIED": 0,
	"STAKED":      1,
	"FLEXIBLE":    2,
}

func (x VaultType) String() string {
	return proto.EnumName(VaultType_name, int32(x))
}

func (VaultType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1a185cec7ba75cfc, []int{0}
}

// buf:lint:ignore ENUM_VALUE_PREFIX
type PausedType int32

const (
	// buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
	NONE   PausedType = 0
	LOCK   PausedType = 1
	UNLOCK PausedType = 2
	ALL    PausedType = 3
)

var PausedType_name = map[int32]string{
	0: "NONE",
	1: "LOCK",
	2: "UNLOCK",
	3: "ALL",
}

var PausedType_value = map[string]int32{
	"NONE":   0,
	"LOCK":   1,
	"UNLOCK": 2,
	"ALL":    3,
}

func (x PausedType) String() string {
	return proto.EnumName(PausedType_name, int32(x))
}

func (PausedType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1a185cec7ba75cfc, []int{1}
}

type Reward struct {
	Index   cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=index,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"index"`
	Total   cosmossdk_io_math.Int       `protobuf:"bytes,2,opt,name=total,proto3,customtype=cosmossdk.io/math.Int" json:"total"`
	Rewards cosmossdk_io_math.Int       `protobuf:"bytes,3,opt,name=rewards,proto3,customtype=cosmossdk.io/math.Int" json:"rewards"`
}

func (m *Reward) Reset()         { *m = Reward{} }
func (m *Reward) String() string { return proto.CompactTextString(m) }
func (*Reward) ProtoMessage()    {}
func (*Reward) Descriptor() ([]byte, []int) {
	return fileDescriptor_1a185cec7ba75cfc, []int{0}
}
func (m *Reward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Reward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Reward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Reward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Reward.Merge(m, src)
}
func (m *Reward) XXX_Size() int {
	return m.Size()
}
func (m *Reward) XXX_DiscardUnknown() {
	xxx_messageInfo_Reward.DiscardUnknown(m)
}

var xxx_messageInfo_Reward proto.InternalMessageInfo

type Position struct {
	Principal cosmossdk_io_math.Int       `protobuf:"bytes,1,opt,name=principal,proto3,customtype=cosmossdk.io/math.Int" json:"principal"`
	Index     cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=index,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"index"`
	Amount    cosmossdk_io_math.Int       `protobuf:"bytes,3,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
	Time      time.Time                   `protobuf:"bytes,4,opt,name=time,proto3,stdtime" json:"time"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_1a185cec7ba75cfc, []int{1}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Position.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(m, src)
}
func (m *Position) XXX_Size() int {
	return m.Size()
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetTime() time.Time {
	if m != nil {
		return m.Time
	}
	return time.Time{}
}

type PositionEntry struct {
	Address   []byte                      `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Vault     VaultType                   `protobuf:"varint,2,opt,name=vault,proto3,enum=noble.dollar.vaults.v1.VaultType" json:"vault,omitempty"`
	Principal cosmossdk_io_math.Int       `protobuf:"bytes,3,opt,name=principal,proto3,customtype=cosmossdk.io/math.Int" json:"principal"`
	Index     cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=index,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"index"`
	Amount    cosmossdk_io_math.Int       `protobuf:"bytes,5,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
	Time      time.Time                   `protobuf:"bytes,6,opt,name=time,proto3,stdtime" json:"time"`
}

func (m *PositionEntry) Reset()         { *m = PositionEntry{} }
func (m *PositionEntry) String() string { return proto.CompactTextString(m) }
func (*PositionEntry) ProtoMessage()    {}
func (*PositionEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_1a185cec7ba75cfc, []int{2}
}
func (m *PositionEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PositionEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PositionEntry.Merge(m, src)
}
func (m *PositionEntry) XXX_Size() int {
	return m.Size()
}
func (m *PositionEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_PositionEntry.DiscardUnknown(m)
}

var xxx_messageInfo_PositionEntry proto.InternalMessageInfo

func (m *PositionEntry) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *PositionEntry) GetVault() VaultType {
	if m != nil {
		return m.Vault
	}
	return UNSPECIFIED
}

func (m *PositionEntry) GetTime() time.Time {
	if m != nil {
		return m.Time
	}
	return time.Time{}
}

func init() {
	proto.RegisterEnum("noble.dollar.vaults.v1.VaultType", VaultType_name, VaultType_value)
	proto.RegisterEnum("noble.dollar.vaults.v1.PausedType", PausedType_name, PausedType_value)
	proto.RegisterType((*Reward)(nil), "noble.dollar.vaults.v1.Reward")
	proto.RegisterType((*Position)(nil), "noble.dollar.vaults.v1.Position")
	proto.RegisterType((*PositionEntry)(nil), "noble.dollar.vaults.v1.PositionEntry")
}

func init() {
	proto.RegisterFile("noble/dollar/vaults/v1/vaults.proto", fileDescriptor_1a185cec7ba75cfc)
}

var fileDescriptor_1a185cec7ba75cfc = []byte{
	// 561 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xbf, 0x6f, 0xd3, 0x40,
	0x14, 0xc7, 0xed, 0xfc, 0xce, 0x6b, 0x01, 0x73, 0x02, 0x14, 0x82, 0x70, 0x4a, 0x58, 0xaa, 0x48,
	0xd8, 0xb4, 0x48, 0x14, 0x09, 0x96, 0xa6, 0x71, 0x44, 0xa8, 0x95, 0x46, 0x69, 0x8a, 0x10, 0x0b,
	0xba, 0xc4, 0x47, 0xb0, 0xb0, 0x7d, 0x96, 0x7d, 0x09, 0x0d, 0x33, 0x03, 0x63, 0xff, 0x07, 0x16,
	0x06, 0x06, 0x06, 0xfe, 0x88, 0x8e, 0x15, 0x13, 0x62, 0x28, 0x28, 0x19, 0x90, 0xf8, 0x2b, 0x90,
	0xef, 0x2e, 0x14, 0x04, 0x4b, 0xab, 0x2c, 0xd6, 0x7b, 0x77, 0xef, 0x7d, 0xee, 0xbd, 0xf7, 0x7d,
	0x32, 0xdc, 0x0c, 0x68, 0xdf, 0x23, 0xa6, 0x43, 0x3d, 0x0f, 0x47, 0xe6, 0x18, 0x8f, 0x3c, 0x16,
	0x9b, 0xe3, 0x35, 0x69, 0x19, 0x61, 0x44, 0x19, 0x45, 0x57, 0x78, 0x90, 0x21, 0x82, 0x0c, 0x79,
	0x35, 0x5e, 0x2b, 0x5f, 0xc4, 0xbe, 0x1b, 0x50, 0x93, 0x7f, 0x45, 0x68, 0xf9, 0xea, 0x80, 0xc6,
	0x3e, 0x8d, 0x9f, 0x71, 0xcf, 0x14, 0x8e, 0xbc, 0xba, 0x34, 0xa4, 0x43, 0x2a, 0xce, 0x13, 0x4b,
	0x9e, 0x56, 0x86, 0x94, 0x0e, 0x3d, 0x62, 0x72, 0xaf, 0x3f, 0x7a, 0x6e, 0x32, 0xd7, 0x27, 0x31,
	0xc3, 0x7e, 0x28, 0x02, 0xaa, 0x3f, 0x55, 0xc8, 0x75, 0xc9, 0x2b, 0x1c, 0x39, 0xc8, 0x86, 0xac,
	0x1b, 0x38, 0x64, 0xbf, 0xa4, 0xae, 0xa8, 0xab, 0xc5, 0xfa, 0xdd, 0xc3, 0xe3, 0x8a, 0xf2, 0xf5,
	0xb8, 0x72, 0x4d, 0x3c, 0x13, 0x3b, 0x2f, 0x0d, 0x97, 0x9a, 0x3e, 0x66, 0x2f, 0x0c, 0x9b, 0x0c,
	0xf1, 0x60, 0xd2, 0x20, 0x83, 0xcf, 0x9f, 0x6e, 0x81, 0xac, 0xa2, 0x41, 0x06, 0xef, 0x7f, 0x7c,
	0xac, 0xa9, 0x5d, 0x01, 0x41, 0x4d, 0xc8, 0x32, 0xca, 0xb0, 0x57, 0x4a, 0x71, 0xda, 0x6d, 0x49,
	0xbb, 0xfc, 0x2f, 0xad, 0x15, 0xb0, 0x3f, 0x38, 0xad, 0x80, 0x49, 0x0e, 0x4f, 0x47, 0x8f, 0x20,
	0x1f, 0xf1, 0xfa, 0xe2, 0x52, 0xfa, 0x8c, 0xa4, 0x39, 0xa0, 0xfa, 0x21, 0x05, 0x85, 0x0e, 0x8d,
	0x5d, 0xe6, 0xd2, 0x00, 0xb5, 0xa1, 0x18, 0x46, 0x6e, 0x30, 0x70, 0x43, 0xec, 0xc9, 0x96, 0x4f,
	0x8f, 0x3e, 0x41, 0x9c, 0x8c, 0x2f, 0xb5, 0x88, 0xf1, 0x3d, 0x84, 0x1c, 0xf6, 0xe9, 0x28, 0x60,
	0x67, 0xee, 0x5a, 0xe6, 0xa3, 0x7b, 0x90, 0x49, 0x44, 0x2f, 0x65, 0x56, 0xd4, 0xd5, 0xa5, 0xf5,
	0xb2, 0x21, 0x36, 0xc2, 0x98, 0x6f, 0x84, 0xd1, 0x9b, 0x6f, 0x44, 0xbd, 0x90, 0xbc, 0x71, 0xf0,
	0xad, 0xa2, 0x76, 0x79, 0x46, 0xf5, 0x4d, 0x1a, 0xce, 0xcd, 0xc7, 0x65, 0x05, 0x2c, 0x9a, 0xa0,
	0x12, 0xe4, 0xb1, 0xe3, 0x44, 0x24, 0x8e, 0xf9, 0xc4, 0x96, 0xbb, 0x73, 0x17, 0x6d, 0x40, 0x96,
	0x6f, 0x2e, 0xef, 0xfe, 0xfc, 0xfa, 0x0d, 0xe3, 0xff, 0x4b, 0x6d, 0x3c, 0x4e, 0xac, 0xde, 0x24,
	0x24, 0x5d, 0x11, 0xff, 0xb7, 0x0c, 0xe9, 0x05, 0xca, 0x90, 0x59, 0xac, 0x0c, 0xd9, 0x05, 0xc9,
	0x90, 0x3b, 0xad, 0x0c, 0xb5, 0x07, 0x50, 0xfc, 0x3d, 0x35, 0x74, 0x01, 0x96, 0xf6, 0xda, 0xbb,
	0x1d, 0x6b, 0xab, 0xd5, 0x6c, 0x59, 0x0d, 0x4d, 0x41, 0x00, 0xb9, 0xdd, 0xde, 0xe6, 0xb6, 0xd5,
	0xd0, 0x54, 0xb4, 0x0c, 0x85, 0xa6, 0x6d, 0x3d, 0x69, 0xd5, 0x6d, 0x4b, 0x4b, 0x95, 0x33, 0x6f,
	0xdf, 0xe9, 0x4a, 0xed, 0x3e, 0x40, 0x07, 0x8f, 0x62, 0xe2, 0xf0, 0xf4, 0x02, 0x64, 0xda, 0x3b,
	0x6d, 0x4b, 0x53, 0x12, 0xcb, 0xde, 0xd9, 0xda, 0xd6, 0xd4, 0x84, 0xb0, 0xd7, 0xe6, 0x76, 0x0a,
	0xe5, 0x21, 0xbd, 0x69, 0xdb, 0x5a, 0x5a, 0x24, 0xd7, 0x37, 0x0e, 0xa7, 0xba, 0x7a, 0x34, 0xd5,
	0xd5, 0xef, 0x53, 0x5d, 0x3d, 0x98, 0xe9, 0xca, 0xd1, 0x4c, 0x57, 0xbe, 0xcc, 0x74, 0xe5, 0xe9,
	0x75, 0xa9, 0xac, 0x90, 0x79, 0x7f, 0xf2, 0xda, 0x64, 0x93, 0x90, 0xc4, 0xf2, 0xcf, 0xd6, 0xcf,
	0xf1, 0xbe, 0xee, 0xfc, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xf7, 0x0f, 0xfd, 0x24, 0x01, 0x05, 0x00,
	0x00,
}

func (m *Reward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Reward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Reward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Rewards.Size()
		i -= size
		if _, err := m.Rewards.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Total.Size()
		i -= size
		if _, err := m.Total.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Index.Size()
		i -= size
		if _, err := m.Index.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Position) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Position) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Position) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n1, err1 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintVaults(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0x22
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Index.Size()
		i -= size
		if _, err := m.Index.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Principal.Size()
		i -= size
		if _, err := m.Principal.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *PositionEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PositionEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintVaults(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x32
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.Index.Size()
		i -= size
		if _, err := m.Index.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.Principal.Size()
		i -= size
		if _, err := m.Principal.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVaults(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.Vault != 0 {
		i = encodeVarintVaults(dAtA, i, uint64(m.Vault))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintVaults(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintVaults(dAtA []byte, offset int, v uint64) int {
	offset -= sovVaults(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Reward) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Index.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = m.Total.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = m.Rewards.Size()
	n += 1 + l + sovVaults(uint64(l))
	return n
}

func (m *Position) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Principal.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = m.Index.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovVaults(uint64(l))
	return n
}

func (m *PositionEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovVaults(uint64(l))
	}
	if m.Vault != 0 {
		n += 1 + sovVaults(uint64(m.Vault))
	}
	l = m.Principal.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = m.Index.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovVaults(uint64(l))
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovVaults(uint64(l))
	return n
}

func sovVaults(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVaults(x uint64) (n int) {
	return sovVaults(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Reward) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVaults
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
			return fmt.Errorf("proto: Reward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Reward: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Index.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Total", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Total.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVaults(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVaults
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
func (m *Position) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVaults
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
			return fmt.Errorf("proto: Position: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Position: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Principal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Principal.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Index.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVaults(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVaults
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
func (m *PositionEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVaults
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
			return fmt.Errorf("proto: PositionEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PositionEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = append(m.Address[:0], dAtA[iNdEx:postIndex]...)
			if m.Address == nil {
				m.Address = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vault", wireType)
			}
			m.Vault = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Vault |= VaultType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Principal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Principal.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Index.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVaults
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
				return ErrInvalidLengthVaults
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVaults
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVaults(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthVaults
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
func skipVaults(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVaults
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
					return 0, ErrIntOverflowVaults
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
					return 0, ErrIntOverflowVaults
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
				return 0, ErrInvalidLengthVaults
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVaults
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVaults
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVaults        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVaults          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVaults = fmt.Errorf("proto: unexpected end of group")
)
