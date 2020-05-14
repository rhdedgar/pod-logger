/*
Package oapi provides functions to gather, format, and log information
about new pods and scan results.
*/
package oapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/clam"
	"github.com/rhdedgar/pod-logger/client"
	"github.com/rhdedgar/pod-logger/docker"

	"net/http"

	"github.com/rhdedgar/pod-logger/config"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	// config.AppSecrets.OAPIURL is the OpenShift API URL.
	//APIURL = config.AppSecrets.Oconfig.AppSecrets.OAPIURL
	// config.LogURL is the URL used by sendData to POST scan.Logs as JSON.
	//logURL = config.config.LogURL
	// config.ScanLogFile is the scan log file path that splunk-forwarder-operator is configured to read.
	//scanLogFile = config.config.ScanLogFile
	// config.PodLogFile is the pod creation log file path that splunk-forwarder-operator is configured to read.
	//podLogFile = config.config.PodLogFile
	// scanLogMX and podLogMX are the mutexes for the scan and pod creation log files.
	scanLogMX, podLogMX sync.Mutex
)

// PrepDockerInfo gathers information about a user before sending it off to the logging service.
func PrepDockerInfo(mStat docker.DockerContainer) {
	podNs := mStat[0].Config.Labels.IoKubernetesPodNamespace
	podName := mStat[0].Config.Labels.IoKubernetesPodName

	//fmt.Println("trying docker pod: ", podNs, podName)
	nsInfo, podInfo, err := GetInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting Docker pod info:", err)
	}

	prepLog(podName, podNs, podInfo, nsInfo)
}

// PrepCrioInfo gathers information about a user before sending it off to the logging service.
func PrepCrioInfo(mStat models.Container) {
	podNs := mStat.Status.Labels.IoKubernetesPodNamespace
	podName := mStat.Status.Labels.IoKubernetesPodName

	//fmt.Println("trying crio pod with URL: ", config.AppSecrets.OAPIURL, podNs, podName)
	nsInfo, podInfo, err := GetInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting Cri-o pod info:", err)
	}

	prepLog(podName, podNs, podInfo, nsInfo)
}

// PrepClamInfo gathers information about a user before logging the scan result.
func PrepClamInfo(scanResult models.ScanResult) {
	podNs := scanResult.NameSpace
	podName := scanResult.PodName

	nsInfo, _, err := GetInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting clam pod info:", err)
	}

	scanResult.UserName = nsInfo.Metadata.Annotations.OpenshiftIoRequester
	//go cloud.UploadScanLog(scanResult)

	scanLogMX.Lock()
	defer scanLogMX.Unlock()

	f, err := os.OpenFile(config.ScanLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0655)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	jBytes, err := json.Marshal(scanResult)
	if err != nil {
		log.Println("Error marshaling scan result to write to disk:", err)
	}

	f.Write(jBytes)
	f.WriteString("\n")

	clam.CheckScanResults(scanResult)
}

// GetInfo GETs json data from the URL designated in the config package.
func GetInfo(podNs, podName string) (apinamespace.APINamespace, apipod.APIPod, error) {
	var podDef apipod.APIPod
	var nsDef apinamespace.APINamespace

	nsURL := fmt.Sprintf("/api/v1/namespaces/%v", podNs)
	podURL := fmt.Sprintf("/api/v1/namespaces/%v/pods/%v/status", podNs, podName)

	//fmt.Println("provided:", podNs, podName)

	// Marshal the pod response from the API server into the podDef struct
	reqPod, err := http.NewRequest("GET", config.AppSecrets.OAPIURL+podURL, nil)
	if err != nil {
		return nsDef, podDef, fmt.Errorf("Error getting pod info: %v \n", err)
	}

	_, err = client.MakeClient(reqPod, &podDef)
	if err != nil {
		return nsDef, podDef, fmt.Errorf("Error making pod request client: %v \n", err)
	}

	// Marshall the namespace response from the API server into the nsDef struct
	reqNs, err := http.NewRequest("GET", config.AppSecrets.OAPIURL+nsURL, nil)
	if err != nil {
		return nsDef, podDef, fmt.Errorf("Error getting pod info: %v \n", err)
	}

	_, err = client.MakeClient(reqNs, &nsDef)
	if err != nil {
		return nsDef, podDef, fmt.Errorf("Error making NS request client: %v\n", err)
	}

	return nsDef, podDef, nil
}

// prepLog neatly formats relevant pod, project, and user data as a models.Log.
// Then prints it for splunk pickup.
func prepLog(podName, podNs string, podDef apipod.APIPod, nsDef apinamespace.APINamespace) {
	mLog := models.Log{
		User:      nsDef.Metadata.Annotations.OpenshiftIoRequester,
		Namespace: podNs,
		PodName:   podName,
		HostIP:    podDef.Status.HostIP,
		PodIP:     podDef.Status.PodIP,
		StartTime: podDef.Status.StartTime,
		UID:       nsDef.Metadata.UID,
	}
	log.Printf("%+v", mLog)

	podLogMX.Lock()
	defer podLogMX.Unlock()

	f, err := os.OpenFile(config.PodLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0655)
	if err != nil {
		log.Printf("error opening file: %v", err)
	}
	defer f.Close()

	jBytes, err := json.Marshal(mLog)
	if err != nil {
		log.Println("Error marshaling pod log to write to disk:", err)
	}

	f.Write(jBytes)
	f.WriteString("\n")
}

// SendData POSTS a models.Log as JSON to config.LogURL
func SendData(mlog models.Log) (int, error) {
	jsonStr, err := json.Marshal(mlog)
	if err != nil {
		log.Println("SendData: Error marshalling json: ", err)
		return 0, err
	}

	req, err := http.NewRequest("POST", config.LogURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println("SendData: Error making HTTP request:", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("SendData: Error making request: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
