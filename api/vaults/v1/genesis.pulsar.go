// Code generated by protoc-gen-go-pulsar. DO NOT EDIT.
package vaultsv1

import (
	_ "cosmossdk.io/api/amino"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	runtime "github.com/cosmos/cosmos-proto/runtime"
	_ "github.com/cosmos/gogoproto/gogoproto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoiface "google.golang.org/protobuf/runtime/protoiface"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	io "io"
	reflect "reflect"
	sync "sync"
)

var _ protoreflect.List = (*_GenesisState_2_list)(nil)

type _GenesisState_2_list struct {
	list *[]*PositionEntry
}

func (x *_GenesisState_2_list) Len() int {
	if x.list == nil {
		return 0
	}
	return len(*x.list)
}

func (x *_GenesisState_2_list) Get(i int) protoreflect.Value {
	return protoreflect.ValueOfMessage((*x.list)[i].ProtoReflect())
}

func (x *_GenesisState_2_list) Set(i int, value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*PositionEntry)
	(*x.list)[i] = concreteValue
}

func (x *_GenesisState_2_list) Append(value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*PositionEntry)
	*x.list = append(*x.list, concreteValue)
}

func (x *_GenesisState_2_list) AppendMutable() protoreflect.Value {
	v := new(PositionEntry)
	*x.list = append(*x.list, v)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_GenesisState_2_list) Truncate(n int) {
	for i := n; i < len(*x.list); i++ {
		(*x.list)[i] = nil
	}
	*x.list = (*x.list)[:n]
}

func (x *_GenesisState_2_list) NewElement() protoreflect.Value {
	v := new(PositionEntry)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_GenesisState_2_list) IsValid() bool {
	return x.list != nil
}

var _ protoreflect.List = (*_GenesisState_3_list)(nil)

type _GenesisState_3_list struct {
	list *[]*Reward
}

func (x *_GenesisState_3_list) Len() int {
	if x.list == nil {
		return 0
	}
	return len(*x.list)
}

func (x *_GenesisState_3_list) Get(i int) protoreflect.Value {
	return protoreflect.ValueOfMessage((*x.list)[i].ProtoReflect())
}

func (x *_GenesisState_3_list) Set(i int, value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*Reward)
	(*x.list)[i] = concreteValue
}

func (x *_GenesisState_3_list) Append(value protoreflect.Value) {
	valueUnwrapped := value.Message()
	concreteValue := valueUnwrapped.Interface().(*Reward)
	*x.list = append(*x.list, concreteValue)
}

func (x *_GenesisState_3_list) AppendMutable() protoreflect.Value {
	v := new(Reward)
	*x.list = append(*x.list, v)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_GenesisState_3_list) Truncate(n int) {
	for i := n; i < len(*x.list); i++ {
		(*x.list)[i] = nil
	}
	*x.list = (*x.list)[:n]
}

func (x *_GenesisState_3_list) NewElement() protoreflect.Value {
	v := new(Reward)
	return protoreflect.ValueOfMessage(v.ProtoReflect())
}

func (x *_GenesisState_3_list) IsValid() bool {
	return x.list != nil
}

var (
	md_GenesisState                          protoreflect.MessageDescriptor
	fd_GenesisState_positions                protoreflect.FieldDescriptor
	fd_GenesisState_rewards                  protoreflect.FieldDescriptor
	fd_GenesisState_total_flexible_principal protoreflect.FieldDescriptor
	fd_GenesisState_paused                   protoreflect.FieldDescriptor
)

func init() {
	file_noble_dollar_vaults_v1_genesis_proto_init()
	md_GenesisState = File_noble_dollar_vaults_v1_genesis_proto.Messages().ByName("GenesisState")
	fd_GenesisState_positions = md_GenesisState.Fields().ByName("positions")
	fd_GenesisState_rewards = md_GenesisState.Fields().ByName("rewards")
	fd_GenesisState_total_flexible_principal = md_GenesisState.Fields().ByName("total_flexible_principal")
	fd_GenesisState_paused = md_GenesisState.Fields().ByName("paused")
}

