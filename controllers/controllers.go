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

package controllers

import (
	"log"

	"github.com/labstack/echo"

	"net/http"

	"github.com/rhdedgar/pod-logger/docker"
	"github.com/rhdedgar/pod-logger/models"
	"github.com/rhdedgar/pod-logger/oapi"
)

// PostCrioPodLog handles received crictl inspect data in json format. It's accessed with:
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

// PostDockerPodLog handles received docker inspect data in json format. It's accessed with:
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

// PostClamScanResult handles received clamAV scan result data in json format. It's accessed with:
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
