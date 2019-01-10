package controllers

import (
	"github.com/labstack/echo"

	"fmt"
	"github.com/rhdedgar/pod-logger/models"
	"github.com/rhdedgar/pod-logger/oapi"
	"net/http"
)

// POST /api/scan/log

// POST /api/pod/log
func PostApiPodLog(c echo.Context) error {
	var container models.Status

	if err := c.Bind(&container); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process content"}
	}

	go oapi.GetInfo(container)

	return c.NoContent(http.StatusOK)
}
