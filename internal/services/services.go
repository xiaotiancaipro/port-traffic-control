package services

import (
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
	"sync"
)

func New(config *configs.Configuration, log *logger.Log, ext *extensions.Extensions, util *utils.Utils) *Services {
	healthService := &HealthService{
		Log:        log,
		DB:         ext.Database,
		StringUtil: util.StringUtil,
	}
	groupsService := &GroupsService{
		Log:  log,
		DB:   ext.Database,
		Lock: &sync.RWMutex{},
	}
	portsService := &PortsService{
		Log: log,
		DB:  ext.Database,
	}
	tcService := &TCService{
		Config:     config.TC,
		Log:        log,
		TC:         ext.TC.TC_,
		Iface:      ext.TC.Iface,
		HandleRoot: ext.TC.HandleRoot,
	}
	return &Services{
		HealthService: healthService,
		GroupsService: groupsService,
		PortsService:  portsService,
		TCService:     tcService,
	}
}
