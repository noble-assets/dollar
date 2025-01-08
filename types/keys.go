package types

import (
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const ModuleName = "dollar"

var (
	ModuleAddress = authtypes.NewModuleAddress(ModuleName)

	YieldName    = fmt.Sprintf("%s/yield", ModuleName)
	YieldAddress = authtypes.NewModuleAddress(YieldName)
)

var (
	IndexKey        = []byte("index")
	PrincipalPrefix = []byte("principal/")
)
