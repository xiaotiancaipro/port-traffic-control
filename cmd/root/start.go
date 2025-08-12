package root

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/controllers"
	"port-traffic-control/internal/extensions"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/services"
	"port-traffic-control/internal/srv"
	"port-traffic-control/internal/utils"
	"syscall"

	"github.com/spf13/cobra"
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

	util := utils.New(log)

	pidFile := filepath.Join(config.API.PIDPath, config.API.PIDFile)
	if _, err = os.Stat(pidFile); !os.IsNotExist(err) {
		pid, err_ := util.ProcessUtil.CheckRunning(pidFile)
		if err_ != nil {
			cmd.PrintErrf("Service already running, PID=%d\n", pid)
			os.Exit(1)
		}
		cmd.PrintErrf("The service is not running, but the process file exists, PIDFile=%s\n", pidFile)
		os.Exit(1)
	}

	if _, err = os.Stat(config.API.PIDPath); os.IsNotExist(err) {
		if err_ := os.Mkdir(config.API.PIDPath, 0755); err_ != nil {
			cmd.PrintErrf("Failed to create path, Path=%s, Error=%v\n", config.API.PIDPath, err_)
			os.Exit(1)
		}
	}

	if _, err = os.Stat(config.Database.Path); os.IsNotExist(err) {
		if err_ := os.Mkdir(config.Database.Path, 0755); err_ != nil {
			cmd.PrintErrf("Failed to create path, Path=%s, Error=%v\n", config.Database.Path, err_)
			os.Exit(1)
		}
	}

	ext, err := extensions.New(config)
	if err != nil {
		cmd.PrintErrf("Middleware loading failed, Error=%v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err_ := ext.Close(); err_ != nil {
			cmd.PrintErrf("Failed to close extension, Error=%v\n", err_)
			os.Exit(1)
		}
	}()

	service := services.New(config, log, ext, util)

	controller := controllers.New(log, service, util)

	server := srv.New(config, log, controller, util)
	server.Start()
	cmd.Printf("Server started successfully, Address=%s\n", server.Server.Addr)

	pid := os.Getpid()
	if err_ := util.ProcessUtil.WritePIDFile(pidFile, pid); err_ != nil {
		cmd.PrintErrf("Failed to write PID file, Error=%v\n", err_)
		os.Exit(1)
	}
	defer func() { _ = os.Remove(pidFile) }()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	info := fmt.Sprintf("Received shutdown signal, service is shutting down, SIG=%v", sig)
	log.Info(info)
	cmd.Println(info)

	if err = server.Stop(); err != nil {
		cmd.PrintErrf("Failed to stop server, Error=%v\n", err)
		os.Exit(1)
	}
	cmd.Println("Server stopped successfully")

}
