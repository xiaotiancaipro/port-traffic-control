package root

import "github.com/spf13/cobra"

const version = "0.0.3"

type Version struct{}

func (v Version) Init(root *cobra.Command) {
	root.AddCommand(v.cmd())
}

func (v Version) cmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "version",
		Short: "Version information",
		Args:  cobra.ExactArgs(0),
		Run:   v.run,
	}
	return command
}

func (Version) run(cmd *cobra.Command, _ []string) {
	cmd.Printf("Version: v%s\n", version)
}