var _ protoreflect.Message = (*fastReflection_GenesisState)(nil)

type fastReflection_GenesisState GenesisState

func (x *GenesisState) ProtoReflect() protoreflect.Message {
	return (*fastReflection_GenesisState)(x)
}

func (x *GenesisState) slowProtoReflect() protoreflect.Message {
	mi := &file_noble_dollar_vaults_v1_genesis_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

var _fastReflection_GenesisState_messageType fastReflection_GenesisState_messageType
var _ protoreflect.MessageType = fastReflection_GenesisState_messageType{}

type fastReflection_GenesisState_messageType struct{}

func (x fastReflection_GenesisState_messageType) Zero() protoreflect.Message {
	return (*fastReflection_GenesisState)(nil)
}
func (x fastReflection_GenesisState_messageType) New() protoreflect.Message {
	return new(fastReflection_GenesisState)
}
func (x fastReflection_GenesisState_messageType) Descriptor() protoreflect.MessageDescriptor {
	return md_GenesisState
}

// Descriptor returns message descriptor, which contains only the protobuf
// type information for the message.
func (x *fastReflection_GenesisState) Descriptor() protoreflect.MessageDescriptor {
	return md_GenesisState
}

// Type returns the message type, which encapsulates both Go and protobuf
// type information. If the Go type information is not needed,
// it is recommended that the message descriptor be used instead.
func (x *fastReflection_GenesisState) Type() protoreflect.MessageType {
	return _fastReflection_GenesisState_messageType
}

// New returns a newly allocated and mutable empty message.
func (x *fastReflection_GenesisState) New() protoreflect.Message {
	return new(fastReflection_GenesisState)
}

// Interface unwraps the message reflection interface and
// returns the underlying ProtoMessage interface.
func (x *fastReflection_GenesisState) Interface() protoreflect.ProtoMessage {
	return (*GenesisState)(x)
}

// Range iterates over every populated field in an undefined order,
// calling f for each field descriptor and value encountered.
// Range returns immediately if f returns false.
// While iterating, mutating operations may only be performed
// on the current field descriptor.
func (x *fastReflection_GenesisState) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if len(x.Positions) != 0 {
		value := protoreflect.ValueOfList(&_GenesisState_2_list{list: &x.Positions})
		if !f(fd_GenesisState_positions, value) {
			return
		}
	}
	if len(x.Rewards) != 0 {
		value := protoreflect.ValueOfList(&_GenesisState_3_list{list: &x.Rewards})
		if !f(fd_GenesisState_rewards, value) {
			return
		}
	}
	if x.TotalFlexiblePrincipal != "" {
		value := protoreflect.ValueOfString(x.TotalFlexiblePrincipal)
		if !f(fd_GenesisState_total_flexible_principal, value) {
			return
		}
	}
	if x.Paused != 0 {
		value := protoreflect.ValueOfEnum((protoreflect.EnumNumber)(x.Paused))
		if !f(fd_GenesisState_paused, value) {
			return
		}
	}
}

// Has reports whether a field is populated.
//
// Some fields have the property of nullability where it is possible to
// distinguish between the default value of a field and whether the field
// was explicitly populated with the default value. Singular message fields,
// member fields of a oneof, and proto2 scalar fields are nullable. Such
// fields are populated only if explicitly set.
//
// In other cases (aside from the nullable cases above),
// a proto3 scalar field is populated if it contains a non-zero value, and
// a repeated field is populated if it is non-empty.
func (x *fastReflection_GenesisState) Has(fd protoreflect.FieldDescriptor) bool {
	switch fd.FullName() {
	case "noble.dollar.vaults.v1.GenesisState.positions":
		return len(x.Positions) != 0
	case "noble.dollar.vaults.v1.GenesisState.rewards":
		return len(x.Rewards) != 0
	case "noble.dollar.vaults.v1.GenesisState.total_flexible_principal":
		return x.TotalFlexiblePrincipal != ""
	case "noble.dollar.vaults.v1.GenesisState.paused":
		return x.Paused != 0
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.dollar.vaults.v1.GenesisState"))
		}
		panic(fmt.Errorf("message noble.dollar.vaults.v1.GenesisState does not contain field %s", fd.FullName()))
	}
}

