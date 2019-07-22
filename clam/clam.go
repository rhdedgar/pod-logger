package clam

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
)

// CheckScanResults compares positive scan logs with the immediate takedown blacklist
func CheckScanResults(scanRes models.ScanResult) {
	for _, result := range scanRes.Results {
		fmt.Printf("Scan result: %+v", result)

		for sig, reason := range config.AppSecrets.TDSigList {
			if sig == result.Name {
				fmt.Println("calling banuser here for:", scanRes.UserName, reason)
				//banUser(scanRes.UserName, reason)
				return
			}
		}
	}
}

func banUser(userName, banReason string) {
	fmt.Println("Banning user: ", userName)

	var jsonStr = []byte(`{"authorization_username":config.AppSecrets.TDAPIUser, is_banned": "true", "takedown_code": banReason}`)
	req, err := http.NewRequest("POST", config.AppSecrets.TDAPIURL, bytes.NewBuffer(jsonStr))

	req.Header.Set("Authorization", "Bearer "+config.AppSecrets.TDAPIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("authorization_username", "application/json")

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

	fmt.Println("response Body:", string(body))
}
