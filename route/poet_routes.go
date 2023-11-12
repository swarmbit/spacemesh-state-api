package route

import (
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/config"
)

type PoetRoutes struct {
	configValues *config.Config
}

func NewPoetRoutes(configValues *config.Config) *PoetRoutes {
	routes := &PoetRoutes{
		configValues: configValues,
	}
	return routes
}

func (p *PoetRoutes) GetPoets(c *gin.Context) {
	c.JSON(200, p.configValues.Poets)
}
