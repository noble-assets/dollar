// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/portal/v1/events.proto

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

// PeerUpdated is an event emitted whenever a peer is updated.
type PeerUpdated struct {
	Chain          uint16 `protobuf:"varint,1,opt,name=chain,proto3,customtype=uint16" json:"chain"`
	OldTransceiver []byte `protobuf:"bytes,2,opt,name=old_transceiver,json=oldTransceiver,proto3" json:"old_transceiver,omitempty"`
	NewTransceiver []byte `protobuf:"bytes,3,opt,name=new_transceiver,json=newTransceiver,proto3" json:"new_transceiver,omitempty"`
	OldManager     []byte `protobuf:"bytes,4,opt,name=old_manager,json=oldManager,proto3" json:"old_manager,omitempty"`
	NewManager     []byte `protobuf:"bytes,5,opt,name=new_manager,json=newManager,proto3" json:"new_manager,omitempty"`
}

func (m *PeerUpdated) Reset()         { *m = PeerUpdated{} }
func (m *PeerUpdated) String() string { return proto.CompactTextString(m) }
func (*PeerUpdated) ProtoMessage()    {}
func (*PeerUpdated) Descriptor() ([]byte, []int) {
	return fileDescriptor_878c0cf9b5833b22, []int{0}
}
func (m *PeerUpdated) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PeerUpdated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PeerUpdated.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PeerUpdated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerUpdated.Merge(m, src)
}
func (m *PeerUpdated) XXX_Size() int {
	return m.Size()
}
func (m *PeerUpdated) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerUpdated.DiscardUnknown(m)
}

var xxx_messageInfo_PeerUpdated proto.InternalMessageInfo

func (m *PeerUpdated) GetOldTransceiver() []byte {
	if m != nil {
		return m.OldTransceiver
	}
	return nil
}

func (m *PeerUpdated) GetNewTransceiver() []byte {
	if m != nil {
		return m.NewTransceiver
	}
	return nil
}

func (m *PeerUpdated) GetOldManager() []byte {
	if m != nil {
		return m.OldManager
	}
	return nil
}

func (m *PeerUpdated) GetNewManager() []byte {
	if m != nil {
		return m.NewManager
	}
	return nil
}

// BridgingPathSet is an event emitted whenever a supported bridging path is set.
type BridgingPathSet struct {
	DestinationChainId uint16 `protobuf:"varint,1,opt,name=destination_chain_id,json=destinationChainId,proto3,casttype=uint16" json:"destination_chain_id,omitempty"`
	DestinationToken   []byte `protobuf:"bytes,2,opt,name=destination_token,json=destinationToken,proto3" json:"destination_token,omitempty"`
	Supported          bool   `protobuf:"varint,3,opt,name=supported,proto3" json:"supported,omitempty"`
}

func (m *BridgingPathSet) Reset()         { *m = BridgingPathSet{} }
func (m *BridgingPathSet) String() string { return proto.CompactTextString(m) }
func (*BridgingPathSet) ProtoMessage()    {}
func (*BridgingPathSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_878c0cf9b5833b22, []int{1}
}
func (m *BridgingPathSet) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BridgingPathSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BridgingPathSet.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BridgingPathSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BridgingPathSet.Merge(m, src)
}
func (m *BridgingPathSet) XXX_Size() int {
	return m.Size()
}
func (m *BridgingPathSet) XXX_DiscardUnknown() {
	xxx_messageInfo_BridgingPathSet.DiscardUnknown(m)
}

var xxx_messageInfo_BridgingPathSet proto.InternalMessageInfo

func (m *BridgingPathSet) GetDestinationChainId() uint16 {
	if m != nil {
		return m.DestinationChainId
	}
	return 0
}

func (m *BridgingPathSet) GetDestinationToken() []byte {
	if m != nil {
		return m.DestinationToken
	}
	return nil
}

