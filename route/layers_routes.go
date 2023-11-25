package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
	"github.com/swarmbit/spacemesh-state-api/types"
	"net/http"
	"strconv"
)

type LayersRoutes struct {
	db           *database.ReadDB
	networkUtils *network.NetworkUtils
	state        *network.NetworkState
}

func NewLayersRoutes(db *database.ReadDB, networkUtils *network.NetworkUtils, state *network.NetworkState) *LayersRoutes {
	routes := &LayersRoutes{
		db:           db,
		networkUtils: networkUtils,
		state:        state,
	}
	return routes
}

func (l *LayersRoutes) GetLayers(c *gin.Context) {

	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	sortStr := c.DefaultQuery("sort", "desc")

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
	if sortStr == "asc" {
		sort = 1
	} else {
		sort = -1
	}

	layers, err := l.db.GetProcessedsLayers(int64(offset), int64(limit), sort)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get layers",
		})
		return
	}

	layersInt := make([]int64, len(layers))

	for i, layer := range layers {
		layersInt[i] = layer.Layer
	}

	c.JSON(200, layersInt)
}

func (l *LayersRoutes) GetLayerTransactions(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	sortStr := c.DefaultQuery("sort", "asc")
	completeStr := c.DefaultQuery("complete", "true")

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

	complete := completeStr == "true"

	layerStr := c.Param("layer")

	layer, err := strconv.Atoi(layerStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "layer must be a valid integer",
		})
		return
	}

	transactions, errRewards := l.db.GetLayerTransactions(layer, int64(offset), int64(limit), sort, complete)
	count, errCount := l.db.CountLayerTransactions(layer)

	if errRewards != nil || errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch transactions for layer",
		})
	} else if transactions != nil {

		transactionsResponse := make([]*types.Transaction, len(transactions))

		for i, v := range transactions {
			method := ""
			if v.Method == 0 {
				method = "Spawn"
			}
			if v.Method == 16 {
				method = "Spend"
			}
			transactionsResponse[i] = &types.Transaction{
				ID:               v.ID,
				Status:           v.Status,
				PrincipalAccount: v.PrincipaAccount,
				ReceiverAccount:  v.ReceiverAccount,
				Fee:              v.Gas * v.GasPrice,
				Amount:           v.Amount,
				Layer:            v.Layer,
				Counter:          v.Counter,
				Method:           method,
				Timestamp:        int64(config.GenesisEpochSeconds + (v.Layer * config.LayerDuration)),
			}
		}

		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, transactionsResponse)
	} else {
		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, make([]*types.Transaction, 0))
	}
}

func (l *LayersRoutes) GetLayerRewards(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")
	sortStr := c.DefaultQuery("sort", "desc")

	layerStr := c.Param("layer")

	layer, err := strconv.Atoi(layerStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "layer must be a valid integer",
			})
		return
	}

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
	if sortStr == "asc" {
		sort = 1
	} else {
		sort = -1
	}

	rewards, errRewards := l.db.GetLayerRewards(layer, int64(offset), int64(limit), sort)
	count, errCount := l.db.CountLayerRewards(layer)

	if errRewards != nil || errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch rewards for account",
			})
	} else if rewards != nil {

		rewardsResponse := make([]*types.Reward, len(rewards))

		for i, v := range rewards {
			rewardsResponse[i] = &types.Reward{
				Account: v.Coinbase,
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
