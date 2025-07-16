package srv

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/controllers"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/middlewares"
	"port-traffic-control/internal/routers"
	"port-traffic-control/internal/utils"
	"time"
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
		Config:    config.API,
		Log:       log,
		Server:    server,
		ServerErr: make(chan error),
	}

}

func (s *Srv) Start() error {
	go func() {
		err := s.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.ServerErr <- err
		}
	}()
	if s.ServerErr != nil {
		err := fmt.Errorf("server start error, Error=%s", <-s.ServerErr)
		s.Log.Error(err)
		return err
	}
	s.Log.Infof("Service started, Address=%s", s.Server.Addr)
	return nil
}

func (s *Srv) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		err = fmt.Errorf("server shutdown error, Error=%s", err.Error())
		s.Log.Error(err)
		return err
	}
	s.Log.Infof("Server stopped")
	return nil
}
