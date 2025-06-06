package pulumi2crd

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kubebuilder/v4/pkg/model/resource"
)

func ConvertResources(spec *schema.Package) ([]*extensionv1.CustomResourceDefinition, error) {
	c := Converter{}

	crds := []*extensionv1.CustomResourceDefinition{}
	for _, r := range spec.Resources {
		if crd, err := c.Convert(r); err != nil {
			return nil, err
		} else {
			crds = append(crds, crd)
		}
	}

	return crds, nil
}

type Converter struct {
	Domain    string
	GroupName string
}

// https://github.com/kubernetes-sigs/controller-tools/blob/main/pkg/crd/gen.go

func (c Converter) Convert(spec *schema.Resource) (*extensionv1.CustomResourceDefinition, error) {
	_, _, typ, err := DecomposeToken(spec.Token)
	if err != nil {
		return nil, err
	}

	plural := resource.RegularPlural(typ)

	return &extensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s.%s", plural, c.Group()),
		},
		Spec: extensionv1.CustomResourceDefinitionSpec{
			Group: c.Group(),
			Names: extensionv1.CustomResourceDefinitionNames{
				Kind:     typ,
				ListKind: ListKind(typ),
				Singular: strings.ToLower(typ),
				Plural:   plural,
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
				Subresources: &extensionv1.CustomResourceSubresources{
					Status: &extensionv1.CustomResourceSubresourceStatus{},
				},
				AdditionalPrinterColumns: []extensionv1.CustomResourceColumnDefinition{},
			}},
		},
	}, nil
}

func (c Converter) Group() string {
	return fmt.Sprintf("%s.%s", c.GroupName, c.Domain)
}

func ConvertTypes(spec schema.Package) map[string]extensionv1.JSONSchemaProps {
	types := map[string]extensionv1.JSONSchemaProps{}
	for _, typ := range spec.Types {
		ConvertType(types, typ)
	}

	return types
}

func ConvertType(types map[string]extensionv1.JSONSchemaProps, typ schema.Type) {
	switch typ := typ.(type) {
	case *schema.ObjectType:
		types[typ.Token] = ConvertObjectType(typ)
	}
}

func PrimitiveType(typ schema.Type) string {
	switch typ := typ.(type) {
	case *schema.OptionalType:
		return PrimitiveType(typ.ElementType)
	case *schema.InputType:
		return PrimitiveType(typ.ElementType)
	default:
		// schema.primitiveType
		return typ.String()
	}
}

func ConvertObjectType(typ *schema.ObjectType) extensionv1.JSONSchemaProps {
	props := map[string]extensionv1.JSONSchemaProps{}
	for _, p := range typ.Properties {
		props[p.Name] = ConvertProperty(p)
	}

	return extensionv1.JSONSchemaProps{
		Description: typ.Comment,
		Properties:  props,
	}
}

func ConvertProperty(prop *schema.Property) extensionv1.JSONSchemaProps {
	return extensionv1.JSONSchemaProps{
		Type: PrimitiveType(prop.Type),
	}
}

func Spec(spec *schema.Resource) extensionv1.JSONSchemaProps {
	props := map[string]extensionv1.JSONSchemaProps{}
	for _, p := range spec.InputProperties {
		props[p.Name] = ConvertProperty(p)
	}

	return extensionv1.JSONSchemaProps{
		Description: spec.Comment,
		Properties:  props,
	}
}

func Status(spec *schema.Resource) extensionv1.JSONSchemaProps {
	props := map[string]extensionv1.JSONSchemaProps{}
	for _, p := range spec.Properties {
		props[p.Name] = ConvertProperty(p)
	}

	return extensionv1.JSONSchemaProps{
		Description: spec.Comment,
		Properties:  props,
	}
}
