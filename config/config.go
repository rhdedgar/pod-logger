package config

import (
	"encoding/json"
	"fmt"
	"github.com/rhdedgar/pod-logger/models"
	"io/ioutil"
)

var (
	Token string
	URL   string
)

func init() {
	var appSecrets models.AppSecrets

	filePath := "/secrets/api_config.json"
	//filePath := "/home/remote/dedgar/ansible/config_secrets.json"
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	err = json.Unmarshal(fileBytes, &appSecrets)
	if err != nil {
		fmt.Println("Error Unmarshaling secrets json: ", err)
	}

	Token = appSecrets.Token
	URL = appSecrets.URL
}
