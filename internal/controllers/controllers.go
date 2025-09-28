package controllers

import (
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/services"
	"port-traffic-control/internal/utils"
)

func New(log *logger.Log, ext *extensions.Extensions, service *services.Services, util *utils.Utils) *Controllers {
	healthController := &HealthController{
		Log:           log,
		HealthService: service.HealthService,
		ResponseUtil:  util.ResponseUtil,
	}
	groupsController := &GroupsController{
		Log:           log,
		GroupsService: service.GroupsService,
		TCService:     service.TCService,
		PortsService:  service.PortsService,
		ResponseUtil:  util.ResponseUtil,
	}
	portsController := &PortsController{
		Log:           log,
		DB:            ext.Database,
		GroupsService: service.GroupsService,
		PortsService:  service.PortsService,
		TCService:     service.TCService,
		ResponseUtil:  util.ResponseUtil,
	}
	return &Controllers{
		HealthController: healthController,
		GroupsController: groupsController,
		PortsController:  portsController,
	}
}
