package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
)

func AddRoutes(readDB *database.ReadDB, router *gin.Engine) {
	networkUtils := network.NewNetworkUtils()
	state := network.NewNetworkState(readDB, networkUtils)
	accountRoutes := NewAccountRoutes(readDB, networkUtils, state)
	networkRoutes := NewNetworkRoutes(state)
	nodeRoutes := NewNodeRoutes(readDB, networkUtils, state)

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

	router.GET("/nodes/:nodeId/rewards", func(c *gin.Context) {
		nodeRoutes.GetNodeRewards(c)
	})

	router.GET("/nodes/:nodeId/rewards/details", func(c *gin.Context) {
		nodeRoutes.GetNodeRewardsDetails(c)
	})

	router.GET("/nodes/:nodeId/rewards/eligibility", func(c *gin.Context) {
		nodeRoutes.GetEligibility(c)
	})
}