// Clear clears the field such that a subsequent Has call reports false.
//
// Clearing an extension field clears both the extension type and value
// associated with the given field number.
//
// Clear is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_GenesisState) Clear(fd protoreflect.FieldDescriptor) {
	switch fd.FullName() {
	case "noble.dollar.vaults.v1.GenesisState.positions":
		x.Positions = nil
	case "noble.dollar.vaults.v1.GenesisState.rewards":
		x.Rewards = nil
	case "noble.dollar.vaults.v1.GenesisState.total_flexible_principal":
		x.TotalFlexiblePrincipal = ""
	case "noble.dollar.vaults.v1.GenesisState.paused":
		x.Paused = 0
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.dollar.vaults.v1.GenesisState"))
		}
		panic(fmt.Errorf("message noble.dollar.vaults.v1.GenesisState does not contain field %s", fd.FullName()))
	}
}

// Get retrieves the value for a field.
//
// For unpopulated scalars, it returns the default value, where
// the default value of a bytes scalar is guaranteed to be a copy.
// For unpopulated composite types, it returns an empty, read-only view
// of the value; to obtain a mutable reference, use Mutable.
func (x *fastReflection_GenesisState) Get(descriptor protoreflect.FieldDescriptor) protoreflect.Value {
	switch descriptor.FullName() {
	case "noble.dollar.vaults.v1.GenesisState.positions":
		if len(x.Positions) == 0 {
			return protoreflect.ValueOfList(&_GenesisState_2_list{})
		}
		listValue := &_GenesisState_2_list{list: &x.Positions}
		return protoreflect.ValueOfList(listValue)
	case "noble.dollar.vaults.v1.GenesisState.rewards":
		if len(x.Rewards) == 0 {
			return protoreflect.ValueOfList(&_GenesisState_3_list{})
		}
		listValue := &_GenesisState_3_list{list: &x.Rewards}
		return protoreflect.ValueOfList(listValue)
	case "noble.dollar.vaults.v1.GenesisState.total_flexible_principal":
		value := x.TotalFlexiblePrincipal
		return protoreflect.ValueOfString(value)
	case "noble.dollar.vaults.v1.GenesisState.paused":
		value := x.Paused
		return protoreflect.ValueOfEnum((protoreflect.EnumNumber)(value))
	default:
		if descriptor.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.dollar.vaults.v1.GenesisState"))
		}
		panic(fmt.Errorf("message noble.dollar.vaults.v1.GenesisState does not contain field %s", descriptor.FullName()))
	}
}

// Set stores the value for a field.
//
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType.
// When setting a composite type, it is unspecified whether the stored value
// aliases the source's memory in any way. If the composite value is an
// empty, read-only value, then it panics.
//
// Set is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_GenesisState) Set(fd protoreflect.FieldDescriptor, value protoreflect.Value) {
	switch fd.FullName() {
	case "noble.dollar.vaults.v1.GenesisState.positions":
		lv := value.List()
		clv := lv.(*_GenesisState_2_list)
		x.Positions = *clv.list
	case "noble.dollar.vaults.v1.GenesisState.rewards":
		lv := value.List()
		clv := lv.(*_GenesisState_3_list)
		x.Rewards = *clv.list
	case "noble.dollar.vaults.v1.GenesisState.total_flexible_principal":
		x.TotalFlexiblePrincipal = value.Interface().(string)
	case "noble.dollar.vaults.v1.GenesisState.paused":
		x.Paused = (PausedType)(value.Enum())
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.dollar.vaults.v1.GenesisState"))
		}
		panic(fmt.Errorf("message noble.dollar.vaults.v1.GenesisState does not contain field %s", fd.FullName()))
	}
}

