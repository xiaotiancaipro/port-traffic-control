package controllers

import (
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
)

func New(log *logger.Log, ext *extensions.Extensions, util *utils.Utils) *Controllers {
	return &Controllers{
		HealthController: &HealthController{
			Log:          log,
			DB:           ext.Database,
			StringUtil:   util.StringUtil,
			ResponseUtil: util.ResponseUtil,
		},
	}
}
