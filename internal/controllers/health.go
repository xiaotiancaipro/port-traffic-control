package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"port-traffic-control/internal/services"
)

func (hc *HealthController) Health(c *gin.Context) {
	err := hc.HealthService.Check()
	if err != nil {
		hc.Log.Errorf("Health Check Error")
		hc.ResponseUtil.Error(c, services.InternalServerError)
		return
	}
	info := fmt.Sprintf("Service is running")
	hc.Log.Info(info)
	hc.ResponseUtil.Success(c, info, nil)
	return
}
