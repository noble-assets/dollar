package vaults

import (
	"fmt"
	"strings"

	vaultsv1 "dollar.noble.xyz/api/vaults/v1"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const SubmoduleName = "dollar/vaults"

var (
	StakedVaultName    = fmt.Sprintf("%s/%s", SubmoduleName, strings.ToLower(vaultsv1.VaultType_STAKED.String()))
	StakedVaultAddress = authtypes.NewModuleAddress(StakedVaultName)

	FlexibleVaultName    = fmt.Sprintf("%s/%s", SubmoduleName, strings.ToLower(vaultsv1.VaultType_FLEXIBLE.String()))
	FlexibleVaultAddress = authtypes.NewModuleAddress(FlexibleVaultName)
)

var (
	PausedKey                 = []byte("paused")
	TotalFlexiblePrincipalKey = []byte("total_flexible_principal")
	PositionPrefix            = []byte("position/")
	RewardsPrefix             = []byte("rewards/")
)
