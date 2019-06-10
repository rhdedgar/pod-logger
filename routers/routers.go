package routers

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/rhdedgar/pod-logger/controllers"
)

var (
	Routers *echo.Echo
)

func init() {
	Routers = echo.New()

	Routers.Use(middleware.Logger())
	Routers.Use(middleware.Recover())

	Routers.POST("/api/crio/log", controllers.PostCrioPodLog)
	Routers.POST("/api/docker/log", controllers.PostDockerPodLog)
}