func (m *BridgingPathSet) GetSupported() bool {
	if m != nil {
		return m.Supported
	}
	return false
}

// OwnershipTransferred is an event emitted whenever an ownership transfer occurs.
type OwnershipTransferred struct {
	PreviousOwner string `protobuf:"bytes,1,opt,name=previous_owner,json=previousOwner,proto3" json:"previous_owner,omitempty"`
	NewOwner      string `protobuf:"bytes,2,opt,name=new_owner,json=newOwner,proto3" json:"new_owner,omitempty"`
}

func (m *OwnershipTransferred) Reset()         { *m = OwnershipTransferred{} }
func (m *OwnershipTransferred) String() string { return proto.CompactTextString(m) }
func (*OwnershipTransferred) ProtoMessage()    {}
func (*OwnershipTransferred) Descriptor() ([]byte, []int) {
	return fileDescriptor_878c0cf9b5833b22, []int{2}
}
func (m *OwnershipTransferred) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OwnershipTransferred) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OwnershipTransferred.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OwnershipTransferred) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OwnershipTransferred.Merge(m, src)
}
func (m *OwnershipTransferred) XXX_Size() int {
	return m.Size()
}
func (m *OwnershipTransferred) XXX_DiscardUnknown() {
	xxx_messageInfo_OwnershipTransferred.DiscardUnknown(m)
}

var xxx_messageInfo_OwnershipTransferred proto.InternalMessageInfo

func (m *OwnershipTransferred) GetPreviousOwner() string {
	if m != nil {
		return m.PreviousOwner
	}
	return ""
}

func (m *OwnershipTransferred) GetNewOwner() string {
	if m != nil {
		return m.NewOwner
	}
	return ""
}

func init() {
	proto.RegisterType((*PeerUpdated)(nil), "noble.dollar.portal.v1.PeerUpdated")
	proto.RegisterType((*BridgingPathSet)(nil), "noble.dollar.portal.v1.BridgingPathSet")
	proto.RegisterType((*OwnershipTransferred)(nil), "noble.dollar.portal.v1.OwnershipTransferred")
}

func init() {
	proto.RegisterFile("noble/dollar/portal/v1/events.proto", fileDescriptor_878c0cf9b5833b22)
}

