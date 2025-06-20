package pulumi2crd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	filev1alpha1 "buf.build/gen/go/unmango/protofs/protocolbuffers/go/dev/unmango/file/v1alpha1"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/spf13/afero"
	"github.com/unmango/go/codec"
	uxv1alpha1 "github.com/unstoppablemango/ux/gen/dev/unmango/ux/v1alpha1"
	"github.com/unstoppablemango/ux/sdk/plugin"
)

var Plugin = plugin.New(
	plugin.WithGenerator(Generator{}),
	plugin.WithCapabilities(&uxv1alpha1.Capability{
		From: "pulumi",
		To:   "crd",
	}),
)

type Generator struct{}

// Generate implements ux.Generator.
func (g Generator) Generate(ctx context.Context, req *uxv1alpha1.GenerateRequest) (*uxv1alpha1.GenerateResponse, error) {
	fs, err := plugin.OutputFs(req)
	if err != nil {
		return nil, err
	}

	outputs := []*filev1alpha1.File{}
	for _, i := range req.Inputs {
		f, err := fs.Open(i.Name)
		if err != nil {
			return nil, err
		}

		data, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		var c codec.Codec
		switch filepath.Ext(f.Name()) {
		case ".json":
			c = codec.Json
		case ".yaml", ".yml":
			c = codec.GoYaml
		default:
			return nil, fmt.Errorf("unsupported file type: %s", f.Name())
		}

		var spec schema.PackageSpec
		if err := c.Unmarshal(data, &spec); err != nil {
			return nil, err
		}

		pkg, err := schema.ImportSpec(spec,
			map[string]schema.Language{},
			schema.ValidationOptions{},
		)
		if err != nil {
			return nil, err
		}

		crds, err := ConvertResources(pkg)
		if err != nil {
			return nil, err
		}

		for _, crd := range crds {
			yaml, err := codec.GoYaml.Marshal(crd)
			if err != nil {
				return nil, err
			}

			name := fmt.Sprint(crd.Name, ".yml")
			if err = afero.WriteFile(fs, name, yaml, os.ModePerm); err != nil {
				return nil, err
			}

			outputs = append(outputs, &filev1alpha1.File{
				Name: name,
			})
		}
	}

	return &uxv1alpha1.GenerateResponse{
		Outputs: outputs,
	}, nil
}
