package controllers

import (
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/services"
	"port-traffic-control/internal/utils"
)

type Controllers struct {
	HealthController *HealthController
	GroupsController *GroupsController
	PortsController  *PortsController
}

type HealthController struct {
	Log           *logger.Log
	HealthService *services.HealthService
	ResponseUtil  *utils.ResponseUtil
}

type GroupsController struct {
	Log           *logger.Log
	GroupsService *services.GroupsService
}

type PortsController struct {
	Log          *logger.Log
	PortsService *services.PortsService
}