// Mutable returns a mutable reference to a composite type.
//
// If the field is unpopulated, it may allocate a composite value.
// For a field belonging to a oneof, it implicitly clears any other field
// that may be currently set within the same oneof.
// For extension fields, it implicitly stores the provided ExtensionType
// if not already stored.
// It panics if the field does not contain a composite type.
//
// Mutable is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_GenesisState) Mutable(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "noble.dollar.vaults.v1.GenesisState.positions":
		if x.Positions == nil {
			x.Positions = []*PositionEntry{}
		}
		value := &_GenesisState_2_list{list: &x.Positions}
		return protoreflect.ValueOfList(value)
	case "noble.dollar.vaults.v1.GenesisState.rewards":
		if x.Rewards == nil {
			x.Rewards = []*Reward{}
		}
		value := &_GenesisState_3_list{list: &x.Rewards}
		return protoreflect.ValueOfList(value)
	case "noble.dollar.vaults.v1.GenesisState.total_flexible_principal":
		panic(fmt.Errorf("field total_flexible_principal of message noble.dollar.vaults.v1.GenesisState is not mutable"))
	case "noble.dollar.vaults.v1.GenesisState.paused":
		panic(fmt.Errorf("field paused of message noble.dollar.vaults.v1.GenesisState is not mutable"))
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.dollar.vaults.v1.GenesisState"))
		}
		panic(fmt.Errorf("message noble.dollar.vaults.v1.GenesisState does not contain field %s", fd.FullName()))
	}
}

// NewField returns a new value that is assignable to the field
// for the given descriptor. For scalars, this returns the default value.
// For lists, maps, and messages, this returns a new, empty, mutable value.
func (x *fastReflection_GenesisState) NewField(fd protoreflect.FieldDescriptor) protoreflect.Value {
	switch fd.FullName() {
	case "noble.dollar.vaults.v1.GenesisState.positions":
		list := []*PositionEntry{}
		return protoreflect.ValueOfList(&_GenesisState_2_list{list: &list})
	case "noble.dollar.vaults.v1.GenesisState.rewards":
		list := []*Reward{}
		return protoreflect.ValueOfList(&_GenesisState_3_list{list: &list})
	case "noble.dollar.vaults.v1.GenesisState.total_flexible_principal":
		return protoreflect.ValueOfString("")
	case "noble.dollar.vaults.v1.GenesisState.paused":
		return protoreflect.ValueOfEnum(0)
	default:
		if fd.IsExtension() {
			panic(fmt.Errorf("proto3 declared messages do not support extensions: noble.dollar.vaults.v1.GenesisState"))
		}
		panic(fmt.Errorf("message noble.dollar.vaults.v1.GenesisState does not contain field %s", fd.FullName()))
	}
}

// WhichOneof reports which field within the oneof is populated,
// returning nil if none are populated.
// It panics if the oneof descriptor does not belong to this message.
func (x *fastReflection_GenesisState) WhichOneof(d protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	switch d.FullName() {
	default:
		panic(fmt.Errorf("%s is not a oneof field in noble.dollar.vaults.v1.GenesisState", d.FullName()))
	}
	panic("unreachable")
}

// GetUnknown retrieves the entire list of unknown fields.
// The caller may only mutate the contents of the RawFields
// if the mutated bytes are stored back into the message with SetUnknown.
func (x *fastReflection_GenesisState) GetUnknown() protoreflect.RawFields {
	return x.unknownFields
}

// SetUnknown stores an entire list of unknown fields.
// The raw fields must be syntactically valid according to the wire format.
// An implementation may panic if this is not the case.
// Once stored, the caller must not mutate the content of the RawFields.
// An empty RawFields may be passed to clear the fields.
//
// SetUnknown is a mutating operation and unsafe for concurrent use.
func (x *fastReflection_GenesisState) SetUnknown(fields protoreflect.RawFields) {
	x.unknownFields = fields
}

