package main

import (
	"fmt"

	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/routers"
)

func main() {
	fmt.Println("TEST")
	config.SetConfig()

	e := routers.Routers
	e.Logger.Info(e.Start(":8080"))
}
