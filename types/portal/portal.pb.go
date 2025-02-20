// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/portal/v1/portal.proto

package portal

import (
	fmt "fmt"
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

// Peer is the type that stores information about a peer.
type Peer struct {
	Transceiver []byte `protobuf:"bytes,1,opt,name=transceiver,proto3" json:"transceiver,omitempty"`
	Manager     []byte `protobuf:"bytes,2,opt,name=manager,proto3" json:"manager,omitempty"`
}

func (m *Peer) Reset()         { *m = Peer{} }
func (m *Peer) String() string { return proto.CompactTextString(m) }
func (*Peer) ProtoMessage()    {}
func (*Peer) Descriptor() ([]byte, []int) {
	return fileDescriptor_d4188d482379355d, []int{0}
}
func (m *Peer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Peer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Peer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Peer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Peer.Merge(m, src)
}
func (m *Peer) XXX_Size() int {
	return m.Size()
}
func (m *Peer) XXX_DiscardUnknown() {
	xxx_messageInfo_Peer.DiscardUnknown(m)
}

var xxx_messageInfo_Peer proto.InternalMessageInfo

func (m *Peer) GetTransceiver() []byte {
	if m != nil {
		return m.Transceiver
	}
	return nil
}

func (m *Peer) GetManager() []byte {
	if m != nil {
		return m.Manager
	}
	return nil
}

// BridgingPath is the type that stores information about a supported bridging path.
type BridgingPath struct {
	DestinationChainId uint16 `protobuf:"varint,1,opt,name=destination_chain_id,json=destinationChainId,proto3,casttype=uint16" json:"destination_chain_id,omitempty"`
	DestinationToken   []byte `protobuf:"bytes,2,opt,name=destination_token,json=destinationToken,proto3" json:"destination_token,omitempty"`
}

func (m *BridgingPath) Reset()         { *m = BridgingPath{} }
func (m *BridgingPath) String() string { return proto.CompactTextString(m) }
func (*BridgingPath) ProtoMessage()    {}
func (*BridgingPath) Descriptor() ([]byte, []int) {
	return fileDescriptor_d4188d482379355d, []int{1}
}
func (m *BridgingPath) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BridgingPath) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BridgingPath.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BridgingPath) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgingPath.Merge(m, src)
}
func (m *BridgingPath) XXX_Size() int {
	return m.Size()
}
func (m *BridgingPath) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgingPath.DiscardUnknown(m)
}

var xxx_messageInfo_BridgingPath proto.InternalMessageInfo

func (m *BridgingPath) GetDestinationChainId() uint16 {
	if m != nil {
		return m.DestinationChainId
	}
	return 0
}

func (m *BridgingPath) GetDestinationToken() []byte {
	if m != nil {
		return m.DestinationToken
	}
	return nil
}

func init() {
	proto.RegisterType((*Peer)(nil), "noble.dollar.portal.v1.Peer")
	proto.RegisterType((*BridgingPath)(nil), "noble.dollar.portal.v1.BridgingPath")
}

func init() {
	proto.RegisterFile("noble/dollar/portal/v1/portal.proto", fileDescriptor_d4188d482379355d)
}

