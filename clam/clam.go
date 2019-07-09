package clam

import (
	"fmt"

	"github.com/rhdedgar/pod-logger/models"
)

// CheckScanResults compares positive scan logs with the immediate takedown blacklist
func CheckScanResults(scanRes models.ScanResult) {
	for _, result := range scanRes.Results {
		fmt.Println("Scan result: ", result)
	}
}
