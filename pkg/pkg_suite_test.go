package pulumi2crd_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"gopkg.in/yaml.v3"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg Suite")
}

func parsePackageString(yml string) *schema.Package {
	GinkgoHelper()

	return parsePackageData([]byte(yml))
}

func parsePackageData(data []byte) *schema.Package {
	GinkgoHelper()

	var spec schema.PackageSpec
	Expect(yaml.Unmarshal(data, &spec)).To(Succeed())
	pkg, err := schema.ImportSpec(spec, map[string]schema.Language{}, schema.ValidationOptions{})
	Expect(err).NotTo(HaveOccurred())

	return pkg
}
