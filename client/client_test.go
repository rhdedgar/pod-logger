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

package client_test

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/pod-logger/client"
)

type clientJSON struct {
	clientMessage string
}

var _ = Describe("Client", func() {
	var (
		e         = echo.New()
		targetAPI = "http://localhost:8080/api/client/test"
	)

	BeforeEach(func() {
		go func() {
			e = echo.New()

			e.Use(middleware.Logger())
			e.Use(middleware.Recover())

			e.POST("/api/client/test", postClient)

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

	Describe("MakeClient", func() {
		Context("Successful client creation and response from server", func() {
			It("Should return HTTP status code 200 with no errors", func() {
				recJSON := make(map[string]interface{})

				req, err := http.NewRequest("POST", targetAPI, nil)
				Expect(err).To(BeNil())

				status, err := MakeClient(req, &recJSON, "")

				Expect(err).To(BeNil())
				Expect(status).To(Equal(200))
				Expect(recJSON).ToNot(Equal(new(map[string]interface{})))
			})
		})
	})
})

// postClient mimicks a namespace deletion request response from the OpenShift API. It's accessed with:
// POST /api/client/test
func postClient(c echo.Context) error {
	j := &clientJSON{
		clientMessage: "Test completed successfully.",
	}

	return c.JSON(http.StatusOK, j)
}
