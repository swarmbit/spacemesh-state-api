package types

import "github.com/spacemeshos/go-spacemesh/common/types"

type SmesherReward struct {
	Layer       types.LayerID
	TotalReward uint64
	LayerReward uint64
	Coinbase    types.Address
	AtxID       types.ATXID
	NodeID      types.NodeID
}
