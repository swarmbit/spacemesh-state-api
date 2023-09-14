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
	} else if account.Address.Bytes() == nil {
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

func (a *AccountRoutes) GetAccountRewardsDetails(c *gin.Context) {
	accountAddress := c.Param("accountAddress")
	address, err := sTypes.StringToAddress(accountAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Wrong account address format",
		})
		return
	}

	currentEpoch, _, err := a.state.GetCurrentEpoch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get current epoch",
		})
		return
	}

	sumChannel := make(chan int64)
	sumErrChannel := make(chan error)
	go func(resultChan chan<- int64, errChan chan<- error) {
		rewardsSum, err := a.state.SumRewards(address)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- rewardsSum
		}
	}(sumChannel, sumErrChannel)

	sumEpochChannel := make(chan int64)
	sumEpochErrChannel := make(chan error)
	go func(resultChan chan<- int64, errChan chan<- error) {
		rewardsEpoch, err := a.state.SumRewardsEpoch(address, currentEpoch)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- rewardsEpoch
		}
	}(sumEpochChannel, sumEpochErrChannel)

	countEpochChannel := make(chan int64)
	countEpochErrChannel := make(chan error)
	go func(resultChan chan<- int64, errChan chan<- error) {
		countEpoch, err := a.state.CountTotalRewardsEpoch(address, currentEpoch)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- countEpoch
		}
	}(countEpochChannel, countEpochErrChannel)

	activeAtxChannel := make(chan []*types.NodeSmesher)
	activeAtxErrChannel := make(chan error)
	go func(resultChan chan<- []*types.NodeSmesher, errChan chan<- error) {
		activeAtx, err := a.state.GetActiveAtxPerAddress(currentEpoch-1, address)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- activeAtx
		}
	}(activeAtxChannel, activeAtxErrChannel)

	var sumResult int64
	var sumEpochResult int64
	var countEpochResult int64

	select {
	case sumResult = <-sumChannel:
	case <-sumErrChannel:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal server error",
			"error":  "Failed get sum",
		})
		return
	}

	select {
	case sumEpochResult = <-sumEpochChannel:
	case <-sumEpochErrChannel:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal server error",
			"error":  "Failed get sum for epoch",
		})
		return
	}

	select {
	case countEpochResult = <-countEpochChannel:
	case <-countEpochErrChannel:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal server error",
			"error":  "Failed get count for epoch",
		})
		return
	}

	c.JSON(200, &types.RewardDetails{
		TotalSum:                 sumResult,
		CurrentEpoch:             int64(currentEpoch),
		CurrentEpochRewardsSum:   sumEpochResult,
		CurrentEpochRewardsCount: countEpochResult,
	})
}

func (a *AccountRoutes) GetAccountRewardsEligibilities(c *gin.Context) {
	accountAddress := c.Param("accountAddress")
	address, err := sTypes.StringToAddress(accountAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Bad Request",
			"error":  "Wrong account address format",
		})
		return
	}

	currentEpoch, _, err := a.state.GetCurrentEpoch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal Error",
			"error":  "Failed to get current epoch",
		})
		return
	}

	activeAtxChannel := make(chan []*types.NodeSmesher)
	activeAtxErrChannel := make(chan error)
	go func(resultChan chan<- []*types.NodeSmesher, errChan chan<- error) {
		activeAtx, err := a.state.GetActiveAtxPerAddress(currentEpoch-1, address)
		if err != nil {
			errChan <- err
		} else {
			resultChan <- activeAtx
		}
	}(activeAtxChannel, activeAtxErrChannel)

	var activeAtxResult []*types.NodeSmesher

	select {
	case activeAtxResult = <-activeAtxChannel:
	case <-activeAtxErrChannel:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Internal server error",
			"error":  "Failed get smeshers",
		})
		return
	}

	eligibilities := make([]*types.Eligibility, len(activeAtxResult))
	eligibilityChannel := make(chan *EligibilityResult, len(activeAtxResult))

	for _, atx := range activeAtxResult {
		go func(results chan<- *EligibilityResult, atx *types.NodeSmesher) {
			count, err := a.state.GetEligibilityCount(atx.NodeID)
			if err != nil && count == -1 {
				results <- &EligibilityResult{
					Err: err,
				}
			} else if err != nil {
				results <- &EligibilityResult{
					Err: err,
				}
			} else {
				predictedRewards, err := a.state.GetLatestAVGLayerReward()
				if err != nil {
					results <- &EligibilityResult{
						Err: err,
					}
				} else {
					results <- &EligibilityResult{
						Value: &types.Eligibility{
							Address:           accountAddress,
							EffectiveNumUnits: atx.EffectiveNumUnits,
							Count:             count,
							PredictedRewards:  predictedRewards * uint64(count),
							SmesherId:         atx.NodeID.String(),
						},
						Err: nil,
					}
				}

			}

		}(eligibilityChannel, atx)
	}

	for i := 0; i < len(activeAtxResult); i++ {
		result := <-eligibilityChannel
		if result.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal server error",
				"error":  "Failed get eligibilities",
			})
			return
		} else {
			eligibilities[i] = result.Value
		}
	}

	c.JSON(200, eligibilities)
}

type EligibilityResult struct {
	Value *types.Eligibility
	Err   error
}
