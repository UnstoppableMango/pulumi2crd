package pulumi2crd

import (
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/unmango/go/codec"
)

func Generate(data []byte, codec codec.Codec) ([]byte, error) {
	var spec schema.PackageSpec
	if err := codec.Unmarshal(data, &spec); err != nil {
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

	return codec.Marshal(crds)
}
