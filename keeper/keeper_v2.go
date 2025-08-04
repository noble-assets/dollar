package keeper

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/types/vaults"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// V2 Collections for share-based vault system
// These are added to the existing Keeper struct via embedding or extension

// V2VaultCollections contains all collections for the V2 vault system
type V2VaultCollections struct {
	// V2 Vault State Collections
	VaultStates       collections.Map[int32, vaultsv2.VaultState]                              // vault_type -> VaultState
	UserPositions     collections.Map[collections.Pair[int32, []byte], vaultsv2.UserPosition]  // (vault_type, address) -> UserPosition
	ExitRequests      collections.Map[string, vaultsv2.ExitRequest]                            // request_id -> ExitRequest
	ExitQueue         collections.Map[collections.Pair[int32, uint64], string]                 // (vault_type, queue_pos) -> request_id
	RemotePositions   collections.Map[collections.Pair[string, []byte], vaults.RemotePosition] // (route_id, address) -> RemotePosition
	InFlightPositions collections.Map[uint64, vaults.InFlightPosition]                         // nonce -> InFlightPosition

	// Configuration Collections
	NAVConfigs       collections.Map[int32, vaultsv2.NAVConfig]      // vault_type -> NAVConfig
	FeeConfigs       collections.Map[int32, vaultsv2.FeeConfig]      // vault_type -> FeeConfig
	CrossChainRoutes collections.Map[string, vaults.CrossChainRoute] // route_id -> CrossChainRoute

	// Fee Management Collections
	FeeAccruals    collections.Map[collections.Triple[int32, int32, int64], vaults.FeeAccrual] // (vault_type, fee_type, period) -> FeeAccrual
	LastFeeAccrual collections.Map[int32, int64]                                               // vault_type -> timestamp
	FeeCollections collections.Map[collections.Pair[int32, int64], vaults.FeeCollection]       // (vault_type, timestamp) -> FeeCollection

	// NAV and Pricing Collections
	NAVUpdates          collections.Map[collections.Pair[int32, int64], vaults.NAVUpdate]                  // (vault_type, timestamp) -> NAVUpdate
	LastNAVUpdate       collections.Map[int32, int64]                                                      // vault_type -> timestamp
	CrossChainSnapshots collections.Map[collections.Pair[int32, int64], vaults.CrossChainPositionSnapshot] // (vault_type, timestamp) -> Snapshot

	// Operational Collections
	DriftAlerts collections.Map[collections.Pair[string, []byte], vaults.DriftAlert] // (route_id, address) -> DriftAlert
	LossEvents  collections.Map[collections.Pair[int32, int64], vaults.LossEvent]    // (vault_type, timestamp) -> LossEvent
}

