package types

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Index: 1e12,
	}
}

func (genesis *GenesisState) Validate() error {
	return nil
}
