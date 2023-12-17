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
	epochRoutes := NewEpochRoutes(readDB, networkUtils, state)
	layersRoutes := NewLayersRoutes(readDB, networkUtils, state)
	transactionRoutes := NewTransactionRoutes(readDB, networkUtils, state)

	router.GET("/account", func(c *gin.Context) {
		accountRoutes.GetAccounts(c)
	})

	router.POST("/account/group", func(c *gin.Context) {
		accountRoutes.GetAccountGroup(c)
	})

	router.GET("/account/post/epoch/:epoch", func(c *gin.Context) {
		accountRoutes.GetAccountsPost(c)
	})

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

	router.POST("/account/:accountAddress/atx/:epoch/filter-active-nodes", func(c *gin.Context) {
		accountRoutes.FilterEpochActiveNodes(c)
	})

	router.GET("/account/:accountAddress/atx/:epoch", func(c *gin.Context) {
		accountRoutes.GetEpochAtx(c)
	})

	router.GET("/network/info", func(c *gin.Context) {
		networkRoutes.GetInfo(c)
	})

	router.GET("/nodes", func(c *gin.Context) {
		nodeRoutes.GetNodes(c)
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

	router.GET("/epochs/:epoch", func(c *gin.Context) {
		epochRoutes.GetEpoch(c)
	})

	router.GET("/epochs/:epoch/atx", func(c *gin.Context) {
		epochRoutes.GetEpochAtx(c)
	})

	router.GET("/layers", func(c *gin.Context) {
		layersRoutes.GetLayers(c)
	})

	router.GET("/layers/:layer/transactions", func(c *gin.Context) {
		layersRoutes.GetLayerTransactions(c)
	})

	router.GET("/layers/:layer/rewards", func(c *gin.Context) {
		layersRoutes.GetLayerRewards(c)
	})

	router.GET("/transactions", func(c *gin.Context) {
		transactionRoutes.GetTransactions(c)
	})

	router.GET("/transactions/:transactionId", func(c *gin.Context) {
		transactionRoutes.GetTransaction(c)
	})

	router.GET("/poets", func(c *gin.Context) {
		poetRoutes.GetPoets(c)
	})
}
