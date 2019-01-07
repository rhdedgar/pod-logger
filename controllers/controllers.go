package controllers

import (
	"github.com/labstack/echo"

	"fmt"
	"github.com/rhdedgar/pod-logger/models"
	"net/http"
)

// POST /api/scan/log

// POST /api/pod/log
func PostApiPodLog(c echo.Context) error {
	var pod models.Pod

	if err := c.Bind(&pod); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Failed to process content"}
	}

	fmt.Println(pod.Labels.IoKubernetesPodNamespace, pod.Labels.IoKubernetesPodName)

	return c.NoContent(http.StatusOK)
}
