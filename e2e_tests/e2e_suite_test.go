package e2e_tests

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKubetries(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e tests")
}
