package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
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
