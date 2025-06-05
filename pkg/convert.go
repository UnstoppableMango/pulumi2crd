package pulumi2crd

import (
	"fmt"

	"github.com/pulumi/pulumi/pkg/codegen/schema"
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

// https://github.com/kubernetes-sigs/controller-tools/blob/main/pkg/crd/gen.go

func Convert(name string, spec schema.ResourceSpec) *extensionv1.CustomResourceDefinition {
	plural := resource.RegularPlural(name)

	return &extensionv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: extensionv1.CustomResourceDefinitionSpec{
			Group: "",
			Names: extensionv1.CustomResourceDefinitionNames{
				Plural:     "",
				Singular:   "",
				ShortNames: []string{},
				Kind:       "",
				ListKind:   "",
				Categories: []string{},
			},
			Scope: extensionv1.NamespaceScoped,
			Versions: []extensionv1.CustomResourceDefinitionVersion{{
				Name:    "",
				Served:  true,
				Storage: true,
				Schema: &extensionv1.CustomResourceValidation{
					OpenAPIV3Schema: &extensionv1.JSONSchemaProps{
						// https://github.com/kubernetes-sigs/kubebuilder/blob/1d79aa1ec8204a11ae6be06cb96bae77dd0210bf/pkg/plugins/golang/deploy-image/v1alpha1/scaffolds/internal/templates/api/types.go#L117
						Description: fmt.Sprintf("%s is the Schema for the %s API", name, plural),
						Type:        "object",
						Properties: map[string]extensionv1.JSONSchemaProps{
							"apiVersion": apiVersionProp,
							"kind":       kindProp,
							"metadata":   {Type: "object"},
							"spec":       Spec(name, spec),
						},
					},
				},
				Subresources:             &extensionv1.CustomResourceSubresources{},
				AdditionalPrinterColumns: []extensionv1.CustomResourceColumnDefinition{},
			}},
		},
	}
}

func ConvertResources(spec *schema.PackageSpec) []*extensionv1.CustomResourceDefinition {
	crds := []*extensionv1.CustomResourceDefinition{}
	for name, r := range spec.Resources {
		crds = append(crds, Convert(name, r))
	}

	return crds
}

func Spec(name string, spec schema.ResourceSpec) extensionv1.JSONSchemaProps {
	return extensionv1.JSONSchemaProps{
		Description: spec.Description,
	}
}