var fileDescriptor_878c0cf9b5833b22 = []byte{
	// 406 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xd2, 0x4f, 0xcf, 0xd2, 0x30,
	0x1c, 0x07, 0xf0, 0xf5, 0xd1, 0xe7, 0xc9, 0x43, 0x11, 0xd0, 0x86, 0x18, 0xe2, 0x9f, 0x41, 0x50,
	0x23, 0x89, 0xc9, 0x16, 0x62, 0xa2, 0x17, 0x4f, 0x78, 0xf2, 0x60, 0x24, 0x13, 0x2f, 0x5c, 0x96,
	0x42, 0x7f, 0x8e, 0xc6, 0xd9, 0x2e, 0x5d, 0xd9, 0xc4, 0x57, 0xe1, 0xd9, 0x37, 0x24, 0x47, 0x8e,
	0xc6, 0x03, 0x31, 0xf0, 0x2e, 0x3c, 0x99, 0xb6, 0x23, 0xec, 0xb9, 0x2d, 0xdf, 0xdf, 0xa7, 0x4d,
	0xbf, 0xcb, 0x0f, 0x3f, 0x11, 0x72, 0x91, 0x42, 0xc8, 0x64, 0x9a, 0x52, 0x15, 0x66, 0x52, 0x69,
	0x9a, 0x86, 0xc5, 0x38, 0x84, 0x02, 0x84, 0xce, 0x83, 0x4c, 0x49, 0x2d, 0xc9, 0x7d, 0x8b, 0x02,
	0x87, 0x02, 0x87, 0x82, 0x62, 0xfc, 0xa0, 0x9b, 0xc8, 0x44, 0x5a, 0x12, 0x9a, 0x2f, 0xa7, 0x87,
	0xbf, 0x10, 0x6e, 0x4e, 0x01, 0xd4, 0xa7, 0x8c, 0x51, 0x0d, 0x8c, 0x3c, 0xc5, 0x97, 0xcb, 0x15,
	0xe5, 0xa2, 0x87, 0x06, 0x68, 0xd4, 0x9a, 0xb4, 0xb7, 0xfb, 0xbe, 0xf7, 0x67, 0xdf, 0xbf, 0x5a,
	0x73, 0xa1, 0xc7, 0xaf, 0x22, 0x37, 0x24, 0xcf, 0x71, 0x47, 0xa6, 0x2c, 0xd6, 0x8a, 0x8a, 0x7c,
	0x09, 0xbc, 0x00, 0xd5, 0xbb, 0x18, 0xa0, 0xd1, 0x9d, 0xa8, 0x2d, 0x53, 0x36, 0x3b, 0xa7, 0x06,
	0x0a, 0x28, 0x6f, 0xc0, 0x5b, 0x0e, 0x0a, 0x28, 0xeb, 0xb0, 0x8f, 0x9b, 0xe6, 0xc6, 0xaf, 0x54,
	0xd0, 0x04, 0x54, 0xef, 0xb6, 0x45, 0x58, 0xa6, 0xec, 0xbd, 0x4b, 0x0c, 0x30, 0x37, 0x9d, 0xc0,
	0xa5, 0x03, 0x02, 0xca, 0x0a, 0x0c, 0x7f, 0x22, 0xdc, 0x99, 0x28, 0xce, 0x12, 0x2e, 0x92, 0x29,
	0xd5, 0xab, 0x8f, 0xa0, 0xc9, 0x1b, 0xdc, 0x65, 0x90, 0x6b, 0x2e, 0xa8, 0xe6, 0x52, 0xc4, 0xf6,
	0xf1, 0x31, 0x67, 0x55, 0x39, 0xfc, 0xef, 0x5c, 0x8c, 0xd4, 0xdc, 0x5b, 0xc3, 0xde, 0x31, 0xf2,
	0x02, 0xdf, 0xab, 0x9f, 0xd6, 0xf2, 0x0b, 0x88, 0xaa, 0xe7, 0xdd, 0xda, 0x60, 0x66, 0x72, 0xf2,
	0x08, 0x37, 0xf2, 0x75, 0x66, 0x7e, 0x37, 0x30, 0xdb, 0xf1, 0x3a, 0x3a, 0x07, 0xc3, 0x39, 0xee,
	0x7e, 0x28, 0x05, 0xa8, 0x7c, 0xc5, 0x33, 0x5b, 0xfb, 0x33, 0x28, 0x05, 0x8c, 0x3c, 0xc3, 0xed,
	0x4c, 0x41, 0xc1, 0xe5, 0x3a, 0x8f, 0xa5, 0x01, 0xf6, 0x69, 0x8d, 0xa8, 0x75, 0x4a, 0xed, 0x29,
	0xf2, 0x10, 0x37, 0x4c, 0x79, 0x27, 0x2e, 0xac, 0xb8, 0x16, 0x50, 0xda, 0xe1, 0xe4, 0xf5, 0xf6,
	0xe0, 0xa3, 0xdd, 0xc1, 0x47, 0x7f, 0x0f, 0x3e, 0xfa, 0x71, 0xf4, 0xbd, 0xdd, 0xd1, 0xf7, 0x7e,
	0x1f, 0x7d, 0x6f, 0xfe, 0xb8, 0x5a, 0x02, 0xb7, 0x11, 0xdf, 0x36, 0xdf, 0x43, 0xbd, 0xc9, 0x20,
	0xaf, 0x36, 0x67, 0x71, 0x65, 0x57, 0xe0, 0xe5, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0x52, 0x6d,
	0x8d, 0x39, 0x57, 0x02, 0x00, 0x00,
}

