package route

import (
	"encoding/base64"
	"encoding/hex"
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

	hexAtx, err := n.getHigestAtx(uint64(epoch - 1))
	if err != nil {
		fmt.Printf("Failed to get highest atx: %s", err.Error())
		return
	}

	base64Atx, err := hexToBase64(hexAtx)
	if err != nil {
		fmt.Printf("Failed to get base64 atx: %s", err.Error())
		return
	}

	var genisesAccounts int64 = 28

	n.networkInfo.Store(INFO_KEY, &types.NetworkInfo{
		Epoch:                  epoch.Uint32(),
		Layer:                  uint64(layer.Layer),
		EffectiveUnitsCommited: atxEpochTotals.TotalEffectiveNumUnits,
		CirculatingSupply:      networkInfo.CirculatingSupply,
		TotalAccounts:          uint64(totalAccounts + genisesAccounts),
		AtxHex:                 hexAtx,
		AtxBase64:              base64Atx,
		TotalActiveSmeshers:    uint64(atxEpoch),
		NextEpoch: &types.NetworkInfoNextEpoch{
			Epoch:                  epoch.Uint32() + 1,
			EffectiveUnitsCommited: int64(atxNextEpochTotals.TotalEffectiveNumUnits),
			TotalActiveSmeshers:    atxNextEpoch,
		},
	})

}

func (n *NetworkRoutes) getHigestAtx(epoch uint64) (string, error) {
	atxs, err := n.db.GetAtxForEpoch(epoch)
	if err != nil {
		return "", err
	}

	malfeasanceNodes, err := n.db.GetMalfeasanceNodes()
	if err != nil {
		return "", err
	}

	malfeasanceNodesMap := make(map[string]bool)

	for _, v := range malfeasanceNodes {
		malfeasanceNodesMap[v.ID] = true
	}

	var maxHeight uint64 = 0
	atxID := ""

	for _, atx := range atxs {

		atxHeight := atx.BaseTick + atx.TickCount
		if atxHeight > uint64(maxHeight) && !malfeasanceNodesMap[atx.NodeID] {
			maxHeight = atxHeight
			atxID = atx.AtxID
		}

	}

	return atxID, nil
}

func hexToBase64(hexString string) (string, error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
