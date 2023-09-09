package route

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql/atxs"
	"github.com/swarmbit/spacemesh-state-api/state"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type SmesherRoutes struct {
	state *state.State
}

func NewSmesherRoutes(state *state.State) *SmesherRoutes {
	return &SmesherRoutes{
		state: state,
	}
}

func (s *SmesherRoutes) GetSmesherEligibility(c *gin.Context) {
	nodeIdStr := c.Param("smesherId")
	nodeBytes, err := hex.DecodeString(nodeIdStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed parse smesher id",
		})
	}

	nodeID := sTypes.BytesToNodeID(nodeBytes)

	count, err := s.state.GetEligibilityCount(nodeID)
	if err != nil && count == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Smesher not found for current epoch",
			"error":  "Failed get eligibility",
		})
		return
	}

	currentEpoch, _, err := s.state.GetCurrentEpoch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get current epoch",
		})
		return
	}

	atx, err := atxs.GetByEpochAndNodeID(s.state.DB, currentEpoch-1, nodeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get arx",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal server error",
			"error":  "Failed get eligibility",
		})
		return
	}

	predictedRewards, err := s.state.GetLatestAVGLayerReward()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed get predicted rewards",
		})
		return
	}
	c.JSON(200, &types.Eligibility{
		Count:             count,
		SmesherId:         nodeIdStr,
		PredictedRewards:  predictedRewards * uint64(count),
		EffectiveNumUnits: int64(atx.EffectiveNumUnits()),
		Address:           atx.Coinbase.String(),
	})
}
