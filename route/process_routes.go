package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
)

func AddRoutes(readDB *database.ReadDB, router *gin.Engine) {
	networkUtils := network.NewNetworkUtils()
	accountRoutes := NewAccountRoutes(readDB, networkUtils)
	networkRoutes := NewNetworkRoutes(readDB, networkUtils)
	nodeRoutes := NewNodeRoutes(readDB, networkUtils)

	router.GET("/account/:accountAddress", func(c *gin.Context) {
		accountRoutes.GetAccount(c)
	})

	router.GET("/account/:accountAddress/rewards", func(c *gin.Context) {
		accountRoutes.GetAccountRewards(c)
	})

	router.GET("/account/:accountAddress/transactions", func(c *gin.Context) {
		accountRoutes.GetAccountTransactions(c)
	})

	router.GET("/account/:accountAddress/rewards/details", func(c *gin.Context) {
		accountRoutes.GetAccountRewardsDetails(c)
	})

	router.GET("/account/:accountAddress/rewards/eligibility", func(c *gin.Context) {
		accountRoutes.GetAccountRewardsEligibilities(c)
	})

	router.GET("/network/info", func(c *gin.Context) {
		networkRoutes.GetInfo(c)
	})

	router.GET("/nodes/:nodeId", func(c *gin.Context) {
		nodeRoutes.GetNode(c)
	})

	router.GET("/nodes/:nodeId/eligibility", func(c *gin.Context) {
		nodeRoutes.GetEligibility(c)
	})
}
