package oapi_test

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/config"
	. "github.com/rhdedgar/pod-logger/oapi"
)

var _ = Describe("Oapi", func() {
	config.AppSecrets.OAPIURL = "http://localhost:8080"
	config.AppSecrets.OAPIToken = "exampletdapitoken"
	config.AppSecrets.TDAPIURL = "http://localhost:8080/api/url/"
	config.AppSecrets.TDAPIUser = "exampletdapiuser"
	config.AppSecrets.TDAPIToken = "exampletdapitoken"

	var (
		e       = echo.New()
		testPod = "example_pod"
		testNS  = "example_namespace"
	)

	BeforeEach(func() {
		go func() {
			e = echo.New()

			e.Use(middleware.Logger())
			e.Use(middleware.Recover())

			e.GET("/api/v1/namespaces/:namespace", getNamespace)
			e.GET("/api/v1/namespaces/:namespace/pods/:pod/status", getPod)

			e.Use(middleware.Logger())

			e.Logger.Info(e.Start(":8080"))
		}()
	})

	AfterEach(func() {
		//e.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	})

	Describe("GetInfo", func() {
		Context("Successful returns of namespace and pod objects from mock openshift API", func() {
			It("Should return namespace and pod object with no errors", func() {
				ns, pod, err := GetInfo(testNS, testPod)

				Expect(err).To(BeNil())
				Expect(ns).ToNot(Equal(apinamespace.APINamespace{}))
				Expect(pod).ToNot(Equal(apipod.APIPod{}))
			})
		})
	})
})

// getPod mimicks a pod definition object response from the OpenShift API. It's accessed with:
// GET /api/v1/namespaces/:namespace/pods/:pod/status
func getPod(c echo.Context) error {
	j := &apipod.APIPod{
		Kind:       "Example pod definition.",
		APIVersion: "v1",
	}

	return c.JSON(http.StatusOK, j)
}

// getPod mimicks a namespace definition object response from the OpenShift API. It's accessed with:
// GET /api/v1/namespaces/:namespace
func getNamespace(c echo.Context) error {
	j := &apinamespace.APINamespace{
		Kind:       "Example namespace definition.",
		APIVersion: "v1",
	}

	return c.JSON(http.StatusOK, j)
}
