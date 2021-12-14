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

package main

import (
	"os"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rhdedgar/pod-logger/routers"
)

func main() {
	// Pod-logger V0.0.12

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
