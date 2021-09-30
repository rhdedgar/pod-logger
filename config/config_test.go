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

package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
)

var _ = Describe("Config", func() {
	var (
		appSecrets     models.AppSecrets
		configFilePath = "./api_config_example.json"
	)

	Describe("LoadJSON", func() {
		Context("Valid JSON read from api_config_example.json file", func() {
			It("Should Unmarshal into type models.AppSecrets", func() {
				err := LoadJSON(&appSecrets, configFilePath)

				Expect(err).To(BeNil())
				Expect(appSecrets).ToNot(Equal(models.AppSecrets{}))
			})
		})
	})
})
