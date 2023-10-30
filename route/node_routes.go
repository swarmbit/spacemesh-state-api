package route

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type NodesRoutes struct {
	db           *database.ReadDB
	networkUtils *network.NetworkUtils
	state        *network.NetworkState
}

func NewNodeRoutes(db *database.ReadDB, networkUtils *network.NetworkUtils, state *network.NetworkState) *NodesRoutes {
	return &NodesRoutes{
		db:           db,
		networkUtils: networkUtils,
		state:        state,
	}
}

func (n *NodesRoutes) GetNode(c *gin.Context) {
	nodeId := c.Param("nodeId")
	node, err := n.db.GetNode(nodeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch node",
		})
		return
	}
	if node.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Not Found",
			"error":  "Node not found",
		})
		return
	}

	c.JSON(200, node)
}

func (n *NodesRoutes) GetNodeRewards(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	sortStr := c.DefaultQuery("sort", "asc")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "offset must be a valid integer",
		})
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "limit must be a valid integer",
		})
		return
	}

	if offset < 0 || limit < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "offset and limit must be greater or equal to 0",
		})
		return
	}

	var sort int8
	if sortStr == "desc" {
		sort = -1
	} else {
		sort = 1
	}

	nodeId := c.Param("nodeId")
	rewards, errRewards := n.db.GetNodeRewards(nodeId, int64(offset), int64(limit), sort)
	count, errCount := n.db.CountNodeRewards(nodeId)

	if errRewards != nil || errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch rewards for node",
		})
	} else if rewards != nil {

		rewardsResponse := make([]*types.Reward, len(rewards))

		for i, v := range rewards {
			rewardsResponse[i] = &types.Reward{
				Rewards: int64(v.TotalReward),
				// legacy
				RewardsDisplay: "",
				Layer:          v.Layer,
				SmesherId:      v.NodeId,
				// legacy
				Time:      "2023-09-05T00:00:00Z",
				Timestamp: config.GenesisEpochSeconds + (v.Layer * config.LayerDuration),
			}
		}

		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, rewardsResponse)
	} else {
		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, make([]*types.Reward, 0))
	}
}

func (n *NodesRoutes) GetNodeRewardsDetails(c *gin.Context) {
	nodeId := c.Param("nodeId")

	networkInfo := n.state.GetInfo()
	epoch := networkInfo.Epoch

	firstLayer := uint32(epoch * config.LayersPerEpoch)
	lastLayer := firstLayer + config.LayersPerEpoch

	countEpochResult, err := n.db.CountNodeRewardsLayers(nodeId, firstLayer, lastLayer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get epoch rewards count",
		})
		return
	}

	sumEpochResult, err := n.db.SumNodeRewardsLayers(nodeId, firstLayer, lastLayer)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get epoch rewards sum",
		})
		return
	}

	total, err := n.db.SumNodeRewardsLayers(nodeId, 0, uint32(networkInfo.Layer))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get epoch rewards sum",
		})
		return
	}

	c.JSON(200, &types.RewardDetails{
		TotalSum:                 int64(total),
		CurrentEpoch:             int64(epoch),
		CurrentEpochRewardsSum:   sumEpochResult,
		CurrentEpochRewardsCount: countEpochResult,
	})
}

func (n *NodesRoutes) GetEligibility(c *gin.Context) {

	networkInfo := n.state.GetInfo()

	nodeId := c.Param("nodeId")

	epoch := networkInfo.Epoch

	nodeAtx, err := n.db.GetAtxWeightNode(nodeId, uint64(epoch-1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get node weight",
		})
		return
	}

	eligibilityCount, err := n.networkUtils.GetNumberOfSlots(uint64(nodeAtx.TotalWeight), networkInfo.TotalWeight, epoch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get eligibility",
		})
		return
	}

	unitReward := networkInfo.EpochSubsidy / networkInfo.TotalWeight
	predictedRewards := unitReward * uint64(nodeAtx.TotalWeight)

	if nodeAtx.TotalWeight == 0 {
		eligibilityCount = -1
		predictedRewards = 0
	}

	c.JSON(200, &types.Eligibility{
		Count:             eligibilityCount,
		EffectiveNumUnits: nodeAtx.TotalEffectiveNumUnits,
		PredictedRewards:  predictedRewards,
	})
}
