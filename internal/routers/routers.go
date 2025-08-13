package routers

import (
	"port-traffic-control/internal/controllers"

	"github.com/gin-gonic/gin"
)

func New(controller *controllers.Controllers) *Routers {
	return &Routers{
		controller,
	}
}

func (r *Routers) Mount(engine *gin.Engine) {

	health := engine.Group("/health")
	{
		health.GET("/", r.HealthController.Health)
	}

	groups := engine.Group("/groups")
	{
		groups.POST("/create", r.GroupsController.Create)
		groups.POST("/get", r.GroupsController.Get)
		groups.POST("/delete", r.GroupsController.Delete)
		groups.GET("/list", r.GroupsController.List)
	}

	ports := engine.Group("/ports")
	{
		ports.POST("/add", r.PortsController.Add)
		ports.POST("/remove", r.PortsController.Remove)
	}

}
