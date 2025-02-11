// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/portal/v1/genesis.proto

package portal

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// GenesisState defines the genesis state of the Noble Dollar Portal submodule.
type GenesisState struct {
	// owner is the account that controls the Noble Dollar Portal.
	Owner string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	// paused contains the genesis paused state of the Noble Dollar Portal.
	Paused bool `protobuf:"varint,2,opt,name=paused,proto3" json:"paused,omitempty"`
	// peers contains the genesis peers of the Noble Dollar Portal.
	Peers map[uint16]Peer `protobuf:"bytes,3,rep,name=peers,proto3,castkey=uint16" json:"peers" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// nonce contains the next available nonce used for transfers out of the Noble Dollar Portal.
	Nonce uint32 `protobuf:"varint,4,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_aeebf6cda203b8c6, []int{0}
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

func (m *GenesisState) GetPaused() bool {
	if m != nil {
		return m.Paused
	}
	return false
}

func (m *GenesisState) GetPeers() map[uint16]Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *GenesisState) GetNonce() uint32 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "noble.dollar.portal.v1.GenesisState")
	proto.RegisterMapType((map[uint16]Peer)(nil), "noble.dollar.portal.v1.GenesisState.PeersEntry")
}

func init() {
	proto.RegisterFile("noble/dollar/portal/v1/genesis.proto", fileDescriptor_aeebf6cda203b8c6)
}

var fileDescriptor_aeebf6cda203b8c6 = []byte{
	// 345 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x86, 0x33, 0xed, 0x4d, 0xb9, 0x77, 0x7a, 0x2b, 0x32, 0x94, 0x12, 0x8b, 0xa6, 0x41, 0x5d,
	0x64, 0xe3, 0x0c, 0xad, 0xa0, 0xe2, 0xce, 0x82, 0xb8, 0x95, 0x14, 0x5d, 0xb8, 0x91, 0xb4, 0x39,
	0x84, 0x62, 0x9c, 0x09, 0x33, 0xd3, 0x6a, 0x5c, 0xfa, 0x04, 0xee, 0x7d, 0x0d, 0x1f, 0xa2, 0xcb,
	0xe2, 0xca, 0x95, 0x4a, 0xfb, 0x22, 0x92, 0x4c, 0x40, 0x17, 0x76, 0x77, 0xce, 0xe1, 0xfb, 0xff,
	0xff, 0x9c, 0x83, 0x77, 0xb9, 0x18, 0x26, 0xc0, 0x22, 0x91, 0x24, 0xa1, 0x64, 0xa9, 0x90, 0x3a,
	0x4c, 0xd8, 0xb4, 0xcb, 0x62, 0xe0, 0xa0, 0xc6, 0x8a, 0xa6, 0x52, 0x68, 0x41, 0x5a, 0x05, 0x45,
	0x0d, 0x45, 0x0d, 0x45, 0xa7, 0xdd, 0xf6, 0xc6, 0x48, 0xa8, 0x5b, 0xa1, 0xae, 0x0b, 0x8a, 0x99,
	0xc6, 0x48, 0xda, 0xcd, 0x58, 0xc4, 0xc2, 0xcc, 0xf3, 0xaa, 0x9c, 0xee, 0xac, 0x88, 0x2b, 0x2d,
	0x0b, 0x68, 0xfb, 0xb9, 0x82, 0xff, 0x9f, 0x99, 0xfc, 0x81, 0x0e, 0x35, 0x10, 0x8a, 0x6d, 0x71,
	0xc7, 0x41, 0x3a, 0xc8, 0x43, 0xfe, 0xbf, 0xbe, 0xf3, 0xfa, 0xb2, 0xd7, 0x2c, 0xc3, 0x4e, 0xa2,
	0x48, 0x82, 0x52, 0x03, 0x2d, 0xc7, 0x3c, 0x0e, 0x0c, 0x46, 0x5a, 0xb8, 0x96, 0x86, 0x13, 0x05,
	0x91, 0x53, 0xf1, 0x90, 0xff, 0x37, 0x28, 0x3b, 0x72, 0x81, 0xed, 0x14, 0x40, 0x2a, 0xa7, 0xea,
	0x55, 0xfd, 0x7a, 0x8f, 0xd1, 0xdf, 0xcf, 0xa2, 0x3f, 0xc3, 0xe9, 0x79, 0xae, 0x38, 0xe5, 0x5a,
	0x66, 0xfd, 0xb5, 0xd9, 0x7b, 0xc7, 0x7a, 0xfc, 0xe8, 0xd4, 0x26, 0x63, 0xae, 0xbb, 0x07, 0x81,
	0x71, 0x23, 0x4d, 0x6c, 0x73, 0xc1, 0x47, 0xe0, 0xfc, 0xf1, 0x90, 0xdf, 0x08, 0x4c, 0xd3, 0xbe,
	0xc4, 0xf8, 0x5b, 0x4a, 0xd6, 0x71, 0xf5, 0x06, 0xb2, 0xe2, 0x80, 0x46, 0x90, 0x97, 0xa4, 0x87,
	0xed, 0x69, 0x98, 0x4c, 0xa0, 0xd8, 0xb1, 0xde, 0xdb, 0x5c, 0xb5, 0x4c, 0x6e, 0x12, 0x18, 0xf4,
	0xb8, 0x72, 0x84, 0xfa, 0x87, 0xb3, 0x85, 0x8b, 0xe6, 0x0b, 0x17, 0x7d, 0x2e, 0x5c, 0xf4, 0xb4,
	0x74, 0xad, 0xf9, 0xd2, 0xb5, 0xde, 0x96, 0xae, 0x75, 0xb5, 0x55, 0x6a, 0x8d, 0xd1, 0x7d, 0xf6,
	0xc0, 0x74, 0x96, 0x82, 0x2a, 0x9f, 0x3b, 0xac, 0x15, 0xdf, 0xdd, 0xff, 0x0a, 0x00, 0x00, 0xff,
	0xff, 0x3e, 0x89, 0x9f, 0x54, 0xf3, 0x01, 0x00, 0x00,
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
	if m.Nonce != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Peers) > 0 {
		for k := range m.Peers {
			v := m.Peers[k]
			baseI := i
			{
				size, err := (&v).MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
			i = encodeVarintGenesis(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintGenesis(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x1a
		}
	}
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
	if m.Paused {
		n += 2
	}
	if len(m.Peers) > 0 {
		for k, v := range m.Peers {
			_ = k
			_ = v
			l = v.Size()
			mapEntrySize := 1 + sovGenesis(uint64(k)) + 1 + l + sovGenesis(uint64(l))
			n += mapEntrySize + 1 + sovGenesis(uint64(mapEntrySize))
		}
	}
	if m.Nonce != 0 {
		n += 1 + sovGenesis(uint64(m.Nonce))
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Paused", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Peers", wireType)
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
			if m.Peers == nil {
				m.Peers = make(map[uint16]Peer)
			}
			var mapkey uint32
			mapvalue := &Peer{}
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
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
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenesis
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthGenesis
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthGenesis
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &Peer{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenesis(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenesis
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Peers[uint16(mapkey)] = *mapvalue
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint32(b&0x7F) << shift
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
