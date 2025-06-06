package pulumi2crd

import (
	"fmt"

	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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

func SchemaDescription(name, plural string) string {
	// https://github.com/kubernetes-sigs/kubebuilder/blob/1d79aa1ec8204a11ae6be06cb96bae77dd0210bf/pkg/plugins/golang/deploy-image/v1alpha1/scaffolds/internal/templates/api/types.go#L117
	return fmt.Sprintf("%s is the Schema for the %s API", name, plural)
}

func ListKind(kind string) string {
	// https://github.com/kubernetes-sigs/controller-tools/blob/67249c598b39237877a67c43301ce92564efd115/pkg/crd/spec.go#L77
	return kind + "List"
}