// IsValid reports whether the message is valid.
//
// An invalid message is an empty, read-only value.
//
// An invalid message often corresponds to a nil pointer of the concrete
// message type, but the details are implementation dependent.
// Validity is not part of the protobuf data model, and may not
// be preserved in marshaling or other operations.
func (x *fastReflection_GenesisState) IsValid() bool {
	return x != nil
}

// ProtoMethods returns optional fastReflectionFeature-path implementations of various operations.
// This method may return nil.
//
// The returned methods type is identical to
// "google.golang.org/protobuf/runtime/protoiface".Methods.
// Consult the protoiface package documentation for details.
func (x *fastReflection_GenesisState) ProtoMethods() *protoiface.Methods {
	size := func(input protoiface.SizeInput) protoiface.SizeOutput {
		x := input.Message.Interface().(*GenesisState)
		if x == nil {
			return protoiface.SizeOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Size:              0,
			}
		}
		options := runtime.SizeInputToOptions(input)
		_ = options
		var n int
		var l int
		_ = l
		if len(x.Positions) > 0 {
			for _, e := range x.Positions {
				l = options.Size(e)
				n += 1 + l + runtime.Sov(uint64(l))
			}
		}
		if len(x.Rewards) > 0 {
			for _, e := range x.Rewards {
				l = options.Size(e)
				n += 1 + l + runtime.Sov(uint64(l))
			}
		}
		l = len(x.TotalFlexiblePrincipal)
		if l > 0 {
			n += 1 + l + runtime.Sov(uint64(l))
		}
		if x.Paused != 0 {
			n += 1 + runtime.Sov(uint64(x.Paused))
		}
		if x.unknownFields != nil {
			n += len(x.unknownFields)
		}
		return protoiface.SizeOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Size:              n,
		}
	}

	marshal := func(input protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
		x := input.Message.Interface().(*GenesisState)
		if x == nil {
			return protoiface.MarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Buf:               input.Buf,
			}, nil
		}
		options := runtime.MarshalInputToOptions(input)
		_ = options
		size := options.Size(x)
		dAtA := make([]byte, size)
		i := len(dAtA)
		_ = i
		var l int
		_ = l
		if x.unknownFields != nil {
			i -= len(x.unknownFields)
			copy(dAtA[i:], x.unknownFields)
		}
		if x.Paused != 0 {
			i = runtime.EncodeVarint(dAtA, i, uint64(x.Paused))
			i--
			dAtA[i] = 0x28
		}
		if len(x.TotalFlexiblePrincipal) > 0 {
			i -= len(x.TotalFlexiblePrincipal)
			copy(dAtA[i:], x.TotalFlexiblePrincipal)
			i = runtime.EncodeVarint(dAtA, i, uint64(len(x.TotalFlexiblePrincipal)))
			i--
			dAtA[i] = 0x22
		}
		if len(x.Rewards) > 0 {
			for iNdEx := len(x.Rewards) - 1; iNdEx >= 0; iNdEx-- {
				encoded, err := options.Marshal(x.Rewards[iNdEx])
				if err != nil {
					return protoiface.MarshalOutput{
						NoUnkeyedLiterals: input.NoUnkeyedLiterals,
						Buf:               input.Buf,
					}, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
				i--
				dAtA[i] = 0x1a
			}
		}
		if len(x.Positions) > 0 {
			for iNdEx := len(x.Positions) - 1; iNdEx >= 0; iNdEx-- {
				encoded, err := options.Marshal(x.Positions[iNdEx])
				if err != nil {
					return protoiface.MarshalOutput{
						NoUnkeyedLiterals: input.NoUnkeyedLiterals,
						Buf:               input.Buf,
					}, err
				}
				i -= len(encoded)
				copy(dAtA[i:], encoded)
				i = runtime.EncodeVarint(dAtA, i, uint64(len(encoded)))
				i--
				dAtA[i] = 0x12
			}
		}
		if input.Buf != nil {
			input.Buf = append(input.Buf, dAtA...)
		} else {
			input.Buf = dAtA
		}
		return protoiface.MarshalOutput{
			NoUnkeyedLiterals: input.NoUnkeyedLiterals,
			Buf:               input.Buf,
		}, nil
	}
	unmarshal := func(input protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
		x := input.Message.Interface().(*GenesisState)
		if x == nil {
			return protoiface.UnmarshalOutput{
				NoUnkeyedLiterals: input.NoUnkeyedLiterals,
				Flags:             input.Flags,
			}, nil
		}
		options := runtime.UnmarshalInputToOptions(input)
		_ = options
		dAtA := input.Buf
		l := len(dAtA)
		iNdEx := 0
		for iNdEx < l {
			preIndex := iNdEx
			var wire uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
				}
				if iNdEx >= l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
			}
			if fieldNum <= 0 {
				return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
			}
			switch fieldNum {
			case 2:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Positions", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Positions = append(x.Positions, &PositionEntry{})
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.Positions[len(x.Positions)-1]); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 3:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
				}
				var msglen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					msglen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if msglen < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + msglen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.Rewards = append(x.Rewards, &Reward{})
				if err := options.Unmarshal(dAtA[iNdEx:postIndex], x.Rewards[len(x.Rewards)-1]); err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				iNdEx = postIndex
			case 4:
				if wireType != 2 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field TotalFlexiblePrincipal", wireType)
				}
				var stringLen uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
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
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				postIndex := iNdEx + intStringLen
				if postIndex < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if postIndex > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				x.TotalFlexiblePrincipal = string(dAtA[iNdEx:postIndex])
				iNdEx = postIndex
			case 5:
				if wireType != 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, fmt.Errorf("proto: wrong wireType = %d for field Paused", wireType)
				}
				x.Paused = 0
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrIntOverflow
					}
					if iNdEx >= l {
						return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					x.Paused |= PausedType(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
			default:
				iNdEx = preIndex
				skippy, err := runtime.Skip(dAtA[iNdEx:])
				if err != nil {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, err
				}
				if (skippy < 0) || (iNdEx+skippy) < 0 {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, runtime.ErrInvalidLength
				}
				if (iNdEx + skippy) > l {
					return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
				}
				if !options.DiscardUnknown {
					x.unknownFields = append(x.unknownFields, dAtA[iNdEx:iNdEx+skippy]...)
				}
				iNdEx += skippy
			}
		}

		if iNdEx > l {
			return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, io.ErrUnexpectedEOF
		}
		return protoiface.UnmarshalOutput{NoUnkeyedLiterals: input.NoUnkeyedLiterals, Flags: input.Flags}, nil
	}
	return &protoiface.Methods{
		NoUnkeyedLiterals: struct{}{},
		Flags:             protoiface.SupportMarshalDeterministic | protoiface.SupportUnmarshalDiscardUnknown,
		Size:              size,
		Marshal:           marshal,
		Unmarshal:         unmarshal,
		Merge:             nil,
		CheckInitialized:  nil,
	}
}

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.0
// 	protoc        (unknown)
// source: noble/dollar/vaults/v1/genesis.proto

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// GenesisState defines the genesis state of the Noble Dollar Vaults submodule.
type GenesisState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// total_flexible_principal contains all the users positions inside Vaults.
	Positions []*PositionEntry `protobuf:"bytes,2,rep,name=positions,proto3" json:"positions,omitempty"`
	// rewards maps the rewards amounts by the index.
	Rewards []*Reward `protobuf:"bytes,3,rep,name=rewards,proto3" json:"rewards,omitempty"`
	// total_flexible_principal contains the total principal amount contained in the flexible Vault.
	TotalFlexiblePrincipal string `protobuf:"bytes,4,opt,name=total_flexible_principal,json=totalFlexiblePrincipal,proto3" json:"total_flexible_principal,omitempty"`
	// paused represents the current pausing state of the Vaults.
	Paused PausedType `protobuf:"varint,5,opt,name=paused,proto3,enum=noble.dollar.vaults.v1.PausedType" json:"paused,omitempty"`
}

