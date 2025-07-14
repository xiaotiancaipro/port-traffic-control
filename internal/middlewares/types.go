package middlewares

import (
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
)

type Middlewares struct {
	Auth *Auth
}

type Auth struct {
	Log          *logger.Log
	Config       *configs.APIConfig
	ResponseUtil *utils.ResponseUtil
}
