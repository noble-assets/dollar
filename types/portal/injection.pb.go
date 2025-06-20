// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/portal/v1/injection.proto

package portal

import (
	fmt "fmt"
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

// MsgDeliverInjection is an internal message type used for delivering Noble
// Dollar Portal messages. It is specifically used to insert VAA's into the top
// of a block via ABCI++.
type MsgDeliverInjection struct {
	Vaa []byte `protobuf:"bytes,1,opt,name=vaa,proto3" json:"vaa,omitempty"`
}

func (m *MsgDeliverInjection) Reset()         { *m = MsgDeliverInjection{} }
func (m *MsgDeliverInjection) String() string { return proto.CompactTextString(m) }
func (*MsgDeliverInjection) ProtoMessage()    {}
func (*MsgDeliverInjection) Descriptor() ([]byte, []int) {
	return fileDescriptor_0bc79fe39604a038, []int{0}
}
func (m *MsgDeliverInjection) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeliverInjection) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeliverInjection.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeliverInjection) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeliverInjection.Merge(m, src)
}
func (m *MsgDeliverInjection) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeliverInjection) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeliverInjection.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeliverInjection proto.InternalMessageInfo

func (m *MsgDeliverInjection) GetVaa() []byte {
	if m != nil {
		return m.Vaa
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgDeliverInjection)(nil), "noble.dollar.portal.v1.MsgDeliverInjection")
}

func init() {
	proto.RegisterFile("noble/dollar/portal/v1/injection.proto", fileDescriptor_0bc79fe39604a038)
}

var fileDescriptor_0bc79fe39604a038 = []byte{
	// 163 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0xcb, 0xcb, 0x4f, 0xca,
	0x49, 0xd5, 0x4f, 0xc9, 0xcf, 0xc9, 0x49, 0x2c, 0xd2, 0x2f, 0xc8, 0x2f, 0x2a, 0x49, 0xcc, 0xd1,
	0x2f, 0x33, 0xd4, 0xcf, 0xcc, 0xcb, 0x4a, 0x4d, 0x2e, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0x12, 0x03, 0xab, 0xd3, 0x83, 0xa8, 0xd3, 0x83, 0xa8, 0xd3, 0x2b, 0x33, 0x54,
	0x52, 0xe7, 0x12, 0xf6, 0x2d, 0x4e, 0x77, 0x49, 0xcd, 0xc9, 0x2c, 0x4b, 0x2d, 0xf2, 0x84, 0x69,
	0x12, 0x12, 0xe0, 0x62, 0x2e, 0x4b, 0x4c, 0x94, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x09, 0x02, 0x31,
	0x9d, 0xac, 0x4e, 0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09,
	0x8f, 0xe5, 0x18, 0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a, 0x01, 0x6a, 0x28,
	0xc4, 0x86, 0x8a, 0xca, 0x2a, 0xfd, 0x32, 0x23, 0xfd, 0x92, 0xca, 0x82, 0xd4, 0x62, 0xa8, 0x7b,
	0x92, 0xd8, 0xc0, 0x6e, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x12, 0x05, 0xff, 0x0f, 0xad,
	0x00, 0x00, 0x00,
}

func (m *MsgDeliverInjection) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeliverInjection) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeliverInjection) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Vaa) > 0 {
		i -= len(m.Vaa)
		copy(dAtA[i:], m.Vaa)
		i = encodeVarintInjection(dAtA, i, uint64(len(m.Vaa)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintInjection(dAtA []byte, offset int, v uint64) int {
	offset -= sovInjection(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgDeliverInjection) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Vaa)
	if l > 0 {
		n += 1 + l + sovInjection(uint64(l))
	}
	return n
}

func sovInjection(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozInjection(x uint64) (n int) {
	return sovInjection(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgDeliverInjection) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInjection
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
			return fmt.Errorf("proto: MsgDeliverInjection: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeliverInjection: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vaa", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInjection
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
				return ErrInvalidLengthInjection
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthInjection
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
			skippy, err := skipInjection(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthInjection
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
func skipInjection(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInjection
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
					return 0, ErrIntOverflowInjection
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
					return 0, ErrIntOverflowInjection
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
				return 0, ErrInvalidLengthInjection
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupInjection
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthInjection
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthInjection        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInjection          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupInjection = fmt.Errorf("proto: unexpected end of group")
)
