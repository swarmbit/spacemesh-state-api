package route

import (
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
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
	epochStr := c.DefaultQuery("epoch", "2")
	epoch, err := strconv.Atoi(epochStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "epoch must be a valid integer",
		})
		return
	}

	count, err := s.state.GetEligibilityCount(nodeID, sTypes.EpochID(epoch))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
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
		Count:            count,
		SmesherId:        nodeIdStr,
		PredictedRewards: predictedRewards * uint64(count),
	})
}
