package clam

import (
	"fmt"

	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
)

// CheckScanResults compares positive scan logs with the immediate takedown blacklist
func CheckScanResults(scanRes models.ScanResult) {
	for _, result := range scanRes.Results {
		fmt.Printf("Scan result: %+v", result)

		for _, sig := range config.AppSecrets.TDSigList {
			if sig == result.Name {
				banUser(scanRes.UserName)
				return
			}
		}
	}
}

func banUser(userName string) {
	fmt.Println("Banning user: ", userName)
}
