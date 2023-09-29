package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
)

type NodesRoutes struct {
	db           *database.ReadDB
	networkUtils *network.NetworkUtils
}

func NewNodeRoutes(db *database.ReadDB, networkUtils *network.NetworkUtils) *NodesRoutes {
	return &NodesRoutes{
		db:           db,
		networkUtils: networkUtils,
	}
}

func (n *NodesRoutes) GetNode(c *gin.Context) {

}

func (n *NodesRoutes) GetEligibility(c *gin.Context) {
	/*
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
	*/
}
