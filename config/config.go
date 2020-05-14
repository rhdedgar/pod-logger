package config

import (
	"encoding/json"
	"fmt"
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
	// LogURL is the URL used by sendData to POST scan.Logs as JSON.
	LogURL = os.Getenv("LOG_WRITER_URL")
	// ScanLogFile is the scan log file path that splunk-forwarder-operator is configured to read.
	ScanLogFile = os.Getenv("SCAN_LOG_FILE")
	// PodLogFile is the pod creation log file path that splunk-forwarder-operator is configured to read.
	PodLogFile = os.Getenv("POD_LOG_FILE")
)

// LoadJSON reads a file at filePath, and Unmarshals the contents into the provided data structure pointer.
func LoadJSON(ds interface{}, filePath string) error {
	fileBytes, err := ioutil.ReadFile(filePath)
	fmt.Println("Config file contents: ", string(fileBytes))

	if err != nil {
		log.Println("Error loading secrets json: ", err)
		return err
	}

	err = json.Unmarshal(fileBytes, ds)
	if err != nil {
		log.Println("Error Unmarshalling secrets json: ", err)
		return err
	}
	return nil
}

// init attempts to populate the AppSecrets var with data needed to run this server.
func init() {
	filePath := "/secrets/api_config.json"
	tokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"
	oAPIURL := "https://openshift.default.svc.cluster.local"

	err := LoadJSON(AppSecrets, filePath)
	if err != nil {
		log.Println("Cannot load AppSecrets JSON:")
	}

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
