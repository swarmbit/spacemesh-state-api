package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/network"
	"github.com/swarmbit/spacemesh-state-api/price"
)

func AddRoutes(readDB *database.ReadDB, router *gin.Engine, priceResolver *price.PriceResolver, configValues *config.Config) {
	networkUtils := network.NewNetworkUtils()
	state := network.NewNetworkState(readDB, networkUtils, priceResolver)
	accountRoutes := NewAccountRoutes(readDB, networkUtils, state, priceResolver)
	networkRoutes := NewNetworkRoutes(state)
	poetRoutes := NewPoetRoutes(configValues)
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

	router.GET("/account/:accountAddress/rewards/details/:epoch", func(c *gin.Context) {
		accountRoutes.GetAccountRewardsDetailsEpoch(c)
	})

	router.GET("/account/:accountAddress/atx/:epoch/filter-active-nodes", func(c *gin.Context) {
		accountRoutes.FilterEpochActiveNodes(c)
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
	
	router.GET("/poets", func(c *gin.Context) {
		poetRoutes.GetPoets(c)
	})
}
