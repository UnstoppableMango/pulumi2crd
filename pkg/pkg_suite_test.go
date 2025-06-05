package pulumi2crd_test

import (
	"embed"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"gopkg.in/yaml.v3"
)

//go:embed testdata
var testdata embed.FS

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg Suite")
}

func readPackageFile(name string) *schema.Package {
	GinkgoHelper()

	data, err := testdata.ReadFile(filepath.Join("testdata", name))
	Expect(err).NotTo(HaveOccurred())

	var spec schema.PackageSpec
	Expect(yaml.Unmarshal(data, &spec)).To(Succeed())
	pkg, err := schema.ImportSpec(spec, map[string]schema.Language{}, schema.ValidationOptions{})
	Expect(err).NotTo(HaveOccurred())

	return pkg
}
