package types

import "cosmossdk.io/math"

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Index: math.LegacyOneDec(),
	}
}

func (genesis *GenesisState) Validate() error {
	return nil
}
