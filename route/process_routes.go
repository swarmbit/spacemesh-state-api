package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
)

func AddRoutes(readDB *database.ReadDB, router *gin.Engine) {
	accountRoutes := NewAccountRoutes(readDB)
	networkRoutes := NewNetworkRoutes(readDB)

	//smesherRoutes := NewSmesherRoutes(state)

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

	router.GET("/network/info", func(c *gin.Context) {
		networkRoutes.GetInfo(c)
	})

	/*
		router.GET("/account/:accountAddress/rewards/eligibility", func(c *gin.Context) {
			accountRoutes.GetAccountRewardsEligibilities(c)
		})

		router.GET("/smesher/:smesherId/eligibility", func(c *gin.Context) {
			smesherRoutes.GetSmesherEligibility(c)
		})*/
}
