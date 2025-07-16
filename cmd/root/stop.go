package root

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"port-traffic-control/internal/configs"
	"port-traffic-control/internal/logger"
	"port-traffic-control/internal/utils"
	"syscall"
	"time"
)

type Stop struct{}

func (s Stop) Init(api *cobra.Command) {
	api.AddCommand(s.cmd())
}

func (s Stop) cmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "stop",
		Short: "Stop api server",
		Args:  cobra.ExactArgs(0),
		Run:   s.run,
	}
	return command
}

func (Stop) run(cmd *cobra.Command, _ []string) {

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

	processUtil := utils.NewProcessUtil(log)

	pidFile := filepath.Join(config.API.PIDPath, config.API.PIDFile)
	if _, err = os.Stat(pidFile); os.IsNotExist(err) {
		cmd.PrintErrf("The process file does not exist, please start the service first")
		os.Exit(1)
	}

	pid, err := processUtil.CheckRunning(pidFile)
	if err != nil {
		cmd.PrintErrf("The service is not started, please start the service first")
		os.Exit(1)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		cmd.PrintErrf("Process not found, PID=%d, Error=%v\n", pid, err)
		os.Exit(1)
	}

	if err_ := process.Signal(syscall.SIGINT); err_ != nil {
		cmd.PrintErrf("Failed to send close signal, PID=%d, Error=%v\n", pid, err_)
		os.Exit(1)
	}
	cmd.Printf("Shutdown signal sent successfully, PID=%d\n", pid)

	startTime := time.Now()
	for {
		pid, err = processUtil.CheckRunning(pidFile)
		if err != nil {
			cmd.Printf("Server stopped successfully, TimeConsumed=%v\n", time.Since(startTime))
			return
		}
		time.Sleep(200 * time.Millisecond)
	}

}
