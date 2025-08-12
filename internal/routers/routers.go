package routers

import (
	"port-traffic-control/internal/controllers"

	"github.com/gin-gonic/gin"
)

const (
	Health = "/health"
	Groups = "/groups"
)

func New(controller *controllers.Controllers) *Routers {
	return &Routers{
		controller,
	}
}

func (r *Routers) Mount(engine *gin.Engine) {

	health := engine.Group(Health)
	{
		health.GET("/", r.HealthController.Health)
	}

	groups := engine.Group(Groups)
	{
		groups.POST("/create", r.GroupsController.Create)
	}

}
