package main

import (
	"github.com/rhdedgar/pod-logger/routers"
)

func main() {
	e := routers.Routers
	e.Logger.Info(e.Start(":8880"))
}
