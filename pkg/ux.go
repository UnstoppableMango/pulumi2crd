package pulumi2crd

import (
	"context"

	"github.com/pulumi/pulumi/pkg/codegen/schema"
	uxv1alpha1 "github.com/unstoppablemango/ux/gen/dev/unmango/ux/v1alpha1"
	"github.com/unstoppablemango/ux/pkg/payload"
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
	codec, err := payload.Codec(req.Payload)
	if err != nil {
		return nil, err
	}

	var spec schema.PackageSpec
	if err = codec.Unmarshal(req.Payload.Data, &spec); err != nil {
		return nil, err
	}

	pack, err := schema.ImportSpec(spec)
	if err != nil {
		return nil, err
	}

	crds := ConvertResources(pack)
	data, err := codec.Marshal(crds)
	if err != nil {
		return nil, err
	}

	res := &uxv1alpha1.GenerateResponse{
		Payload: &uxv1alpha1.Payload{
			ContentType: "application/yaml", // TODO
			Data:        data,
		},
	}

	return res, nil
}
