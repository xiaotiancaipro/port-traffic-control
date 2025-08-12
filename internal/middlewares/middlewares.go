package middlewares

import (
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"

	"github.com/gin-gonic/gin"
)

func New(log *logger.Log, config *configs.Configuration, util *utils.Utils) *Middlewares {
	return &Middlewares{
		Auth: &Auth{
			Log:          log,
			Config:       config.API,
			ResponseUtil: util.ResponseUtil,
		},
	}
}

func (m *Middlewares) Mount(engine *gin.Engine) {
	engine.Use(func(c *gin.Context) { m.Auth.mount(c) })
}
