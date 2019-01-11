package controllers

import (
	"github.com/labstack/echo"

	"fmt"
	"net/http"

	"github.com/rhdedgar/pod-logger/models"
	"github.com/rhdedgar/pod-logger/oapi"
)

// POST /api/scan/log

// POST /api/pod/log
func PostApiPodLog(c echo.Context) error {
	var container models.Container

	if err := c.Bind(&container); err != nil {
		fmt.Println("Error binding received data:\n", err)
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process content"}
	}

	go oapi.GetInfo(container)

	return c.NoContent(http.StatusOK)
}
