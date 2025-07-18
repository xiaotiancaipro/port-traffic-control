package services

import (
	"github.com/florianl/go-tc"
	"gorm.io/gorm"
	"net"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
	"sync"
)

type Services struct {
	HealthService *HealthService
	GroupsService *GroupsService
	PortsService  *PortsService
	TCService     *TCService
}

type HealthService struct {
	Log        *logger.Log
	DB         *gorm.DB
	StringUtil *utils.StringUtil
}

type GroupsService struct {
	Log  *logger.Log
	DB   *gorm.DB
	Lock sync.RWMutex
}

type PortsService struct {
	Log *logger.Log
	DB  *gorm.DB
}

type TCService struct {
	Config     *configs.TCConfig
	Log        *logger.Log
	TC         *tc.Tc
	Iface      *net.Interface
	HandleRoot uint32
}
