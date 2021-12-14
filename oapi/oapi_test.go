/*
Copyright 2019 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oapi_test

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
	. "github.com/rhdedgar/pod-logger/oapi"
)

var _ = Describe("Oapi", func() {
	config.AppSecrets.OAPIURL = "http://localhost:8080"
	config.AppSecrets.OAPIToken = "exampletdapitoken"
	config.AppSecrets.TDAPIURL = "http://localhost:8080/api/url/"
	config.AppSecrets.TDAPIUser = "exampletdapiuser"
	config.AppSecrets.TDAPIToken = "exampletdapitoken"
	config.LogURL = "http://localhost:8080/api/log"

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

			e.POST("/api/log", postLog)
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
		Context("Successful returns of namespace and pod objects from a mock openshift API", func() {
			It("Should return namespace and pod object with no errors", func() {
				ns, pod, err := GetInfo(testNS, testPod)

				Expect(err).To(BeNil())
				Expect(ns).ToNot(Equal(apinamespace.APINamespace{}))
				Expect(pod).ToNot(Equal(apipod.APIPod{}))
			})
		})
	})

	Describe("SendData", func() {
		Context("Successful POST of JSON to a mock LogWriter server", func() {
			It("Should return HTTP status code 200 with no errors", func() {
				status, err := SendData(models.Log{
					User:      "testUser",
					Namespace: "testNamespace",
					PodName:   "testPod",
					HostIP:    "10.10.10.10",
					PodIP:     "10.10.10.11",
					StartTime: time.Now(),
					UID:       "testuid",
				})

				Expect(err).To(BeNil())
				Expect(status).To(Equal(200))
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

// postLog mimicks a LogWriter server which receives models.Log POSTs as JSON. It's accessed with:
// POST /api/log
func postLog(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
