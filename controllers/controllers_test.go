package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//. "github.com/rhdedgar/pod-logger/controllers"
	"github.com/rhdedgar/pod-logger/controllers"
	"github.com/rhdedgar/pod-logger/docker"
	"github.com/rhdedgar/pod-logger/models"
)

var _ = Describe("Controllers", func() {
	var (
		crioPodLogURL   = "http://localhost:8080/api/crio/log"
		dockerPodLogURL = "http://localhost:8080/api/docker/log"
		clamLogURL      = "http://localhost:8080/api/clam/scanresult"
		crictlFilePath  = "./crictl_inspect_example.json"
		dockerFilePath  = "./docker_inspect_example.json"
		clamFilePath    = "./clam_scan_result_example.json"
		e               = echo.New()
	)

	BeforeEach(func() {
		go func() {
			e = echo.New()

			e.Use(middleware.Logger())
			e.Use(middleware.Recover())

			e.POST("/api/crio/log", controllers.PostCrioPodLog)
			e.POST("/api/docker/log", controllers.PostDockerPodLog)
			e.POST("/api/clam/scanresult", controllers.PostClamScanResult)

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

	Describe("PostDockerPodLog", func() {
		Context("Valid docker output POSTed to /api/docker/log", func() {
			var container docker.DockerContainer

			It("Should return HTTP status code 200", func() {
				statusCode, err := senderTest(container, dockerPodLogURL, dockerFilePath)
				Expect(err).To(BeNil())
				Expect(statusCode).To(Equal(200))
			})
		})
	})

	Describe("PostClamScanResult", func() {
		var scanResult models.ScanResult

		Context("Valid clam scan result POSTed to /api/clam/scanresult", func() {
			It("Should return HTTP status code 200", func() {
				statusCode, err := senderTest(scanResult, clamLogURL, clamFilePath)
				Expect(err).To(BeNil())
				Expect(statusCode).To(Equal(200))
			})
		})
	})

	Describe("PostCrioPodLog", func() {
		Context("Valid crictl output POSTed to /api/crio/log", func() {
			var container models.Container

			It("Should return HTTP status code 200", func() {
				statusCode, err := senderTest(container, crioPodLogURL, crictlFilePath)
				Expect(err).To(BeNil())
				Expect(statusCode).To(Equal(200))
			})
		})
	})

})

// senderTest POSTs the specified URL with prepopulated example data, and unmarshals it into the provide data structure.
func senderTest(ds interface{}, url, filePath string) (int, error) {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
		return 0, err
	}

	err = json.Unmarshal(fileBytes, &ds)
	if err != nil {
		fmt.Println("Error Unmarshaling secrets json: ", err)
		return 0, err
	}

	jsonStr, err := json.Marshal(ds)
	if err != nil {
		fmt.Println("Error Marshaling secrets json: ", err)
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error creating new HTTP request:", err)
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending to pod-logger at %v: %v \n", url, err)
		fmt.Printf("Could not send %v \n", string(jsonStr[:]))
		return 0, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	return resp.StatusCode, nil
}