func (m *PeerUpdated) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PeerUpdated) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PeerUpdated) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.NewManager) > 0 {
		i -= len(m.NewManager)
		copy(dAtA[i:], m.NewManager)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.NewManager)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.OldManager) > 0 {
		i -= len(m.OldManager)
		copy(dAtA[i:], m.OldManager)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.OldManager)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.NewTransceiver) > 0 {
		i -= len(m.NewTransceiver)
		copy(dAtA[i:], m.NewTransceiver)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.NewTransceiver)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.OldTransceiver) > 0 {
		i -= len(m.OldTransceiver)
		copy(dAtA[i:], m.OldTransceiver)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.OldTransceiver)))
		i--
		dAtA[i] = 0x12
	}
	if m.Chain != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.Chain))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *BridgingPathSet) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BridgingPathSet) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BridgingPathSet) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Supported {
		i--
		if m.Supported {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if len(m.DestinationToken) > 0 {
		i -= len(m.DestinationToken)
		copy(dAtA[i:], m.DestinationToken)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.DestinationToken)))
		i--
		dAtA[i] = 0x12
	}
	if m.DestinationChainId != 0 {
		i = encodeVarintEvents(dAtA, i, uint64(m.DestinationChainId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OwnershipTransferred) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OwnershipTransferred) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OwnershipTransferred) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.NewOwner) > 0 {
		i -= len(m.NewOwner)
		copy(dAtA[i:], m.NewOwner)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.NewOwner)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PreviousOwner) > 0 {
		i -= len(m.PreviousOwner)
		copy(dAtA[i:], m.PreviousOwner)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.PreviousOwner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvents(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvents(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PeerUpdated) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Chain != 0 {
		n += 1 + sovEvents(uint64(m.Chain))
	}
	l = len(m.OldTransceiver)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.NewTransceiver)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.OldManager)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.NewManager)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func (m *BridgingPathSet) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.DestinationChainId != 0 {
		n += 1 + sovEvents(uint64(m.DestinationChainId))
	}
	l = len(m.DestinationToken)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	if m.Supported {
		n += 2
	}
	return n
}

func (m *OwnershipTransferred) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PreviousOwner)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	l = len(m.NewOwner)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PeerUpdated) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: PeerUpdated: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PeerUpdated: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Chain", wireType)
			}
			m.Chain = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OldTransceiver", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OldTransceiver = append(m.OldTransceiver[:0], dAtA[iNdEx:postIndex]...)
			if m.OldTransceiver == nil {
				m.OldTransceiver = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewTransceiver", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewTransceiver = append(m.NewTransceiver[:0], dAtA[iNdEx:postIndex]...)
			if m.NewTransceiver == nil {
				m.NewTransceiver = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OldManager", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OldManager = append(m.OldManager[:0], dAtA[iNdEx:postIndex]...)
			if m.OldManager == nil {
				m.OldManager = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewManager", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewManager = append(m.NewManager[:0], dAtA[iNdEx:postIndex]...)
			if m.NewManager == nil {
				m.NewManager = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *BridgingPathSet) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: BridgingPathSet: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BridgingPathSet: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationChainId", wireType)
			}
			m.DestinationChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationToken = append(m.DestinationToken[:0], dAtA[iNdEx:postIndex]...)
			if m.DestinationToken == nil {
				m.DestinationToken = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Supported", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
			m.Supported = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func (m *OwnershipTransferred) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvents
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
			return fmt.Errorf("proto: OwnershipTransferred: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OwnershipTransferred: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousOwner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PreviousOwner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NewOwner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvents
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
				return ErrInvalidLengthEvents
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvents
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NewOwner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvents(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvents
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
func skipEvents(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
					return 0, ErrIntOverflowEvents
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
				return 0, ErrInvalidLengthEvents
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvents
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvents
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvents        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvents          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvents = fmt.Errorf("proto: unexpected end of group")
)
