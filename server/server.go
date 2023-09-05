package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/sql/accounts"
	"github.com/swarmbit/spacemesh-db-connector/state"
)

const GenesisEpochSeconds = 1689321600
const LayerDuration = 300

type Account struct {
	Balance        int64  `json:"balance"`
	BalanceDisplay string `json:"balanceDisplay"`
	Counter        int64  `json:"counter"`
	Address        string `json:"address"`
}

type Reward struct {
	Rewards        int64  `json:"rewards"`
	RewardsDisplay string `json:"rewardsDisplay"`
	Layer          int64  `json:"layer"`
	SmesherId      string `json:"smesherId"`
	Time           string `json:"time"`
	Timestamp      int64  `json:"timestamp"`
}

func StartServer() {

	db, err := state.StartDB("<state.sql path>", 10)
	if err != nil {
		fmt.Print("Failed to open db")
	}

	router := gin.Default()
	router.GET("/account/:accountAddress", func(c *gin.Context) {
		accountAddress := c.Param("accountAddress")
		address, err := types.StringToAddress(accountAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Bad Request",
				"error":  "Wrong account address format",
			})
			return
		}

		account, err := accounts.Latest(db, address)
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
			c.JSON(200, &Account{
				Balance: int64(account.Balance),
				// legacy
				BalanceDisplay: "",
				Address:        accountAddress,
				Counter:        int64(account.NextNonce),
			})
		}
	})

	router.GET("/account/:accountAddress/rewards", func(c *gin.Context) {

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
		address, err := types.StringToAddress(accountAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "Bad Request",
				"error":  "Wrong account address format",
			})
			return
		}

		rewards, errRewards := state.ListRewardsPaginated(db, address, int64(offset), int64(limit))
		count, errCount := state.CountTotalRewards(db, address)

		if errRewards != nil || errCount != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Internal Error",
				"error":  "Failed to fetch rewards for account",
			})
		} else if rewards != nil {

			rewardsResponse := make([]*Reward, len(rewards))

			for i, v := range rewards {
				rewardsResponse[i] = &Reward{
					Rewards: int64(v.TotalReward),
					// legacy
					RewardsDisplay: "",
					Layer:          int64(v.Layer),
					SmesherId:      v.NodeID.String(),
					// legacy
					Time:      "2023-09-05T00:00:00Z",
					Timestamp: GenesisEpochSeconds + (int64(v.Layer) * LayerDuration),
				}
			}

			c.Header("total", strconv.FormatInt(count, 10))
			c.JSON(200, rewardsResponse)
		} else {
			c.Header("total", strconv.FormatInt(count, 10))
			c.JSON(200, make([]*Reward, 0))
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}
