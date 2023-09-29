package route

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
	"github.com/swarmbit/spacemesh-state-api/types"
)

type AccountRoutes struct {
	db           *database.ReadDB
	networkUtils *network.NetworkUtils
}

func NewAccountRoutes(readDB *database.ReadDB, networkUtils *network.NetworkUtils) *AccountRoutes {
	return &AccountRoutes{
		db:           readDB,
		networkUtils: networkUtils,
	}
}

func (a *AccountRoutes) GetAccount(c *gin.Context) {
	accountAddress := c.Param("accountAddress")
	account, err := a.db.GetAccount(accountAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch account",
		})
		return
	}
	if account.Address == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "Not Found",
			"error":  "Account not found",
		})
		return
	}
	numberOfTransactions, err := a.db.CountTransactions(accountAddress)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch account",
		})
		return
	}
	numberOfRewards, err := a.db.CountRewards(accountAddress)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch account",
		})
		return
	}

	c.JSON(200, &types.Account{
		Balance: account.Balance,
		// legacy
		BalanceDisplay:       "",
		Address:              accountAddress,
		TotalRewards:         account.TotalRewards,
		NumberOfTransactions: numberOfTransactions,
		Counter:              numberOfTransactions,
		NumberOfRewards:      numberOfRewards,
	})
}

func (a *AccountRoutes) GetAccountRewards(c *gin.Context) {
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

	accountAddress := c.Param("accountAddress")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Wrong account address format",
		})
		return
	}

	rewards, errRewards := a.db.GetRewards(accountAddress, int64(offset), int64(limit), sort)
	count, errCount := a.db.CountRewards(accountAddress)

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

func (a *AccountRoutes) GetAccountTransactions(c *gin.Context) {
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

	accountAddress := c.Param("accountAddress")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Wrong account address format",
		})
		return
	}

	transactions, errRewards := a.db.GetTransactions(accountAddress, int64(offset), int64(limit), sort, complete)
	count, errCount := a.db.CountTransactions(accountAddress)

	if errRewards != nil || errCount != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to fetch transactions for account",
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
				ID:              v.ID,
				Status:          v.Status,
				PrincipaAccount: v.PrincipaAccount,
				ReceiverAccount: v.ReceiverAccount,
				Fee:             v.Gas * v.GasPrice,
				Amount:          v.Amount,
				Layer:           v.Layer,
				Counter:         v.Counter,
				Method:          method,
			}
		}

		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, transactionsResponse)
	} else {
		c.Header("total", strconv.FormatInt(count, 10))
		c.JSON(200, make([]*types.Transaction, 0))
	}
}

func (a *AccountRoutes) GetAccountRewardsDetails(c *gin.Context) {
	accountAddress := c.Param("accountAddress")

	layer, err := a.db.GetLastProcessedLayer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get current epoch",
		})
		return
	}

	epoch := a.networkUtils.GetEpoch(uint64(layer.Layer))

	account, err := a.db.GetAccount(accountAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get rewards sum",
		})
		return
	}

	firstLayer := uint32(epoch * config.LayersPerEpoch)
	lastLayer := firstLayer + config.LayersPerEpoch

	countEpochResult, err := a.db.CountRewardsLayers(accountAddress, firstLayer, lastLayer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get epoch rewards count",
		})
		return
	}

	sumEpochResult, err := a.db.SumRewardsLayers(accountAddress, firstLayer, lastLayer)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get epoch rewards sum",
		})
		return
	}

	c.JSON(200, &types.RewardDetails{
		TotalSum:                 int64(account.TotalRewards),
		CurrentEpoch:             int64(epoch),
		CurrentEpochRewardsSum:   sumEpochResult,
		CurrentEpochRewardsCount: countEpochResult,
	})
}

func (a *AccountRoutes) GetAccountRewardsEligibilities(c *gin.Context) {
	accountAddress := c.Param("accountAddress")

	layer, err := a.db.GetLastProcessedLayer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get current epoch",
		})
		return
	}

	epoch := a.networkUtils.GetEpoch(uint64(layer.Layer))

	accountAtx, err := a.db.GetAtxWeightAccount(accountAddress, uint64(epoch-1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get account weight",
		})
		return
	}

	atxEpoch, err := a.db.GetAtxEpoch(uint64(epoch - 1))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get total weight",
		})
		return
	}

	eligibilityCount, err := a.networkUtils.GetNumberOfSlots(uint64(accountAtx.TotalWeight), atxEpoch.TotalWeight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get eligibility",
		})
		return
	}

	c.JSON(200, &types.Eligibility{
		Address:           accountAddress,
		Count:             eligibilityCount,
		EffectiveNumUnits: accountAtx.TotalEffectiveNumUnits,
	})
}
