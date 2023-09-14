package types

import "github.com/spacemeshos/go-spacemesh/common/types"

type NodeSmesherReward struct {
	Address     types.Address
	Layer       types.LayerID
	TotalReward uint64
	LayerReward uint64
	Coinbase    types.Address
	AtxID       types.ATXID
	NodeID      types.NodeID
}

type NodeSmesher struct {
	Coinbase          types.Address
	NodeID            types.NodeID
	EffectiveNumUnits int64
}

type NodeLayer struct {
	Layer types.LayerID
}

type NodeAccount struct {
	Address      types.Address
	LayerUpdated types.LayerID
	Balance      int64
	NextNonce    int
	Template     string
	State        []byte
}
