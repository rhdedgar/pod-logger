package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	// AppSecrets is the populated struct of secrets needed for OpenShift and AWS API auth.
	AppSecrets models.AppSecrets
	// ClusterName is the name of the current OpenShift cluster. Used in clamav logs
	ClusterName = os.Getenv("CLUSTER_NAME")
)

// init attempts to populate the AppSecrets var with data needed to run this server.
func init() {
	filePath := "/secrets/api_config.json"
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	fmt.Println("Config file contents: ", string(fileBytes))

	err = json.Unmarshal(fileBytes, &AppSecrets)
	if err != nil {
		fmt.Println("Error Unmarshalling secrets json: ", err)
	}

	if AppSecrets.OAPIToken == "" {
		fmt.Println("Secrets were not loaded, application will fail.")
	}
}
