package network

import (
	"github.com/swarmbit/spacemesh-state-api/config"
    "math/big"

    "github.com/spacemeshos/economics/rewards"
	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/proposals/util"
	"github.com/spacemeshos/go-spacemesh/tortoise"
)

const (
	OneYear       = 105120               // 365 days, in 5-minute intervals
	VestStart     = OneYear              // one year, in layers
	VestEnd       = 4 * OneYear
	OneSmesh = 1000000000 // 1e9 (1bn) smidge per smesh
	TotalVaulted  = OneSmesh * 150000000 // 150mn smesh
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


func (n *NetworkUtils) Vested(layer uint64) uint64 {
	lid := sTypes.LayerID(layer)
	if lid.Before(VestStart) {
		return 0
	}
	if !lid.Before(VestEnd) {
		return TotalVaulted
	}
	vested := new(big.Int).SetUint64(TotalVaulted)
	vested.Mul(vested, new(big.Int).SetUint64(uint64(lid.Difference(VestStart))))
	// Note: VestingStart may equal VestingEnd but division by zero is not possible here since in this case
	// one of the first two conditionals above would have been triggered and the method would already have
	// returned.
	vested.Div(vested, new(big.Int).SetUint64(uint64(VestEnd - VestStart)))
	return vested.Uint64()
}
