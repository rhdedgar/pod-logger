package main

import (
	"github.com/rhdedgar/pod-logger/routers"
)

func main() {
	// Pod-logger V0.0.5

	e := routers.Routers
	e.HideBanner = true
	e.HidePort = true

	e.Logger.Fatal(e.Start(":8080"))
}
