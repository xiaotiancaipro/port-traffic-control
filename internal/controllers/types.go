package controllers

import (
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/services"
	"port-traffic-control/internal/utils"
)

type Controllers struct {
	HealthController *HealthController
}

type HealthController struct {
	Log           *logger.Log
	HealthService *services.HealthService
	ResponseUtil  *utils.ResponseUtil
}
