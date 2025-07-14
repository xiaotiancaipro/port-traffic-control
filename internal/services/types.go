package services

import (
	"gorm.io/gorm"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
)

type Services struct {
	HealthService *HealthService
}

type HealthService struct {
	Log        *logger.Log
	DB         *gorm.DB
	StringUtil *utils.StringUtil
}
