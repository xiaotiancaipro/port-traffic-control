package services

import (
	"gorm.io/gorm"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
)

type Services struct {
	HealthService *HealthService
	GroupsService *GroupsService
	PortsService  *PortsService
}

type HealthService struct {
	Log        *logger.Log
	DB         *gorm.DB
	StringUtil *utils.StringUtil
}

type GroupsService struct {
	Log *logger.Log
	DB  *gorm.DB
}

type PortsService struct {
	Log *logger.Log
	DB  *gorm.DB
}
