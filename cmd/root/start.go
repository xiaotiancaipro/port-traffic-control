package root

import (
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

}
