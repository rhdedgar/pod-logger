package oapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oapi Suite")
}
