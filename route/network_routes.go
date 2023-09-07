package route

import (
	"encoding/base64"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/swarmbit/spacemesh-state-api/state"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type NetworkRoutes struct {
	state *state.State
	db    sql.Executor
}

func NewNetworkRoutes(state *state.State) *NetworkRoutes {
	return &NetworkRoutes{
		state: state,
		db:    state.DB,
	}
}

func (n *NetworkRoutes) GetHighestAtx(c *gin.Context) {
	highest, err := n.state.GetHighestAtx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get higest ATX",
		})
		return
	}
	c.JSON(200, &types.HigestAtx{
		AtxHex:    hex.EncodeToString(highest.Bytes()),
		AtxBase64: base64.StdEncoding.EncodeToString(highest.Bytes()),
	})
}

func (n *NetworkRoutes) GetInfo(c *gin.Context) {

	epoch, _, err := n.state.GetCurrentEpoch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get current epoch",
		})
		return
	}
	totalCommited, err := n.state.GetTotalCommittedForEpoch(epoch)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get totalCommited",
		})
		return
	}
	circulatingSupply, err := n.state.GetCirculatingSupply()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get circulating supply",
		})
		return
	}
	totalAccounts, err := n.state.GetNumberOfAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get number of accounts",
		})
		return
	}
	c.JSON(200, &types.NetworkInfo{
		Epoch:             epoch.Uint32(),
		TotalCommited:     totalCommited,
		CirculatingSupply: circulatingSupply,
		TotalAccounts:     totalAccounts,
	})
}
