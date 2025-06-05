package pulumi2crd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/pkg/codegen/schema"
	"gopkg.in/yaml.v3"

	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
)

var _ = Describe("Convert", func() {
	var pkg schema.PackageSpec

	BeforeEach(func() {
		data, err := testdata.ReadFile("testdata/schema.yml")
		Expect(err).NotTo(HaveOccurred())
		Expect(yaml.Unmarshal(data, &pkg)).To(Succeed())
	})

	It("should work", func() {
		const name = "baremetal:coreutils:Cat"
		r := pkg.Resources[name]

		crd := pulumi2crd.Convert(name, r)

		Expect(crd).NotTo(BeNil())
		Expect(crd.Name).To(Equal(name))
		Expect(crd.Spec.Names.Plural).To(Equal("Cats"))
	})
})