// InitializeV2Collections initializes all V2 collections in the keeper
func (k *Keeper) InitializeV2Collections(builder *collections.SchemaBuilder) V2VaultCollections {
	return V2VaultCollections{
		// V2 Vault State Collections
		VaultStates: collections.NewMap(
			builder,
			collections.NewPrefix(200),
			"v2_vault_states",
			collections.Int32Key,
			codec.CollValue[vaultsv2.VaultState](k.cdc),
		),
		UserPositions: collections.NewMap(
			builder,
			collections.NewPrefix(201),
			"v2_user_positions",
			collections.PairKeyCodec(collections.Int32Key, collections.BytesKey),
			codec.CollValue[vaultsv2.UserPosition](k.cdc),
		),
		ExitRequests: collections.NewMap(
			builder,
			collections.NewPrefix(202),
			"v2_exit_requests",
			collections.StringKey,
			codec.CollValue[vaultsv2.ExitRequest](k.cdc),
		),
		ExitQueue: collections.NewMap(
			builder,
			collections.NewPrefix(203),
			"v2_exit_queue",
			collections.PairKeyCodec(collections.Int32Key, collections.Uint64Key),
			collections.StringValue,
		),
		RemotePositions: collections.NewMap(
			builder,
			collections.NewPrefix(204),
			"v2_remote_positions",
			collections.PairKeyCodec(collections.StringKey, collections.BytesKey),
			codec.CollValue[vaults.RemotePosition](k.cdc),
		),
		InFlightPositions: collections.NewMap(
			builder,
			collections.NewPrefix(205),
			"v2_inflight_positions",
			collections.Uint64Key,
			codec.CollValue[vaults.InFlightPosition](k.cdc),
		),

		// Configuration Collections
		NAVConfigs: collections.NewMap(
			builder,
			collections.NewPrefix(220),
			"v2_nav_configs",
			collections.Int32Key,
			codec.CollValue[vaultsv2.NAVConfig](k.cdc),
		),
		FeeConfigs: collections.NewMap(
			builder,
			collections.NewPrefix(221),
			"v2_fee_configs",
			collections.Int32Key,
			codec.CollValue[vaultsv2.FeeConfig](k.cdc),
		),
		CrossChainRoutes: collections.NewMap(
			builder,
			collections.NewPrefix(222),
			"v2_cross_chain_routes",
			collections.StringKey,
			codec.CollValue[vaults.CrossChainRoute](k.cdc),
		),

		// Fee Management Collections
		FeeAccruals: collections.NewMap(
			builder,
			collections.NewPrefix(230),
			"v2_fee_accruals",
			collections.TripleKeyCodec(collections.Int32Key, collections.Int32Key, collections.Int64Key),
			codec.CollValue[vaults.FeeAccrual](k.cdc),
		),
		LastFeeAccrual: collections.NewMap(
			builder,
			collections.NewPrefix(231),
			"v2_last_fee_accrual",
			collections.Int32Key,
			collections.Int64Value,
		),
		FeeCollections: collections.NewMap(
			builder,
			collections.NewPrefix(232),
			"v2_fee_collections",
			collections.PairKeyCodec(collections.Int32Key, collections.Int64Key),
			codec.CollValue[vaults.FeeCollection](k.cdc),
		),

		// NAV and Pricing Collections
		NAVUpdates: collections.NewMap(
			builder,
			collections.NewPrefix(240),
			"v2_nav_updates",
			collections.PairKeyCodec(collections.Int32Key, collections.Int64Key),
			codec.CollValue[vaults.NAVUpdate](k.cdc),
		),
		LastNAVUpdate: collections.NewMap(
			builder,
			collections.NewPrefix(241),
			"v2_last_nav_update",
			collections.Int32Key,
			collections.Int64Value,
		),
		CrossChainSnapshots: collections.NewMap(
			builder,
			collections.NewPrefix(242),
			"v2_cross_chain_snapshots",
			collections.PairKeyCodec(collections.Int32Key, collections.Int64Key),
			codec.CollValue[vaults.CrossChainPositionSnapshot](k.cdc),
		),

		// Operational Collections
		DriftAlerts: collections.NewMap(
			builder,
			collections.NewPrefix(251),
			"v2_drift_alerts",
			collections.PairKeyCodec(collections.StringKey, collections.BytesKey),
			codec.CollValue[vaults.DriftAlert](k.cdc),
		),
		LossEvents: collections.NewMap(
			builder,
			collections.NewPrefix(252),
			"v2_loss_events",
			collections.PairKeyCodec(collections.Int32Key, collections.Int64Key),
			codec.CollValue[vaults.LossEvent](k.cdc),
		),
	}
}

// Helper functions for creating composite keys

// V2VaultUserKey creates a key for vault-user combinations
func V2VaultUserKey(vaultType vaults.VaultType, address sdk.AccAddress) collections.Pair[int32, []byte] {
	return collections.Join(int32(vaultType), address.Bytes())
}

// V2ExitQueueKey creates a key for exit queue positions
func V2ExitQueueKey(vaultType vaults.VaultType, queuePosition uint64) collections.Pair[int32, uint64] {
	return collections.Join(int32(vaultType), queuePosition)
}

// V2RemotePositionKey creates a key for remote positions
func V2RemotePositionKey(routeID string, address sdk.AccAddress) collections.Pair[string, []byte] {
	return collections.Join(routeID, address.Bytes())
}

// V2FeeAccrualKey creates a key for fee accruals
func V2FeeAccrualKey(vaultType vaults.VaultType, feeType vaults.FeeType, period int64) collections.Triple[int32, int32, int64] {
	return collections.Join3(int32(vaultType), int32(feeType), period)
}

// V2NAVHistoryKey creates a key for NAV history
func V2NAVHistoryKey(vaultType vaults.VaultType, timestamp int64) collections.Pair[int32, int64] {
	return collections.Join(int32(vaultType), timestamp)
}

// V2CrossChainSnapshotKey creates a key for cross-chain snapshots
func V2CrossChainSnapshotKey(vaultType vaults.VaultType, timestamp int64) collections.Pair[int32, int64] {
	return collections.Join(int32(vaultType), timestamp)
}

// V2FeeCollectionKey creates a key for fee collections
func V2FeeCollectionKey(vaultType vaults.VaultType, timestamp int64) collections.Pair[int32, int64] {
	return collections.Join(int32(vaultType), timestamp)
}

// V2LossEventKey creates a key for loss events
func V2LossEventKey(vaultType vaults.VaultType, timestamp int64) collections.Pair[int32, int64] {
	return collections.Join(int32(vaultType), timestamp)
}

// V2DriftAlertKey creates a key for drift alerts
func V2DriftAlertKey(routeID string, address sdk.AccAddress) collections.Pair[string, []byte] {
	return collections.Join(routeID, address.Bytes())
}
