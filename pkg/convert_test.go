package pulumi2crd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	extensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
)

var _ = Describe("Convert", func() {
	It("should convert simple resources", func() {
		pkg := parsePackageString(`# Simple resources
name: simple
resources:
  simple:coreutils:Cat:
    type: object
    description: Test description
    inputProperties:
      args:
        type: string
    properties:
      args:
        type: string`)

		crds, err := pulumi2crd.ConvertResources(pkg)

		Expect(err).NotTo(HaveOccurred())
		Expect(crds).NotTo(BeEmpty())

		crd := &extensionv1.CustomResourceDefinition{}
		Expect(crds).To(ContainElement(HaveField("Name", "cats.."), &crd))
		Expect(crd.Spec.Group).To(Equal("."))
		Expect(crd.Spec.Names.Kind).To(Equal("Cat"))
		Expect(crd.Spec.Names.ListKind).To(Equal("CatList"))
		Expect(crd.Spec.Names.Singular).To(Equal("cat"))
		Expect(crd.Spec.Names.Plural).To(Equal("cats"))
		Expect(crd.Spec.Scope).To(Equal(extensionv1.NamespaceScoped))
		Expect(crd.Spec.Versions).NotTo(BeEmpty())

		var version extensionv1.CustomResourceDefinitionVersion
		Expect(crd.Spec.Versions).To(ContainElement(
			HaveField("Name", "v1alpha1"), &version,
		))
		Expect(version.Deprecated).To(BeFalseBecause("Deprecated is false"))
		Expect(version.DeprecationWarning).To(BeNil())
		Expect(version.Served).To(BeTrueBecause("Served is true"))
		Expect(version.Storage).To(BeTrueBecause("Storage is true"))
		Expect(version.Subresources.Status).NotTo(BeNil())
		Expect(version.Schema.OpenAPIV3Schema).NotTo(BeNil())
		Expect(version.Schema.OpenAPIV3Schema.Type).To(Equal("object"))

		Expect(version.Schema.OpenAPIV3Schema.Properties).To(MatchAllKeys(Keys{
			"apiVersion": Not(BeNil()),
			"kind":       Not(BeNil()),
			"metadata":   Not(BeNil()),
			"spec": MatchFields(IgnoreExtras, Fields{
				"Description": Equal("Test description"),
				"Properties": MatchAllKeys(Keys{
					"args": MatchFields(IgnoreExtras, Fields{
						"Type": Equal("string"),
					}),
				}),
			}),
			"status": MatchFields(IgnoreExtras, Fields{
				"Properties": MatchAllKeys(Keys{
					"args": MatchFields(IgnoreExtras, Fields{
						"Type": Equal("string"),
					}),
				}),
			}),
		}))
	})
})
