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
	Routers.POST("/api/clam/scanresult", controllers.PostClamScanResult)
}
