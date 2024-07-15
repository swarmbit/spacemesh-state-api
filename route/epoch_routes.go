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
			"error": "Failed to get epoch rewards",
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

func (e *EpochRoutes) GetEpochAtx(c *gin.Context) {
	epochStr := c.Param("epoch")
	epoch, err := strconv.Atoi(epochStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "epoch must be a valid integer",
		})
		return
	}

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

	atxs, errAtx := e.db.GetAtxForEpochPaginated(uint64(epoch-1), int64(offset), int64(limit), sort)
	count, errCount := e.db.CountAtxEpoch(uint64(epoch - 1))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get epoch atx",
		})
		return
	}

	if errAtx != nil || errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch atx for epoch",
		})
	} else if atxs != nil {

		atxResponse := make([]*types.Atx, len(atxs))

		for i, a := range atxs {
			atxResponse[i] = &types.Atx{
				NodeId:            a.NodeID,
				AtxId:             a.AtxID,
				EffectiveNumUnits: a.EffectiveNumUnits,
				Weight:            a.Weight,
				Received:          a.Received,
			}
		}

		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, atxResponse)
	} else {
		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, make([]*types.Atx, 0))
	}

}
