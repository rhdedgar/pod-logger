package clam

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/rhdedgar/pod-logger/client"
	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
)

// CheckScanResults compares positive scan logs with the immediate takedown blacklist
func CheckScanResults(scanRes models.ScanResult) {
	log.Printf("%+v", scanRes)
	for _, result := range scanRes.Results {
		//fmt.Printf("Scan result: %+v", result)

		for sig, reason := range config.AppSecrets.TDSigList {
			//fmt.Printf("comparing: %v\n to %v\n", sig, result.Description)
			if sig == strings.TrimSuffix(result.Description, " FOUND") {
				log.Printf("User %v matched blacklist for %v:", scanRes.UserName, reason)
				if strings.HasPrefix(config.ClusterName, "starter") {
					banUser(scanRes.UserName, reason)
					return
				}
				deleteNS(scanRes.NameSpace)
				return
			}
		}
	}
}

func banUser(userName, banReason string) {
	for _, excluded := range config.AppSecrets.UserWhitelist {
		if userName == excluded {
			//fmt.Printf("NOT banning user %q\n", excluded)
			return
		}
	}

	log.Println("Banning user: ", userName)
	var newBan = models.BanAPICall{AuthUser: config.AppSecrets.TDAPIUser, IsBanned: "true", TakedownCode: banReason}

	jsonStr, err := json.Marshal(newBan)
	if err != nil {
		log.Println("Error marshalling banUser json: ", err)
		return
	}

	req, err := http.NewRequest("PUT", config.AppSecrets.TDAPIURL+userName+"/ban", bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "Bearer "+config.AppSecrets.TDAPIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization_username", config.AppSecrets.TDAPIURL+userName)

	//fmt.Println("BanRequest", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("banUser: Error making API request: ")
	}

	defer resp.Body.Close()

	// TODO Prometheus to check header response
	log.Println("Successfully called ban API: ", resp.Status)
}

func deleteNS(ns string) {
	recJSON := make(map[string]interface{})
	log.Println("Deleting namespace: ", ns)

	req, err := http.NewRequest("DELETE", config.AppSecrets.OAPIURL+"/apis/project.openshift.io/v1/projects/"+ns, nil)

	if err != nil {
		log.Println("Error creating request to delete namespace: ", err)
	}

	err = client.MakeClient(req, &recJSON)
	if err != nil {
		log.Printf("Error making delete request client: %v \n", err)
	}
}
