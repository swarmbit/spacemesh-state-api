package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/network"
)

type NetworkRoutes struct {
	state *network.NetworkState
}

func NewNetworkRoutes(state *network.NetworkState) *NetworkRoutes {
	routes := &NetworkRoutes{
		state: state,
	}
	return routes
}

func (n *NetworkRoutes) GetInfo(c *gin.Context) {
	c.JSON(200, n.state.GetInfo())
}
