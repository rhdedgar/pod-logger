package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	Token  string
	APIURL string
)

func init() {
	var appSecrets models.AppSecrets

	filePath := "/secrets/api_config.json"
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	err = json.Unmarshal(fileBytes, &appSecrets)
	if err != nil {
		fmt.Println("Error Unmarshalling secrets json: ", err)
	}

	Token = appSecrets.Token
	APIURL = appSecrets.APIURL
}