func (x *GenesisState) Reset() {
	*x = GenesisState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_noble_dollar_vaults_v1_genesis_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenesisState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenesisState) ProtoMessage() {}

// Deprecated: Use GenesisState.ProtoReflect.Descriptor instead.
func (*GenesisState) Descriptor() ([]byte, []int) {
	return file_noble_dollar_vaults_v1_genesis_proto_rawDescGZIP(), []int{0}
}

func (x *GenesisState) GetPositions() []*PositionEntry {
	if x != nil {
		return x.Positions
	}
	return nil
}

func (x *GenesisState) GetRewards() []*Reward {
	if x != nil {
		return x.Rewards
	}
	return nil
}

func (x *GenesisState) GetTotalFlexiblePrincipal() string {
	if x != nil {
		return x.TotalFlexiblePrincipal
	}
	return ""
}

func (x *GenesisState) GetPaused() PausedType {
	if x != nil {
		return x.Paused
	}
	return PausedType_NONE
}

var File_noble_dollar_vaults_v1_genesis_proto protoreflect.FileDescriptor

var file_noble_dollar_vaults_v1_genesis_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2f, 0x64, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x2f, 0x76,
	0x61, 0x75, 0x6c, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x64, 0x6f,
	0x6c, 0x6c, 0x61, 0x72, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x11,
	0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2f, 0x61, 0x6d, 0x69, 0x6e, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x19, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x67, 0x6f,
	0x67, 0x6f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x23, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2f, 0x64, 0x6f, 0x6c, 0x6c, 0x61, 0x72,
	0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x75, 0x6c, 0x74,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc1, 0x02, 0x0a, 0x0c, 0x47, 0x65, 0x6e, 0x65,
	0x73, 0x69, 0x73, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x49, 0x0a, 0x09, 0x70, 0x6f, 0x73, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x6e, 0x6f,
	0x62, 0x6c, 0x65, 0x2e, 0x64, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x09, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x12, 0x3e, 0x0a, 0x07, 0x72, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x64, 0x6f, 0x6c,
	0x6c, 0x61, 0x72, 0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x77, 0x61, 0x72, 0x64, 0x42, 0x04, 0xc8, 0xde, 0x1f, 0x00, 0x52, 0x07, 0x72, 0x65, 0x77, 0x61,
	0x72, 0x64, 0x73, 0x12, 0x6a, 0x0a, 0x18, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x66, 0x6c, 0x65,
	0x78, 0x69, 0x62, 0x6c, 0x65, 0x5f, 0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x30, 0xc8, 0xde, 0x1f, 0x00, 0xda, 0xde, 0x1f, 0x15, 0x63,
	0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x73, 0x64, 0x6b, 0x2e, 0x69, 0x6f, 0x2f, 0x6d, 0x61, 0x74, 0x68,
	0x2e, 0x49, 0x6e, 0x74, 0xd2, 0xb4, 0x2d, 0x0a, 0x63, 0x6f, 0x73, 0x6d, 0x6f, 0x73, 0x2e, 0x49,
	0x6e, 0x74, 0xa8, 0xe7, 0xb0, 0x2a, 0x01, 0x52, 0x16, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x46, 0x6c,
	0x65, 0x78, 0x69, 0x62, 0x6c, 0x65, 0x50, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x12,
	0x3a, 0x0a, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x22, 0x2e, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x64, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x2e, 0x76,
	0x61, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x75, 0x73, 0x65, 0x64, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x06, 0x70, 0x61, 0x75, 0x73, 0x65, 0x64, 0x42, 0xdb, 0x01, 0x0a, 0x1a,
	0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x64, 0x6f, 0x6c, 0x6c, 0x61, 0x72,
	0x2e, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x47, 0x65, 0x6e, 0x65,
	0x73, 0x69, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x34, 0x64, 0x6f, 0x6c, 0x6c,
	0x61, 0x72, 0x2e, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x78, 0x79, 0x7a, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x6e, 0x6f, 0x62, 0x6c, 0x65, 0x2f, 0x64, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x2f, 0x76, 0x61,
	0x75, 0x6c, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x76, 0x31,
	0xa2, 0x02, 0x03, 0x4e, 0x44, 0x56, 0xaa, 0x02, 0x16, 0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x2e, 0x44,
	0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x2e, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x16, 0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x5c, 0x44, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x5c, 0x56,
	0x61, 0x75, 0x6c, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x22, 0x4e, 0x6f, 0x62, 0x6c, 0x65,
	0x5c, 0x44, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x5c, 0x56, 0x61, 0x75, 0x6c, 0x74, 0x73, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x19,
	0x4e, 0x6f, 0x62, 0x6c, 0x65, 0x3a, 0x3a, 0x44, 0x6f, 0x6c, 0x6c, 0x61, 0x72, 0x3a, 0x3a, 0x56,
	0x61, 0x75, 0x6c, 0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_noble_dollar_vaults_v1_genesis_proto_rawDescOnce sync.Once
	file_noble_dollar_vaults_v1_genesis_proto_rawDescData = file_noble_dollar_vaults_v1_genesis_proto_rawDesc
)

