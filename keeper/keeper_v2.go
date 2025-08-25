package keeper

import (
	"cosmossdk.io/collections"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"dollar.noble.xyz/v2/keeper/crosschain"
	vaults "dollar.noble.xyz/v2/types/vaults"
	vaultsv2 "dollar.noble.xyz/v2/types/vaults/v2"
)

// V2 Collections for share-based vault system
// These are added to the existing Keeper struct via embedding or extension

// V2VaultCollections contains all collections for the V2 vault system
type V2VaultCollections struct {
	// V2 Vault State Collections
	VaultStates          collections.Item[vaultsv2.VaultState]                                      // Single vault state
	UserPositions        collections.Map[[]byte, vaultsv2.UserPosition]                             // address -> UserPosition
	ExitRequests         collections.Map[string, vaultsv2.ExitRequest]                              // request_id -> ExitRequest
	ExitQueue            collections.Map[uint64, string]                                            // queue_pos -> request_id
	RemotePositions      collections.Map[collections.Pair[string, []byte], vaultsv2.RemotePosition] // (route_id, address) -> RemotePosition
	InFlightPositions    collections.Map[uint64, vaultsv2.InFlightPosition]                         // nonce -> InFlightPosition
	InflightValueByRoute collections.Map[string, string]                                            // route_id -> inflight value
	TotalInflightValue   collections.Item[string]                                                   // total inflight value

	// Configuration Collections
	NAVConfig        collections.Item[vaultsv2.NAVConfig]              // Single NAV config
	FeeConfig        collections.Item[vaultsv2.FeeConfig]              // Single fee config
	CrossChainRoutes collections.Map[string, vaultsv2.CrossChainRoute] // route_id -> CrossChainRoute

	// Fee Management Collections
	FeeAccruals    collections.Map[collections.Pair[int32, int64], vaultsv2.FeeAccrual] // (fee_type, period) -> FeeAccrual
	LastFeeAccrual collections.Item[int64]                                              // timestamp
	FeeCollections collections.Map[int64, vaultsv2.FeeCollection]                       // timestamp -> FeeCollection

	// NAV and Pricing Collections
	NAVUpdates          collections.Map[int64, vaultsv2.NAVUpdate]                  // timestamp -> NAVUpdate
	LastNAVUpdate       collections.Item[int64]                                     // timestamp
	CrossChainSnapshots collections.Map[int64, vaultsv2.CrossChainPositionSnapshot] // timestamp -> Snapshot

	// Operational Collections
	DriftAlerts collections.Map[collections.Pair[string, []byte], vaultsv2.DriftAlert] // (route_id, address) -> DriftAlert
	LossEvents  collections.Map[int64, vaultsv2.LossEvent]                             // timestamp -> LossEvent

	// Cross-chain configuration
	CrossChainConfig collections.Item[vaultsv2.CrossChainConfig] // Global cross-chain configuration

	// Cross-chain keeper
	CrossChainStore *crosschain.CrossChainKeeper
}

