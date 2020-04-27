package main

import (
	"os"

	"github.com/labstack/echo/middleware"
	"github.com/rhdedgar/pod-logger/routers"
)

func main() {
	// Pod-logger V0.0.10

	e := routers.Routers
	e.HideBanner = true
	e.HidePort = true

	if os.Getenv("DEBUG_APP") == "true" {
		e.Use(middleware.Logger())
		e.Logger.Info(e.Start(":8080"))
	} else {
		e.Start(":8080")
	}
}
