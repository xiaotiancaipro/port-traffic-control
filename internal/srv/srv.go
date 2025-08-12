package srv

import (
	"context"
	"fmt"
	"net/http"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/controllers"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/middlewares"
	"port-traffic-control/internal/routers"
	"port-traffic-control/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func New(config *configs.Configuration, log *logger.Log, controller *controllers.Controllers, util *utils.Utils) *Srv {

	gin.SetMode(gin.ReleaseMode)
	if config.API.Env != "production" {
		gin.SetMode(gin.DebugMode)
		log.Warning("The current service is in debug mode")
	}
	engine := gin.Default()

	middleware := middlewares.New(log, config, util)
	middleware.Mount(engine)

	router := routers.New(controller)
	router.Mount(engine)

	addr := fmt.Sprintf("%s:%d", config.API.Host, config.API.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	return &Srv{
		Config: config.API,
		Log:    log,
		Server: server,
	}

}

func (s *Srv) Start() {
	go func() {
		_ = s.Server.ListenAndServe()
	}()
	s.Log.Infof("Server started successfully, Address=%s", s.Server.Addr)
}

func (s *Srv) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		err = fmt.Errorf("server shutdown error, Error=%s", err.Error())
		s.Log.Error(err)
		return err
	}
	s.Log.Infof("Server stopped successfully")
	return nil
}
