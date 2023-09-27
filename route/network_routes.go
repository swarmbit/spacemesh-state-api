package route

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
	"github.com/swarmbit/spacemesh-state-api/types"
)

const INFO_KEY = "info"

type NetworkRoutes struct {
	db          *database.ReadDB
	networkInfo *sync.Map
}

func NewNetworkRoutes(db *database.ReadDB) *NetworkRoutes {
	routes := &NetworkRoutes{
		db:          db,
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

	layer, err := n.db.GetLastProcessedLayer()
	if err != nil {
		fmt.Printf("Failed to get last processed layer: %s", err.Error())
		return
	}

	epoch := network.GetEpoch(uint64(layer.Layer))

	atxEpoch, err := n.db.CountAtxEpoch(uint64(epoch - 1))
	if err != nil {
		fmt.Printf("Failed to count atx epoch: %s", err.Error())
		return
	}

	atxNextEpoch, err := n.db.CountAtxEpoch(uint64(epoch))
	if err != nil {
		fmt.Printf("Failed to count next atx epoch: %s", err.Error())
		return
	}

	totalAccounts, err := n.db.CountAccounts()
	if err != nil {
		fmt.Printf("Failed to count accounts: %s", err.Error())
		return
	}

	networkInfo, err := n.db.GetNetworkInfo()
	if err != nil {
		fmt.Printf("Failed to get network info: %s", err.Error())
		return
	}

	atxEpochTotals, err := n.db.GetAtxEpoch(uint64(epoch - 1))
	if err != nil {
		fmt.Printf("Failed to get epoch totals: %s", err.Error())
		return
	}

	atxNextEpochTotals, err := n.db.GetAtxEpoch(uint64(epoch))
	if err != nil {
		fmt.Printf("Failed to get next epoch totals: %s", err.Error())
		return
	}

	n.networkInfo.Store(INFO_KEY, &types.NetworkInfo{
		Epoch:                  epoch.Uint32(),
		Layer:                  uint64(layer.Layer),
		EffectiveUnitsCommited: atxEpochTotals.TotalEffectiveNumUnits,
		CirculatingSupply:      networkInfo.CirculatingSupply,
		TotalAccounts:          uint64(totalAccounts),
		//AtxHex:                 hex.EncodeToString(highest.Bytes()),
		//AtxBase64:              base64.StdEncoding.EncodeToString(highest.Bytes()),
		TotalActiveSmeshers: uint64(atxEpoch),
		NextEpoch: &types.NetworkInfoNextEpoch{
			Epoch:                  epoch.Uint32() + 1,
			EffectiveUnitsCommited: int64(atxNextEpochTotals.TotalEffectiveNumUnits),
			TotalActiveSmeshers:    atxNextEpoch,
		},
	})
}
