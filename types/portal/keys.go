package portal

import (
	"fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

const SubmoduleName = "dollar/portal"

var (
	TransceiverAddress = make([]byte, 32)
	ManagerAddress     = make([]byte, 32)
)

var (
	OwnerKey   = []byte("owner")
	PeerPrefix = []byte("peer/")
)

func init() {
	transceiverAddress := authtypes.NewModuleAddress(fmt.Sprintf("%s/transceiver", SubmoduleName))
	copy(TransceiverAddress[12:], transceiverAddress)

	managerAddress := authtypes.NewModuleAddress(fmt.Sprintf("%s/manager", SubmoduleName))
	copy(ManagerAddress[12:], managerAddress)
}
