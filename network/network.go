package network

import (
	"github.com/swarmbit/spacemesh-state-api/config"

	"github.com/spacemeshos/economics/rewards"
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

func (n *NetworkUtils) GetEpochFirst(epoch uint64) sTypes.LayerID {
	return sTypes.LayerID(sTypes.EpochID(epoch)).Mul(config.LayersPerEpoch)
}

func (n *NetworkUtils) GetNumberOfSlots(weight uint64, totalWeight uint64, epoch uint32) (int32, error) {
	layerSize := n.tortoiseConfig.LayerSize
	minimalWeight := uint64(7_879_129_244)
	if epoch < 8 {
		minimalWeight = 107467138
	}

	slots, err := util.GetNumEligibleSlots(weight, minimalWeight, totalWeight, layerSize, 4032)
	return int32(slots), err
}

func (n *NetworkUtils) FirstEffectiveGenesis() sTypes.LayerID {
	return sTypes.LayerID(config.LayersPerEpoch*2 - 1)
}

func (n *NetworkUtils) GetEpochSubsidy(epoch uint64) uint64 {
	genisesLayer := n.FirstEffectiveGenesis()
	epochFirstLayer := n.GetEpochFirst(epoch)
	var totalEpochSubsidy uint64 = 0
	for i := epochFirstLayer; i < epochFirstLayer+config.LayersPerEpoch; i++ {
		totalEpochSubsidy += rewards.TotalSubsidyAtLayer(i.Difference(genisesLayer))
	}
	return totalEpochSubsidy
}
