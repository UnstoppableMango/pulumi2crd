package pulumi2crd

import (
	"fmt"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
)

var apiVersionProp = extensionv1.JSONSchemaProps{
	Description: `APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources`,
}

var kindProp = extensionv1.JSONSchemaProps{
	Description: `Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds`,
}

var metadataProp = extensionv1.JSONSchemaProps{Type: "object"}

type Converter struct {
	Domain string
}

// https://github.com/kubernetes-sigs/controller-tools/blob/main/pkg/crd/gen.go

func (c Converter) Convert(spec *schema.Resource) *extensionv1.CustomResourceDefinition {
	plural := resource.RegularPlural(spec.Token) // TODO: Extract name

	return &extensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: spec.Token,
		},
		Spec: extensionv1.CustomResourceDefinitionSpec{
			Group: c.Domain,
			Names: extensionv1.CustomResourceDefinitionNames{
				Plural:   plural,
				Singular: spec.Token,
				Kind:     "",
				ListKind: "",
			},
			Scope: extensionv1.NamespaceScoped,
			Versions: []extensionv1.CustomResourceDefinitionVersion{{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &extensionv1.CustomResourceValidation{
					OpenAPIV3Schema: &extensionv1.JSONSchemaProps{
						Description: SchemaDescription(spec.Token, plural),
						Type:        "object",
						Properties: map[string]extensionv1.JSONSchemaProps{
							"apiVersion": apiVersionProp,
							"kind":       kindProp,
							"metadata":   metadataProp,
							"spec":       Spec(spec),
							"status":     Status(spec),
						},
					},
				},
				Subresources:             &extensionv1.CustomResourceSubresources{},
				AdditionalPrinterColumns: []extensionv1.CustomResourceColumnDefinition{},
			}},
		},
	}
}

func ConvertResources(spec *schema.Package) ([]*extensionv1.CustomResourceDefinition, error) {
	c := Converter{}

	crds := []*extensionv1.CustomResourceDefinition{}
	for _, r := range spec.Resources {
		crds = append(crds, c.Convert(r))
	}

	return crds, nil
}

func ConvertTypes(spec schema.Package) map[string]extensionv1.JSONSchemaProps {
	types := map[string]extensionv1.JSONSchemaProps{}
	for _, t := range spec.Types {
		switch typ := t.(type) {
		case *schema.ObjectType:
			types[typ.Token] = ConvertObjectType(typ)
		}
	}

	return types
}

func ConvertObjectType(typ *schema.ObjectType) extensionv1.JSONSchemaProps {
	props := map[string]extensionv1.JSONSchemaProps{}
	for _, p := range typ.Properties {
		props[p.Name] = extensionv1.JSONSchemaProps{} // TODO
	}

	return extensionv1.JSONSchemaProps{
		Description: typ.Comment,
		Properties:  props,
	}
}

func Spec(spec *schema.Resource) extensionv1.JSONSchemaProps {
	props := extensionv1.JSONSchemaProps{
		Description: spec.Comment,
	}

	for _, p := range spec.InputProperties {
		props.Properties[p.Name] = extensionv1.JSONSchemaProps{
			Description: p.Comment,
			Type:        p.Type.String(),
		}
	}

	return props
}

func Status(spec *schema.Resource) extensionv1.JSONSchemaProps {
	return extensionv1.JSONSchemaProps{}
}

func SchemaDescription(name, plural string) string {
	// https://github.com/kubernetes-sigs/kubebuilder/blob/1d79aa1ec8204a11ae6be06cb96bae77dd0210bf/pkg/plugins/golang/deploy-image/v1alpha1/scaffolds/internal/templates/api/types.go#L117
	return fmt.Sprintf("%s is the Schema for the %s API", name, plural)
}
