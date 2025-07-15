package services

import (
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
)

func New(log *logger.Log, ext *extensions.Extensions, util *utils.Utils) *Services {
	return &Services{
		HealthService: &HealthService{
			Log:        log,
			DB:         ext.Database,
			StringUtil: util.StringUtil,
		},
		GroupsService: &GroupsService{
			Log: log,
			DB:  ext.Database,
		},
		PortsService: &PortsService{
			Log: log,
			DB:  ext.Database,
		},
	}
}
