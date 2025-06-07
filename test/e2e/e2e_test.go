package e2e_test

import (
	. "github.com/onsi/ginkgo/v2"

	pulumi2crd "github.com/unstoppablemango/pulumi2crd/pkg"
	"github.com/unstoppablemango/ux/pkg/plugin/conformance"
)

var _ = Describe("E2e", func() {
	_ = conformance.NewSuite(conformance.SuiteOptions{
		Plugin: pulumi2crd.Plugin,
	})
})
