package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/state"
)

func AddRoutes(state *state.State, router *gin.Engine) {
	accountRoutes := NewAccountRoutes(state)
	smesherRoutes := NewSmesherRoutes(state)
	networkRoutes := NewNetworkRoutes(state)

	router.GET("/account/:accountAddress", func(c *gin.Context) {
		accountRoutes.GetAccount(c)
	})

	router.GET("/account/:accountAddress/rewards", func(c *gin.Context) {
		accountRoutes.GetAccountRewards(c)
	})

	router.GET("/smesher/:smesherId/eligibility", func(c *gin.Context) {
		smesherRoutes.GetSmesherEligibility(c)
	})

	router.GET("/network/higestatx", func(c *gin.Context) {
		networkRoutes.GetHighestAtx(c)
	})

	router.GET("/network/info", func(c *gin.Context) {
		networkRoutes.GetInfo(c)
	})
}