var fileDescriptor_d4188d482379355d = []byte{
	// 264 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xce, 0xcb, 0x4f, 0xca,
	0x49, 0xd5, 0x4f, 0xc9, 0xcf, 0xc9, 0x49, 0x2c, 0xd2, 0x2f, 0xc8, 0x2f, 0x2a, 0x49, 0xcc, 0xd1,
	0x2f, 0x33, 0x84, 0xb2, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xc4, 0xc0, 0x8a, 0xf4, 0x20,
	0x8a, 0xf4, 0xa0, 0x52, 0x65, 0x86, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x25, 0xfa, 0x20,
	0x16, 0x44, 0xb5, 0x92, 0x13, 0x17, 0x4b, 0x40, 0x6a, 0x6a, 0x91, 0x90, 0x02, 0x17, 0x77, 0x49,
	0x51, 0x62, 0x5e, 0x71, 0x72, 0x6a, 0x66, 0x59, 0x6a, 0x91, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4f,
	0x10, 0xb2, 0x90, 0x90, 0x04, 0x17, 0x7b, 0x6e, 0x62, 0x5e, 0x62, 0x7a, 0x6a, 0x91, 0x04, 0x13,
	0x58, 0x16, 0xc6, 0x55, 0xaa, 0xe4, 0xe2, 0x71, 0x2a, 0xca, 0x4c, 0x49, 0xcf, 0xcc, 0x4b, 0x0f,
	0x48, 0x2c, 0xc9, 0x10, 0xb2, 0xe1, 0x12, 0x49, 0x49, 0x2d, 0x2e, 0xc9, 0xcc, 0x4b, 0x2c, 0xc9,
	0xcc, 0xcf, 0x8b, 0x4f, 0xce, 0x48, 0xcc, 0xcc, 0x8b, 0xcf, 0x4c, 0x01, 0x1b, 0xca, 0xeb, 0xc4,
	0xf5, 0xeb, 0x9e, 0x3c, 0x5b, 0x69, 0x66, 0x5e, 0x89, 0xa1, 0x59, 0x90, 0x10, 0x92, 0x3a, 0x67,
	0x90, 0x32, 0xcf, 0x14, 0x21, 0x6d, 0x2e, 0x41, 0x64, 0xdd, 0x25, 0xf9, 0xd9, 0xa9, 0x79, 0x50,
	0x1b, 0x05, 0x90, 0x24, 0x42, 0x40, 0xe2, 0x4e, 0xe6, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24,
	0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78,
	0x2c, 0xc7, 0x10, 0x25, 0x0b, 0x0d, 0x00, 0x48, 0x68, 0x54, 0x54, 0x56, 0xe9, 0x97, 0x54, 0x16,
	0xa4, 0x16, 0x43, 0xc3, 0x2a, 0x89, 0x0d, 0xec, 0x7d, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0xf8, 0x16, 0x2b, 0x85, 0x53, 0x01, 0x00, 0x00,
}

func (m *Peer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Peer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Peer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Manager) > 0 {
		i -= len(m.Manager)
		copy(dAtA[i:], m.Manager)
		i = encodeVarintPortal(dAtA, i, uint64(len(m.Manager)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Transceiver) > 0 {
		i -= len(m.Transceiver)
		copy(dAtA[i:], m.Transceiver)
		i = encodeVarintPortal(dAtA, i, uint64(len(m.Transceiver)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *BridgingPath) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BridgingPath) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BridgingPath) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DestinationToken) > 0 {
		i -= len(m.DestinationToken)
		copy(dAtA[i:], m.DestinationToken)
		i = encodeVarintPortal(dAtA, i, uint64(len(m.DestinationToken)))
		i--
		dAtA[i] = 0x12
	}
	if m.DestinationChainId != 0 {
		i = encodeVarintPortal(dAtA, i, uint64(m.DestinationChainId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPortal(dAtA []byte, offset int, v uint64) int {
	offset -= sovPortal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Peer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Transceiver)
	if l > 0 {
		n += 1 + l + sovPortal(uint64(l))
	}
	l = len(m.Manager)
	if l > 0 {
		n += 1 + l + sovPortal(uint64(l))
	}
	return n
}

func (m *BridgingPath) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DestinationChainId != 0 {
		n += 1 + sovPortal(uint64(m.DestinationChainId))
	}
	l = len(m.DestinationToken)
	if l > 0 {
		n += 1 + l + sovPortal(uint64(l))
	}
	return n
}

func sovPortal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPortal(x uint64) (n int) {
	return sovPortal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Peer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPortal
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
			return fmt.Errorf("proto: Peer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Peer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Transceiver", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortal
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
				return ErrInvalidLengthPortal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPortal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Transceiver = append(m.Transceiver[:0], dAtA[iNdEx:postIndex]...)
			if m.Transceiver == nil {
				m.Transceiver = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Manager", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortal
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
				return ErrInvalidLengthPortal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPortal
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
			skippy, err := skipPortal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPortal
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
func (m *BridgingPath) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPortal
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
			return fmt.Errorf("proto: BridgingPath: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BridgingPath: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationChainId", wireType)
			}
			m.DestinationChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DestinationChainId |= uint16(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationToken", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortal
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
				return ErrInvalidLengthPortal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthPortal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationToken = append(m.DestinationToken[:0], dAtA[iNdEx:postIndex]...)
			if m.DestinationToken == nil {
				m.DestinationToken = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPortal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPortal
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
func skipPortal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPortal
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
					return 0, ErrIntOverflowPortal
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
					return 0, ErrIntOverflowPortal
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
				return 0, ErrInvalidLengthPortal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPortal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPortal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPortal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPortal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPortal = fmt.Errorf("proto: unexpected end of group")
)
