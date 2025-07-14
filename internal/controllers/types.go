package controllers

import (
	"gorm.io/gorm"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
)

type Controllers struct {
	HealthController *HealthController
}

type HealthController struct {
	Log          *logger.Log
	DB           *gorm.DB
	StringUtil   *utils.StringUtil
	ResponseUtil *utils.ResponseUtil
}
