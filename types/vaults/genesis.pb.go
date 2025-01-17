// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/vaults/v1/genesis.proto

package vaults

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/tx/amino"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// GenesisState defines the genesis state of the Noble Dollar Vaults submodule.
type GenesisState struct {
	// owner is the account that controls the Noble Dollar Vaults.
	Owner string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	// total_flexible_principal contains all the users positions inside Vaults.
	Positions []PositionEntry `protobuf:"bytes,2,rep,name=positions,proto3" json:"positions"`
	// rewards maps the rewards amounts by the index.
	Rewards []Reward `protobuf:"bytes,3,rep,name=rewards,proto3" json:"rewards"`
	// total_flexible_principal contains the total principal amount contained in the flexible Vault.
	TotalFlexiblePrincipal cosmossdk_io_math.Int `protobuf:"bytes,4,opt,name=total_flexible_principal,json=totalFlexiblePrincipal,proto3,customtype=cosmossdk.io/math.Int" json:"total_flexible_principal"`
	Paused                 PausedType            `protobuf:"varint,5,opt,name=paused,proto3,enum=noble.dollar.vaults.v1.PausedType" json:"paused,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_6175d63afe7542ef, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *GenesisState) GetPositions() []PositionEntry {
	if m != nil {
		return m.Positions
	}
	return nil
}

func (m *GenesisState) GetRewards() []Reward {
	if m != nil {
		return m.Rewards
	}
	return nil
}

func (m *GenesisState) GetPaused() PausedType {
	if m != nil {
		return m.Paused
	}
	return NONE
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "noble.dollar.vaults.v1.GenesisState")
}

func init() {
	proto.RegisterFile("noble/dollar/vaults/v1/genesis.proto", fileDescriptor_6175d63afe7542ef)
}

var fileDescriptor_6175d63afe7542ef = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xcd, 0xea, 0x13, 0x31,
	0x14, 0xc5, 0x67, 0xfa, 0x25, 0x8d, 0x22, 0x38, 0xd4, 0x32, 0x16, 0x9c, 0x96, 0xaa, 0x50, 0x84,
	0x26, 0xb6, 0x2e, 0x04, 0x17, 0x82, 0x05, 0x95, 0xee, 0xca, 0xd4, 0x95, 0x9b, 0x92, 0x76, 0xe2,
	0x18, 0x4d, 0x93, 0x21, 0x49, 0x3f, 0xc6, 0xa7, 0xf0, 0x31, 0x5c, 0xba, 0xe8, 0x43, 0x74, 0x59,
	0xba, 0x12, 0x17, 0x45, 0xda, 0x85, 0x6b, 0xdf, 0x40, 0x9a, 0xa4, 0xb8, 0xf9, 0x77, 0x33, 0xcc,
	0xbd, 0xf9, 0x9d, 0x73, 0xef, 0xe1, 0x82, 0xc7, 0x5c, 0x4c, 0x19, 0x41, 0x89, 0x60, 0x0c, 0x4b,
	0xb4, 0xc4, 0x0b, 0xa6, 0x15, 0x5a, 0xf6, 0x50, 0x4a, 0x38, 0x51, 0x54, 0xc1, 0x4c, 0x0a, 0x2d,
	0x82, 0xba, 0xa1, 0xa0, 0xa5, 0xa0, 0xa5, 0xe0, 0xb2, 0xd7, 0xb8, 0x87, 0xe7, 0x94, 0x0b, 0x64,
	0xbe, 0x16, 0x6d, 0x3c, 0x98, 0x09, 0x35, 0x17, 0x6a, 0x62, 0x2a, 0x64, 0x0b, 0xf7, 0x54, 0x4b,
	0x45, 0x2a, 0x6c, 0xff, 0xfc, 0xe7, 0xba, 0x8f, 0xae, 0x6c, 0xe0, 0xa6, 0x18, 0xa8, 0xfd, 0xb7,
	0x00, 0xee, 0xbc, 0xb3, 0x2b, 0x8d, 0x35, 0xd6, 0x24, 0x80, 0xa0, 0x2c, 0x56, 0x9c, 0xc8, 0xd0,
	0x6f, 0xf9, 0x9d, 0xea, 0x20, 0xdc, 0x6f, 0xba, 0x35, 0x37, 0xec, 0x75, 0x92, 0x48, 0xa2, 0xd4,
	0x58, 0x4b, 0xca, 0xd3, 0xd8, 0x62, 0xc1, 0x10, 0x54, 0x33, 0xa1, 0xa8, 0xa6, 0x82, 0xab, 0xb0,
	0xd0, 0x2a, 0x76, 0x6e, 0xf7, 0x9f, 0xc0, 0x9b, 0x53, 0xc1, 0x91, 0x03, 0xdf, 0x70, 0x2d, 0xf3,
	0x41, 0x69, 0x7b, 0x68, 0x7a, 0xf1, 0x7f, 0x75, 0xf0, 0x0a, 0xdc, 0x92, 0x64, 0x85, 0x65, 0xa2,
	0xc2, 0xa2, 0x31, 0x8a, 0xae, 0x19, 0xc5, 0x06, 0x73, 0x0e, 0x17, 0x51, 0xf0, 0x19, 0x84, 0x5a,
	0x68, 0xcc, 0x26, 0x1f, 0x19, 0x59, 0xd3, 0x29, 0x23, 0x93, 0x4c, 0x52, 0x3e, 0xa3, 0x19, 0x66,
	0x61, 0xc9, 0xa4, 0x79, 0x76, 0x16, 0xfc, 0x3a, 0x34, 0xef, 0xdb, 0x44, 0x2a, 0xf9, 0x02, 0xa9,
	0x40, 0x73, 0xac, 0x3f, 0xc1, 0x21, 0xd7, 0xfb, 0x4d, 0x17, 0xb8, 0xa8, 0x43, 0xae, 0xbf, 0xff,
	0xf9, 0xf1, 0xd4, 0x8f, 0xeb, 0xc6, 0xf1, 0xad, 0x33, 0x1c, 0x5d, 0xfc, 0x82, 0x97, 0xa0, 0x92,
	0xe1, 0x85, 0x22, 0x49, 0x58, 0x6e, 0xf9, 0x9d, 0xbb, 0xfd, 0xf6, 0xd5, 0xcc, 0x86, 0x7a, 0x9f,
	0x67, 0x24, 0x76, 0x8a, 0xc1, 0x8b, 0xed, 0x31, 0xf2, 0x77, 0xc7, 0xc8, 0xff, 0x7d, 0x8c, 0xfc,
	0x6f, 0xa7, 0xc8, 0xdb, 0x9d, 0x22, 0xef, 0xe7, 0x29, 0xf2, 0x3e, 0x3c, 0x74, 0x72, 0xeb, 0xb5,
	0xce, 0xbf, 0x22, 0x9d, 0x67, 0x44, 0xb9, 0x93, 0x4d, 0x2b, 0xe6, 0x66, 0xcf, 0xff, 0x05, 0x00,
	0x00, 0xff, 0xff, 0x4b, 0x4d, 0xa7, 0x1a, 0x5c, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Paused != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Paused))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.TotalFlexiblePrincipal.Size()
		i -= size
		if _, err := m.TotalFlexiblePrincipal.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.Rewards) > 0 {
		for iNdEx := len(m.Rewards) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Rewards[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Positions) > 0 {
		for iNdEx := len(m.Positions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Positions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if len(m.Positions) > 0 {
		for _, e := range m.Positions {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Rewards) > 0 {
		for _, e := range m.Rewards {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	l = m.TotalFlexiblePrincipal.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.Paused != 0 {
		n += 1 + sovGenesis(uint64(m.Paused))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Positions = append(m.Positions, PositionEntry{})
			if err := m.Positions[len(m.Positions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rewards = append(m.Rewards, Reward{})
			if err := m.Rewards[len(m.Rewards)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalFlexiblePrincipal", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalFlexiblePrincipal.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Paused", wireType)
			}
			m.Paused = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
