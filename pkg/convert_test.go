package pulumi2crd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"

	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
)

var _ = Describe("Convert", func() {
	var pkg *schema.Package

	BeforeEach(func() {
		pkg = readPackageFile("schema.yml")
	})

	It("should work", func() {
		const name = "baremetal:coreutils:Cat"

		crds, err := pulumi2crd.ConvertResources(pkg)

		Expect(err).NotTo(HaveOccurred())
		Expect(crds).NotTo(BeEmpty())
		// Expect(crds.Name).To(Equal(name))
		// Expect(crds.Spec.Names.Plural).To(Equal("Cats"))
	})
})
