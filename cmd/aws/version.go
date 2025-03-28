package main

import (
	"get.porter.sh/plugin/aws/pkg/aws"
	"get.porter.sh/porter/pkg/porter/version"
	"github.com/spf13/cobra"
)

func buildVersionCommand(p *aws.Plugin) *cobra.Command {
	opts := version.Options{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the plugin version",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.PrintVersion(opts)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.RawFormat, "output", "o", string(version.DefaultVersionFormat),
		"Specify an output format.  Allowed values: json, plaintext")

	return cmd
}
