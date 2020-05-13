package clam_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestClam(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Clam Suite")
}
