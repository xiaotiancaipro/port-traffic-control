package cmd

import (
	"port-traffic-control/cmd/root"

	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return Root{}.Init()
}

type Root struct{}

func (r Root) Init() *cobra.Command {
	return r.cmd()
}

func (r Root) cmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "port-traffic-control",
		Short: "Linux system port traffic control API interface server",
		Run:   r.run,
	}
	root.Version{}.Init(command)
	root.Start{}.Init(command)
	root.Stop{}.Init(command)
	return command
}

func (Root) run(cmd *cobra.Command, _ []string) {
	_ = cmd.Help()
}