// InitializeV2Collections initializes all V2 collections in the keeper
func (k *Keeper) InitializeV2Collections(builder *collections.SchemaBuilder) V2VaultCollections {
	v2Collections := V2VaultCollections{
		// V2 Vault State Collections
		VaultStates: collections.NewItem(
			builder,
			collections.NewPrefix(200),
			"v2_vault_state",
			codec.CollValue[vaultsv2.VaultState](k.cdc),
		),
		UserPositions: collections.NewMap(
			builder,
			collections.NewPrefix(201),
			"v2_user_positions",
			collections.BytesKey,
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
			collections.Uint64Key,
			collections.StringValue,
		),
		RemotePositions: collections.NewMap(
			builder,
			collections.NewPrefix(204),
			"v2_remote_positions",
			collections.PairKeyCodec(collections.StringKey, collections.BytesKey),
			codec.CollValue[vaultsv2.RemotePosition](k.cdc),
		),
		InFlightPositions: collections.NewMap(
			builder,
			collections.NewPrefix(205),
			"v2_inflight_positions",
			collections.Uint64Key,
			codec.CollValue[vaultsv2.InFlightPosition](k.cdc),
		),
		InflightValueByRoute: collections.NewMap(
			builder,
			collections.NewPrefix(206),
			"v2_inflight_value_by_route",
			collections.StringKey,
			collections.StringValue,
		),
		TotalInflightValue: collections.NewItem(
			builder,
			collections.NewPrefix(207),
			"v2_total_inflight_value",
			collections.StringValue,
		),

		// Configuration Collections
		NAVConfig: collections.NewItem(
			builder,
			collections.NewPrefix(220),
			"v2_nav_config",
			codec.CollValue[vaultsv2.NAVConfig](k.cdc),
		),
		FeeConfig: collections.NewItem(
			builder,
			collections.NewPrefix(221),
			"v2_fee_config",
			codec.CollValue[vaultsv2.FeeConfig](k.cdc),
		),
		CrossChainRoutes: collections.NewMap(
			builder,
			collections.NewPrefix(222),
			"v2_cross_chain_routes",
			collections.StringKey,
			codec.CollValue[vaultsv2.CrossChainRoute](k.cdc),
		),

		// Fee Management Collections
		FeeAccruals: collections.NewMap(
			builder,
			collections.NewPrefix(230),
			"v2_fee_accruals",
			collections.PairKeyCodec(collections.Int32Key, collections.Int64Key),
			codec.CollValue[vaultsv2.FeeAccrual](k.cdc),
		),
		LastFeeAccrual: collections.NewItem(
			builder,
			collections.NewPrefix(231),
			"v2_last_fee_accrual",
			collections.Int64Value,
		),
		FeeCollections: collections.NewMap(
			builder,
			collections.NewPrefix(232),
			"v2_fee_collections",
			collections.Int64Key,
			codec.CollValue[vaultsv2.FeeCollection](k.cdc),
		),

		// NAV and Pricing Collections
		NAVUpdates: collections.NewMap(
			builder,
			collections.NewPrefix(240),
			"v2_nav_updates",
			collections.Int64Key,
			codec.CollValue[vaultsv2.NAVUpdate](k.cdc),
		),
		LastNAVUpdate: collections.NewItem(
			builder,
			collections.NewPrefix(241),
			"v2_last_nav_update",
			collections.Int64Value,
		),
		CrossChainSnapshots: collections.NewMap(
			builder,
			collections.NewPrefix(242),
			"v2_cross_chain_snapshots",
			collections.Int64Key,
			codec.CollValue[vaultsv2.CrossChainPositionSnapshot](k.cdc),
		),

		// Operational Collections
		DriftAlerts: collections.NewMap(
			builder,
			collections.NewPrefix(251),
			"v2_drift_alerts",
			collections.PairKeyCodec(collections.StringKey, collections.BytesKey),
			codec.CollValue[vaultsv2.DriftAlert](k.cdc),
		),
		LossEvents: collections.NewMap(
			builder,
			collections.NewPrefix(252),
			"v2_loss_events",
			collections.Int64Key,
			codec.CollValue[vaultsv2.LossEvent](k.cdc),
		),
		CrossChainConfig: collections.NewItem(
			builder,
			collections.NewPrefix(253),
			"v2_cross_chain_config",
			codec.CollValue[vaultsv2.CrossChainConfig](k.cdc),
		),
	}

	// TODO: Initialize cross-chain keeper after updating collection types
	// v2Collections.CrossChainStore = crosschain.NewCrossChainKeeper(
	//	v2Collections.CrossChainRoutes,
	//	v2Collections.RemotePositions,
	//	v2Collections.InFlightPositions,
	//	v2Collections.CrossChainSnapshots,
	//	v2Collections.DriftAlerts,
	//	v2Collections.CrossChainConfig,
	//	nonceCounter,
	// )

	return v2Collections
}

// Helper functions for creating composite keys

// V2VaultUserKey creates a key for user positions (single vault)
func V2VaultUserKey(address sdk.AccAddress) []byte {
	return address.Bytes()
}

// V2ExitQueueKey creates a key for exit queue positions (single vault)
func V2ExitQueueKey(queuePosition uint64) uint64 {
	return queuePosition
}

// V2RemotePositionKey creates a key for remote positions
func V2RemotePositionKey(routeID string, address sdk.AccAddress) collections.Pair[string, []byte] {
	return collections.Join(routeID, address.Bytes())
}

// V2FeeAccrualKey creates a key for fee accruals
func V2FeeAccrualKey(feeType vaults.FeeType, period int64) collections.Pair[int32, int64] {
	return collections.Join(int32(feeType), period)
}

// V2NAVHistoryKey creates a key for NAV history
func V2NAVHistoryKey(timestamp int64) int64 {
	return timestamp
}

// V2CrossChainSnapshotKey creates a key for cross-chain snapshots
func V2CrossChainSnapshotKey(timestamp int64) int64 {
	return timestamp
}

// RegisterCrossChainProvider registers a cross-chain provider
func (collections *V2VaultCollections) RegisterCrossChainProvider(provider crosschain.CrossChainProvider) {
	if collections.CrossChainStore != nil {
		collections.CrossChainStore.RegisterProvider(provider)
	}
}

// GetCrossChainKeeper returns the cross-chain keeper
func (collections *V2VaultCollections) GetCrossChainKeeper() *crosschain.CrossChainKeeper {
	return collections.CrossChainStore
}

// V2FeeCollectionKey creates a key for fee collections
func V2FeeCollectionKey(timestamp int64) int64 {
	return timestamp
}

// V2LossEventKey creates a key for loss events
func V2LossEventKey(timestamp int64) int64 {
	return timestamp
}

// V2DriftAlertKey creates a key for drift alerts
func V2DriftAlertKey(routeID string, address sdk.AccAddress) collections.Pair[string, []byte] {
	return collections.Join(routeID, address.Bytes())
}
