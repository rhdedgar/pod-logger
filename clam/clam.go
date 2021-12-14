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

/*
Package clam compares positive AV signtures against a blacklist.
Various levels of actions for blacklist matches are available (immediate ban, project deletion, etc.).
*/
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
				if strings.HasPrefix(config.ClusterUUID, "starter") { // deprecated, needs replacing if needed again, as UUID will never eval to ^starter-
					BanUser(scanRes.UserName, reason)
					return
				}
				DeleteNS(scanRes.NameSpace)
				return
			}
		}
	}
}

// BanUser bans the provided user via takedown API call for the specified reason.
func BanUser(userName, banReason string) (int, error) {
	for _, excluded := range config.AppSecrets.UserWhitelist {
		if userName == excluded {
			//fmt.Printf("NOT banning user %q\n", excluded)
			return 200, nil
		}
	}

	log.Println("Banning user: ", userName)
	var newBan = models.BanAPICall{AuthUser: config.AppSecrets.TDAPIUser, IsBanned: "true", TakedownCode: banReason}

	jsonStr, err := json.Marshal(newBan)
	if err != nil {
		log.Println("Error marshalling banUser json: ", err)
		return 0, err
	}

	req, err := http.NewRequest("PUT", config.AppSecrets.TDAPIURL+userName+"/ban", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("Error creating request to ban user: ", err)
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+config.AppSecrets.TDAPIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization_username", config.AppSecrets.TDAPIURL+userName)

	//fmt.Println("BanRequest", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("banUser: Error making API request: ")
		return 0, err
	}

	defer resp.Body.Close()

	log.Println("Successfully called ban API: ", resp.Status)
	return resp.StatusCode, nil
}

// DeleteNS takes a string name of the namespace to delete via OpenShift API call.
func DeleteNS(ns string) (int, error) {
	recJSON := make(map[string]interface{})
	log.Println("Deleting namespace: ", ns)

	req, err := http.NewRequest("DELETE", config.AppSecrets.OAPIURL+"/apis/project.openshift.io/v1/projects/"+ns, nil)

	if err != nil {
		log.Println("Error creating request to delete namespace: ", err)
		return 0, err
	}

	status, err := client.MakeClient(req, &recJSON, config.AppSecrets.OAPIToken)
	if err != nil {
		log.Printf("Error making delete request client: %v \n", err)
		return 0, err
	}
	return status, nil
}
