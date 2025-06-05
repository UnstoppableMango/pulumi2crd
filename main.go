package main

import (
	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
	"github.com/unstoppablemango/ux/sdk/plugin/cli"
	"github.com/unstoppablemango/ux/sdk/plugin/cmd"
)

var Cli = cli.New(pulumi2crd.Plugin)

func main() {
	if err := cmd.Execute("pulumi2crd", Cli); err != nil {
		cli.Fail(err)
	}
}
