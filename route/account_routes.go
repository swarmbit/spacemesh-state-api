package route

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	sTypes "github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql"
	"github.com/spacemeshos/go-spacemesh/sql/accounts"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/state"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type AccountRoutes struct {
	state *state.State
	db    sql.Executor
}

func NewAccountRoutes(state *state.State) *AccountRoutes {
	return &AccountRoutes{
		state: state,
		db:    state.DB,
	}
}

func (a *AccountRoutes) GetAccount(c *gin.Context) {
	accountAddress := c.Param("accountAddress")
	address, err := sTypes.StringToAddress(accountAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Wrong account address format",
		})
		return
	}

	account, err := accounts.Latest(a.db, address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch account",
		})
	} else if account.TemplateAddress == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Not Found",
			"error":  "Account not found",
		})
	} else {
		c.JSON(200, &types.Account{
			Balance: int64(account.Balance),
			// legacy
			BalanceDisplay: "",
			Address:        accountAddress,
			Counter:        int64(account.NextNonce),
		})
	}
}

func (a *AccountRoutes) GetAccountRewards(c *gin.Context) {
	offsetStr := c.DefaultQuery("offset", "0")
	limitStr := c.DefaultQuery("limit", "20")

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

	accountAddress := c.Param("accountAddress")
	address, err := sTypes.StringToAddress(accountAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Wrong account address format",
		})
		return
	}

	rewards, errRewards := a.state.ListRewardsPaginated(address, int64(offset), int64(limit))
	count, errCount := a.state.CountTotalRewards(address)

	if errRewards != nil || errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch rewards for account",
		})
	} else if rewards != nil {

		rewardsResponse := make([]*types.Reward, len(rewards))

		for i, v := range rewards {
			rewardsResponse[i] = &types.Reward{
				Rewards: int64(v.TotalReward),
				// legacy
				RewardsDisplay: "",
				Layer:          int64(v.Layer),
				SmesherId:      v.NodeID.String(),
				// legacy
				Time:      "2023-09-05T00:00:00Z",
				Timestamp: config.GenesisEpochSeconds + (int64(v.Layer) * config.LayerDuration),
			}
		}

		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, rewardsResponse)
	} else {
		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, make([]*types.Reward, 0))
	}
}
