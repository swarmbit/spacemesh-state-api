package network

import (
	"github.com/swarmbit/spacemesh-state-api/config"

	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
)

func GetEpoch(layer uint64) sTypes.EpochID {
	return sTypes.EpochID(layer / config.LayersPerEpoch)
}
