package controllers

import (
	"log"

	"github.com/labstack/echo"

	"net/http"

	"github.com/rhdedgar/pod-logger/docker"
	"github.com/rhdedgar/pod-logger/models"
	"github.com/rhdedgar/pod-logger/oapi"
)

// POST /api/crio/log
func PostCrioPodLog(c echo.Context) error {
	var container models.Container

	if err := c.Bind(&container); err != nil {
		log.Println("Error binding received crio data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process crio content"}
	}

	go oapi.PrepCrioInfo(container)

	return c.NoContent(http.StatusOK)
}

// POST /api/docker/log
func PostDockerPodLog(c echo.Context) error {
	var container docker.DockerContainer

	if err := c.Bind(&container); err != nil {
		log.Println("Error binding received docker data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process docker content"}
	}

	go oapi.PrepDockerInfo(container)

	return c.NoContent(http.StatusOK)
}

// POST /api/clam/scanresult
func PostClamScanResult(c echo.Context) error {
	var scanResult models.ScanResult

	if err := c.Bind(&scanResult); err != nil {
		log.Println("Error binding received scan result data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process scan result"}
	}

	//fmt.Println("Scan result bound to:")
	log.Printf("%+v", scanResult)

	go oapi.PrepClamInfo(scanResult)

	return c.NoContent(http.StatusOK)
}
