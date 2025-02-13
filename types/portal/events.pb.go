// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: noble/dollar/portal/v1/events.proto

package portal

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Delivered is the event emitted when a vaa is successfully delivered.
type Delivered struct {
	Vaa []byte `protobuf:"bytes,1,opt,name=vaa,proto3" json:"vaa,omitempty"`
}

func (m *Delivered) Reset()         { *m = Delivered{} }
func (m *Delivered) String() string { return proto.CompactTextString(m) }
func (*Delivered) ProtoMessage()    {}
func (*Delivered) Descriptor() ([]byte, []int) {
	return fileDescriptor_878c0cf9b5833b22, []int{0}
}
func (m *Delivered) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Delivered) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Delivered.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Delivered) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Delivered.Merge(m, src)
}
func (m *Delivered) XXX_Size() int {
	return m.Size()
}
func (m *Delivered) XXX_DiscardUnknown() {
	xxx_messageInfo_Delivered.DiscardUnknown(m)
}

var xxx_messageInfo_Delivered proto.InternalMessageInfo

func (m *Delivered) GetVaa() []byte {
	if m != nil {
		return m.Vaa
	}
	return nil
}

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
	return fileDescriptor_878c0cf9b5833b22, []int{1}
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

// StatePausedUpdated is an event emitted whenever the paused state is updated.
type StatePausedUpdated struct {
	Paused bool `protobuf:"varint,1,opt,name=paused,proto3" json:"paused,omitempty"`
}

func (m *StatePausedUpdated) Reset()         { *m = StatePausedUpdated{} }
func (m *StatePausedUpdated) String() string { return proto.CompactTextString(m) }
func (*StatePausedUpdated) ProtoMessage()    {}
func (*StatePausedUpdated) Descriptor() ([]byte, []int) {
	return fileDescriptor_878c0cf9b5833b22, []int{3}
}

func (m *StatePausedUpdated) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}

func (m *StatePausedUpdated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_StatePausedUpdated.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}

func (m *StatePausedUpdated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatePausedUpdated.Merge(m, src)
}

func (m *StatePausedUpdated) XXX_Size() int {
	return m.Size()
}

func (m *StatePausedUpdated) XXX_DiscardUnknown() {
	xxx_messageInfo_StatePausedUpdated.DiscardUnknown(m)
}

var xxx_messageInfo_StatePausedUpdated proto.InternalMessageInfo

func (m *StatePausedUpdated) GetPaused() bool {
	if m != nil {
		return m.Paused
	}
	return false
}

func init() {
	proto.RegisterType((*Delivered)(nil), "noble.dollar.portal.v1.Delivered")
	proto.RegisterType((*PeerUpdated)(nil), "noble.dollar.portal.v1.PeerUpdated")
	proto.RegisterType((*BridgingPathSet)(nil), "noble.dollar.portal.v1.BridgingPathSet")
	proto.RegisterType((*OwnershipTransferred)(nil), "noble.dollar.portal.v1.OwnershipTransferred")
	proto.RegisterType((*StatePausedUpdated)(nil), "noble.dollar.portal.v1.StatePausedUpdated")
}

func init() {
	proto.RegisterFile("noble/dollar/portal/v1/events.proto", fileDescriptor_878c0cf9b5833b22)
}

