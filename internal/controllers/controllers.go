package controllers

import (
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/services"
	"port-traffic-control/internal/utils"
)

func New(log *logger.Log, service *services.Services, util *utils.Utils) *Controllers {
	return &Controllers{
		HealthController: &HealthController{
			Log:           log,
			HealthService: service.HealthService,
			ResponseUtil:  util.ResponseUtil,
		},
		GroupsController: &GroupsController{
			Log:           log,
			GroupsService: service.GroupsService,
			TCService:     service.TCService,
			ResponseUtil:  util.ResponseUtil,
		},
		PortsController: &PortsController{
			Log:          log,
			PortsService: service.PortsService,
			ResponseUtil: util.ResponseUtil,
		},
	}
}
