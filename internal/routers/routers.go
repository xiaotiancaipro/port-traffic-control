package routers

import (
	"github.com/gin-gonic/gin"
	"port-traffic-control/internal/controllers"
)

const (
	HealthRouter = "/health"
)

func New(controller *controllers.Controllers) *Routers {
	return &Routers{
		controller,
	}
}

func (r *Routers) Mount(engine *gin.Engine) {

	health := engine.Group(HealthRouter)
	{
		health.GET("/", r.HealthController.Health)
	}

}
