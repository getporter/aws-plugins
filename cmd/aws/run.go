package main

import (
	"get.porter.sh/plugin/aws/pkg/aws"
	"github.com/spf13/cobra"
)

func buildRunCommand(p *aws.Plugin) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run PLUGIN_IMPLEMENTATION",
		Short: "Run the plugin and listen for client connections.",
		Run: func(cmd *cobra.Command, args []string) {
			p.Run(args)
		},
	}

	return cmd
}
