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

type EpochRoutes struct {
	db           *database.ReadDB
	networkUtils *network.NetworkUtils
	state        *network.NetworkState
}

func NewEpochRoutes(db *database.ReadDB, networkUtils *network.NetworkUtils, state *network.NetworkState) *EpochRoutes {
	routes := &EpochRoutes{
		db:           db,
		networkUtils: networkUtils,
		state:        state,
	}
	return routes
}

func (e *EpochRoutes) GetEpoch(c *gin.Context) {

	epochStr := c.Param("epoch")
	epoch, err := strconv.Atoi(epochStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "epoch must be a valid integer",
		})
		return
	}

	atxEpoch, err := e.db.CountAtxEpoch(uint64(epoch - 1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to count atx for epoch",
		})
		return
	}

	atxEpochTotals, err := e.db.GetAtxEpoch(uint64(epoch - 1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get atx for epoch",
		})
		return
	}

	firstLayer := uint32(epoch * config.LayersPerEpoch)
	lastLayer := firstLayer + config.LayersPerEpoch

	rewardsTotal, err := e.db.SumRewardsLayers("", firstLayer, lastLayer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to epoch rewards",
		})
		return
	}
	c.JSON(200, &types.Epoch{
		EffectiveUnitsCommited: atxEpochTotals.TotalEffectiveNumUnits,
		EpochSubsidy:           e.state.GetEpochSubsidy(uint32(epoch)),
		TotalWeight:            atxEpochTotals.TotalWeight,
		TotalRewards:           rewardsTotal,
		TotalActiveSmeshers:    uint64(atxEpoch),
	})
}
