package v2

import (
	"time"

	"cosmossdk.io/math"
)

// InflightType enumerates types of funds in transit between the Noble
// vault and remote positions.
type InflightType int32

const (
	// Funds moving from the vault to a remote position
	InflightDepositToPosition InflightType = iota
	// Funds returning from a remote position back to the vault
	InflightWithdrawalFromPosition
	// Funds moving between two remote positions
	InflightRebalanceBetweenPositions
	// Deposits received but not yet deployed to any position
	InflightPendingDeployment
	// Funds returned from positions but not yet distributed to claimants
	InflightPendingWithdrawalDistribution
)

// InflightFund tracks USDN amounts that are currently in transit. These
// values are included in NAV calculations while they are moving between
// the vault and remote positions.
// It captures identifiers, timing and status for each cross-chain transfer.
type InflightFund struct {
	// Identifier of the Hyperlane route these funds are traversing.
	HyperlaneRouteID uint32
	// Provider specific transaction identifier.
	TransactionID string
	// Classification of the inflight movement.
	Type InflightType
	// Amount of USDN in transit.
	Amount math.Int
	// Source Hyperlane domain identifier.
	SourceDomain uint32
	// Destination Hyperlane domain identifier.
	DestDomain uint32
	// Optional position identifiers for rebalancing operations.
	SourcePosition *uint64
	DestPosition   *uint64
	// Block time when the transfer was initiated.
	InitiatedAt time.Time
	// Expected completion time for the transfer.
	ExpectedAt time.Time
	// Current status of the inflight funds.
	Status InFlightStatus
	// Value of the funds at the time the transfer was initiated.
	ValueAtInitiation math.Int
}
