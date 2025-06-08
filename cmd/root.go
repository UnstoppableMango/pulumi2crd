package cmd

import (
	"github.com/spf13/cobra"
	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
	"github.com/unstoppablemango/ux/sdk/plugin/cli"
	uxcmd "github.com/unstoppablemango/ux/sdk/plugin/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "pulumi2crd",
	Short: "Converts Pulumi package specs to Custom Resource Definitions",
	Run: func(cmd *cobra.Command, args []string) {
		app := cli.New(pulumi2crd.Plugin)
		if err := uxcmd.Execute(cmd.Use, app); err != nil {
			cli.Fail(err)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}
