package controllers

import (
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/services"
	"port-traffic-control/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	TCService     *services.TCService
	PortsService  *services.PortsService
	ResponseUtil  *utils.ResponseUtil
}

type PortsController struct {
	Log           *logger.Log
	DB            *gorm.DB
	GroupsService *services.GroupsService
	PortsService  *services.PortsService
	TCService     *services.TCService
	ResponseUtil  *utils.ResponseUtil
}

// RequestBody

type RequestBodyGroupsCreate struct {
	Bandwidth  int32 `json:"bandwidth"`
	PortMaxNum int32 `json:"portMaxNum"`
}

type RequestBodyGroupsGet struct {
	GroupID string `json:"groupID"`
}

type RequestBodyGroupsDelete struct {
	GroupID string `json:"groupID"`
}

type RequestBodyPorts struct {
	GroupID  string  `json:"groupID"`
	PortList []int32 `json:"portList"`
}

// ResponseBody

type ResponseBodyGroupsCreate struct {
	GroupID uuid.UUID `json:"groupId"`
}

type ResponseBodyGroupsGet struct {
	Bandwidth  int32   `json:"bandwidth"`
	PortMaxNum int32   `json:"PortMaxNum"`
	PortList   []int32 `json:"portList"`
}

type ResponseBodyGroupsList struct {
	Groups []uuid.UUID `json:"groups"`
}

type ResponseBodyPorts struct {
	Successful []int32                       `json:"successful"`
	Failed     []ResponseBodyPortsFailedItem `json:"failed"`
}

type ResponseBodyPortsFailedItem struct {
	Port  int32
	Error string
}
