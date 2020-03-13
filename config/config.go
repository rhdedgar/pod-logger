package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	// AppSecrets is the populated struct of secrets needed for OpenShift and AWS API auth.
	AppSecrets models.AppSecrets
	// ClusterName is the name of the current OpenShift cluster. Used in clamav logs
	ClusterName = os.Getenv("CLUSTER_NAME")
)

func loadJSON(filePath string) {
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Println("Error loading secrets json: ", err)
	}

	//fmt.Println("Config file contents: ", string(fileBytes))

	err = json.Unmarshal(fileBytes, &AppSecrets)
	if err != nil {
		log.Println("Error Unmarshalling secrets json: ", err)
	}
}

// init attempts to populate the AppSecrets var with data needed to run this server.
func init() {
	filePath := "/secrets/api_config.json"
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	oAPIURL := "https://openshift.default.svc.cluster.local"

	loadJSON(filePath)
	fileBytes, err := ioutil.ReadFile(tokenPath)

	if err != nil {
		log.Println("Error loading service account token file: ", err)
	}

	AppSecrets.OAPIURL = oAPIURL
	AppSecrets.OAPIToken = string(fileBytes)

	if AppSecrets.OAPIToken == "" {
		log.Println("Secrets were not loaded, application will fail.")
	}
}
