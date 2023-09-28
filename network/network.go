package network

import (
	"github.com/swarmbit/spacemesh-state-api/config"

	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/proposals/util"
	"github.com/spacemeshos/go-spacemesh/tortoise"
)

type NetworkUtils struct {
	tortoiseConfig tortoise.Config
}

func NewNetworkUtils() *NetworkUtils {
	return &NetworkUtils{
		tortoiseConfig: tortoise.DefaultConfig(),
	}
}

func (n *NetworkUtils) GetEpoch(layer uint64) sTypes.EpochID {
	return sTypes.EpochID(layer / config.LayersPerEpoch)
}

func (n *NetworkUtils) GetNumberOfSlots(weight uint64, totalWeight uint64) (int32, error) {
	layerSize := n.tortoiseConfig.LayerSize
	minimalWeight := n.tortoiseConfig.MinimalActiveSetWeight

	slots, err := util.GetNumEligibleSlots(weight, minimalWeight, totalWeight, layerSize, 4032)
	return int32(slots), err
}
