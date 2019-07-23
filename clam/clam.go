package clam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
)

// CheckScanResults compares positive scan logs with the immediate takedown blacklist
func CheckScanResults(scanRes models.ScanResult) {
	for _, result := range scanRes.Results {
		fmt.Printf("Scan result: %+v", result)

		for sig, reason := range config.AppSecrets.TDSigList {
			//fmt.Printf("comparing: %v\n to %v\n", sig, result.Description)
			if sig == strings.TrimSuffix(result.Description, " FOUND") {
				fmt.Println("calling banuser here for:", scanRes.UserName, reason)
				banUser(scanRes.UserName, reason)
				return
			}
		}
	}
}

func banUser(userName, banReason string) {
	fmt.Println("Banning user: ", userName)
	var newBan = models.BanAPICall{AuthUser: config.AppSecrets.TDAPIUser, IsBanned: "true", TakedownCode: banReason}

	jsonStr, err := json.Marshal(newBan)
	if err != nil {
		fmt.Println("Error marshalling banUser json: ", err)
		return
	}

	req, err := http.NewRequest("POST", config.AppSecrets.TDAPIURL+userName+"/ban", bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "Bearer "+config.AppSecrets.TDAPIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization_username", "application/json")

	fmt.Println("BanRequest", req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("banUser: Error making API request: ")
	}

	defer resp.Body.Close()

	// TODO Prometheus to check header response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("banUser: Error reading response body: ")
	}

	fmt.Println("Successfully called ban API")
	fmt.Println("response Body:", string(body))
}
