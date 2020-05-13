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
