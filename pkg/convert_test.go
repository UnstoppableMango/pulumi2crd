package pulumi2crd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/pkg/codegen/schema"
	"gopkg.in/yaml.v3"

	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
)

var _ = Describe("Convert", func() {
	var pkg *schema.Package

	BeforeEach(func() {
		data, err := testdata.ReadFile("testdata/schema.yml")
		Expect(err).NotTo(HaveOccurred())

		var spec schema.PackageSpec
		Expect(yaml.Unmarshal(data, &spec)).To(Succeed())
		pkg, err = schema.ImportSpec(spec)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should work", func() {
		const name = "baremetal:coreutils:Cat"
		r := pkg.Resources[0]

		crd := pulumi2crd.Convert(r)

		Expect(crd).NotTo(BeNil())
		Expect(crd.Name).To(Equal(name))
		Expect(crd.Spec.Names.Plural).To(Equal("Cats"))
	})
})
