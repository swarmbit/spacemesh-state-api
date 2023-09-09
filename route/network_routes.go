package route

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/swarmbit/spacemesh-state-api/state"
	"github.com/swarmbit/spacemesh-state-api/types"
)

const INFO_KEY = "info"

type NetworkRoutes struct {
	state       *state.State
	db          sql.Executor
	networkInfo *sync.Map
}

func NewNetworkRoutes(state *state.State) *NetworkRoutes {
	routes := &NetworkRoutes{
		state:       state,
		db:          state.DB,
		networkInfo: &sync.Map{},
	}
	routes.fetchNetworkInfo()
	routes.periodicNetworkInfoFetch()
	return routes
}

func (n *NetworkRoutes) GetInfo(c *gin.Context) {
	networkInfo, exists := n.networkInfo.Load(INFO_KEY)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch network info",
		})
		return
	}
	c.JSON(200, networkInfo)
}

func (n *NetworkRoutes) periodicNetworkInfoFetch() {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for range ticker.C {
			n.fetchNetworkInfo()
		}
	}()
}

func (n *NetworkRoutes) fetchNetworkInfo() {
	epoch, layer, err := n.state.GetCurrentEpoch()
	if err != nil {
		return
	}
	highest, err := n.state.GetHighestAtx()
	if err != nil {
		return
	}
	effectiveUnitsCommited, err := n.state.GetTotalCommittedForEpoch(epoch - 1)
	if err != nil {
		return
	}
	circulatingSupply, err := n.state.GetCirculatingSupply()
	if err != nil {
		return
	}
	totalAccounts, err := n.state.GetNumberOfAccounts()
	if err != nil {
		return
	}
	totalActiveSmeshers, err := n.state.GetActiveAtxCount(epoch - 1)
	if err != nil {
		return
	}

	n.networkInfo.Store(INFO_KEY, &types.NetworkInfo{
		Epoch:                  epoch.Uint32(),
		Layer:                  int64(layer),
		EffectiveUnitsCommited: effectiveUnitsCommited,
		CirculatingSupply:      circulatingSupply,
		TotalAccounts:          totalAccounts,
		AtxHex:                 hex.EncodeToString(highest.Bytes()),
		AtxBase64:              base64.StdEncoding.EncodeToString(highest.Bytes()),
		TotalActiveSmeshers:    totalActiveSmeshers,
	})
}
