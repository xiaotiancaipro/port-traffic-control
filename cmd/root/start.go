package root

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/controllers"
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/middlewares"
	"port-traffic-control/internal/utils"
)

type Start struct{}

func (s Start) Init(api *cobra.Command) {
	api.AddCommand(s.cmd())
}

func (s Start) cmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "start",
		Short: "Start api server",
		Args:  cobra.ExactArgs(0),
		Run:   s.run,
	}
	return command
}

func (Start) run(cmd *cobra.Command, _ []string) {

	config, err := configs.New()
	if err != nil {
		cmd.PrintErrf("Configuration loading failed, Error=%v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(config.API.LogPath, config.API.LogPrefix)
	if err != nil {
		cmd.PrintErrf("Failed to initialize log, Error=%v\n", err)
		os.Exit(1)
	}

	ext, err := extensions.New(config)
	if err != nil {
		cmd.PrintErrf("Middleware loading failed, Error=%v\n", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	if config.API.Env != "production" {
		gin.SetMode(gin.DebugMode)
		log.Warning("The current service is in debug mode")
	}
	server := gin.Default()

	util := utils.New(log)
	controller := controllers.New(log, ext, util)

	middleware := middlewares.New(log, config, util)
	middleware.Mount(server)

	_, _, _ = ext, server, controller // TODO

}