var fileDescriptor_878c0cf9b5833b22 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xd1, 0x41, 0x4f, 0xf2, 0x30,
	0x18, 0x07, 0xf0, 0x0d, 0x5e, 0x08, 0x94, 0x17, 0x34, 0x0b, 0x21, 0x44, 0xc3, 0x30, 0x53, 0xa3,
	0x07, 0xb3, 0x85, 0x98, 0xe8, 0x9d, 0x78, 0x35, 0x92, 0xa9, 0x17, 0x2e, 0xa4, 0xd0, 0x47, 0x58,
	0x52, 0xdb, 0xa5, 0x2b, 0x9b, 0xf8, 0x29, 0xfc, 0x56, 0x72, 0xe4, 0x68, 0x3c, 0x10, 0x03, 0x5f,
	0xc4, 0xb4, 0x1d, 0x51, 0x6f, 0xdd, 0xd3, 0xdf, 0xfe, 0xc9, 0xd3, 0x3f, 0x3a, 0x66, 0x7c, 0x4c,
	0x21, 0x20, 0x9c, 0x52, 0x2c, 0x82, 0x98, 0x0b, 0x89, 0x69, 0x90, 0xf6, 0x02, 0x48, 0x81, 0xc9,
	0xc4, 0x8f, 0x05, 0x97, 0xdc, 0x69, 0x69, 0xe4, 0x1b, 0xe4, 0x1b, 0xe4, 0xa7, 0xbd, 0x83, 0xe6,
	0x94, 0x4f, 0xb9, 0x26, 0x81, 0x3a, 0x19, 0xed, 0x75, 0x50, 0xf5, 0x06, 0x68, 0x94, 0x82, 0x00,
	0xe2, 0xec, 0xa3, 0x62, 0x8a, 0x71, 0xdb, 0x3e, 0xb2, 0xcf, 0xff, 0x87, 0xea, 0xe8, 0xbd, 0xdb,
	0xa8, 0x36, 0x00, 0x10, 0x8f, 0x31, 0xc1, 0x12, 0x88, 0x73, 0x82, 0x4a, 0x93, 0x19, 0x8e, 0x98,
	0x36, 0xf5, 0x7e, 0x63, 0xb9, 0xee, 0x5a, 0x9f, 0xeb, 0x6e, 0x79, 0x1e, 0x31, 0xd9, 0xbb, 0x0a,
	0xcd, 0xa5, 0x73, 0x86, 0xf6, 0x38, 0x25, 0x23, 0x29, 0x30, 0x4b, 0x26, 0xa0, 0xc2, 0xdb, 0x05,
	0x9d, 0xd9, 0xe0, 0x94, 0x3c, 0xfc, 0x4c, 0x15, 0x64, 0x90, 0xfd, 0x81, 0x45, 0x03, 0x19, 0x64,
	0xbf, 0x61, 0x17, 0xd5, 0x54, 0xe2, 0x33, 0x66, 0x78, 0x0a, 0xa2, 0xfd, 0x4f, 0x23, 0xc4, 0x29,
	0xb9, 0x35, 0x13, 0x05, 0x54, 0xd2, 0x0e, 0x94, 0x0c, 0x60, 0x90, 0xe5, 0xc0, 0x1b, 0xa2, 0xe6,
	0x5d, 0xc6, 0x40, 0x24, 0xb3, 0x28, 0xd6, 0xc9, 0x4f, 0x20, 0xd4, 0xce, 0xa7, 0xa8, 0x11, 0x0b,
	0x48, 0x23, 0x3e, 0x4f, 0x46, 0x5c, 0x01, 0xbd, 0x5a, 0x35, 0xac, 0xef, 0xa6, 0xfa, 0x2f, 0xe7,
	0x10, 0x55, 0x55, 0xbe, 0x11, 0x05, 0x2d, 0x2a, 0x0c, 0x32, 0x7d, 0xe9, 0x5d, 0x20, 0xe7, 0x5e,
	0x62, 0x09, 0x03, 0x3c, 0x4f, 0x80, 0xec, 0xde, 0xaa, 0x85, 0xca, 0xb1, 0x1e, 0xe8, 0xc4, 0x4a,
	0x98, 0x7f, 0xf5, 0xaf, 0x97, 0x1b, 0xd7, 0x5e, 0x6d, 0x5c, 0xfb, 0x6b, 0xe3, 0xda, 0x6f, 0x5b,
	0xd7, 0x5a, 0x6d, 0x5d, 0xeb, 0x63, 0xeb, 0x5a, 0xc3, 0x4e, 0x5e, 0x9a, 0x69, 0xf0, 0x65, 0xf1,
	0x1a, 0xc8, 0x45, 0x0c, 0x49, 0xde, 0xf4, 0xb8, 0xac, 0x2b, 0xbb, 0xfc, 0x0e, 0x00, 0x00, 0xff,
	0xff, 0x59, 0xe5, 0xbd, 0x2a, 0x07, 0x02, 0x00, 0x00,
}

func (m *Delivered) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Delivered) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Delivered) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Vaa) > 0 {
		i -= len(m.Vaa)
		copy(dAtA[i:], m.Vaa)
		i = encodeVarintEvents(dAtA, i, uint64(len(m.Vaa)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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

func (m *StatePausedUpdated) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *StatePausedUpdated) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *StatePausedUpdated) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
		dAtA[i] = 0x8
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
<<<<<<< HEAD
=======
func (m *Delivered) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Vaa)
	if l > 0 {
		n += 1 + l + sovEvents(uint64(l))
	}
	return n
}
>>>>>>> b9d4c8e (add delivered vaa event)

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

func (m *StatePausedUpdated) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Paused {
		n += 2
	}
	return n
}

func sovEvents(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}

func sozEvents(x uint64) (n int) {
	return sovEvents(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Delivered) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Delivered: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Delivered: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Vaa", wireType)
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
			m.Vaa = append(m.Vaa[:0], dAtA[iNdEx:postIndex]...)
			if m.Vaa == nil {
				m.Vaa = []byte{}
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

func (m *StatePausedUpdated) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: StatePausedUpdated: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: StatePausedUpdated: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Paused", wireType)
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
			m.Paused = bool(v != 0)
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