func file_noble_dollar_vaults_v1_genesis_proto_rawDescGZIP() []byte {
	file_noble_dollar_vaults_v1_genesis_proto_rawDescOnce.Do(func() {
		file_noble_dollar_vaults_v1_genesis_proto_rawDescData = protoimpl.X.CompressGZIP(file_noble_dollar_vaults_v1_genesis_proto_rawDescData)
	})
	return file_noble_dollar_vaults_v1_genesis_proto_rawDescData
}

var file_noble_dollar_vaults_v1_genesis_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_noble_dollar_vaults_v1_genesis_proto_goTypes = []interface{}{
	(*GenesisState)(nil),  // 0: noble.dollar.vaults.v1.GenesisState
	(*PositionEntry)(nil), // 1: noble.dollar.vaults.v1.PositionEntry
	(*Reward)(nil),        // 2: noble.dollar.vaults.v1.Reward
	(PausedType)(0),       // 3: noble.dollar.vaults.v1.PausedType
}
var file_noble_dollar_vaults_v1_genesis_proto_depIdxs = []int32{
	1, // 0: noble.dollar.vaults.v1.GenesisState.positions:type_name -> noble.dollar.vaults.v1.PositionEntry
	2, // 1: noble.dollar.vaults.v1.GenesisState.rewards:type_name -> noble.dollar.vaults.v1.Reward
	3, // 2: noble.dollar.vaults.v1.GenesisState.paused:type_name -> noble.dollar.vaults.v1.PausedType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_noble_dollar_vaults_v1_genesis_proto_init() }
func file_noble_dollar_vaults_v1_genesis_proto_init() {
	if File_noble_dollar_vaults_v1_genesis_proto != nil {
		return
	}
	file_noble_dollar_vaults_v1_vaults_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_noble_dollar_vaults_v1_genesis_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenesisState); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_noble_dollar_vaults_v1_genesis_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_noble_dollar_vaults_v1_genesis_proto_goTypes,
		DependencyIndexes: file_noble_dollar_vaults_v1_genesis_proto_depIdxs,
		MessageInfos:      file_noble_dollar_vaults_v1_genesis_proto_msgTypes,
	}.Build()
	File_noble_dollar_vaults_v1_genesis_proto = out.File
	file_noble_dollar_vaults_v1_genesis_proto_rawDesc = nil
	file_noble_dollar_vaults_v1_genesis_proto_goTypes = nil
	file_noble_dollar_vaults_v1_genesis_proto_depIdxs = nil
}
