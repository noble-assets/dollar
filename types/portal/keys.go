package portal

const SubmoduleName = "dollar/portal"

// NOTE: These variables are initialized when creating the module keeper.
var (
	PaddedTransceiverAddress = make([]byte, 32)
	TransceiverAddress       = ""
	PaddedManagerAddress     = make([]byte, 32)
	ManagerAddress           = ""
	RawToken                 = make([]byte, 32)
)

var (
	OwnerKey   = []byte("owner")
	PeerPrefix = []byte("peer/")
	NonceKey   = []byte("nonce")
)
