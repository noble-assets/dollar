package types

import (
	"dollar.noble.xyz/types/portal"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	portal.RegisterLegacyAminoCodec(cdc)

	cdc.RegisterConcrete(&MsgClaimYield{}, "dollar/ClaimYield", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	portal.RegisterInterfaces(registry)

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgClaimYield{})

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
