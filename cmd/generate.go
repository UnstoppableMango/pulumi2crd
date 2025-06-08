package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/unmango/go/codec"
	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
	"github.com/unstoppablemango/ux/sdk/plugin/cli"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate directly without ux",
	Aliases: []string{"gen"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		input := args[0]
		f, err := os.Open(input)
		if err != nil {
			cli.Fail(err)
		}

		data, err := io.ReadAll(f)
		if err != nil {
			cli.Fail(err)
		}

		out, err := pulumi2crd.Generate(data, getCodec(input))
		if err != nil {
			cli.Fail(err)
		}

		fmt.Println(string(out))
	},
}

func getCodec(path string) codec.Codec {
	if ext := filepath.Ext(path); ext == ".json" {
		return codec.Json
	} else {
		return codec.GoYaml
	}
}
