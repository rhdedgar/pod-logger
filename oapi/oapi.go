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
	// APIURL is the OpenShift API URL.
	APIURL = config.AppSecrets.OAPIURL
	// logURL is th URL used by sendData to POST scan.Logs as JSON.
	logURL = os.Getenv("LOG_WRITER_URL")
	// scanLog is the file path that splunk-forwarder-operator is configured to read.
	scanLog = os.Getenv("SCAN_LOG_FILE")
)

// PrepDockerInfo gathers information about a user before sending it off to the logging service.
func PrepDockerInfo(mStat docker.DockerContainer) {
	podNs := mStat[0].Config.Labels.IoKubernetesPodNamespace
	podName := mStat[0].Config.Labels.IoKubernetesPodName

	//fmt.Println("trying docker pod: ", podNs, podName)
	podInfo, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting Docker pod info:", err)
	}

	prepLog(podName, podNs, podInfo, nsInfo)
}

// PrepCrioInfo gathers information about a user before sending it off to the logging service.
func PrepCrioInfo(mStat models.Container) {
	podNs := mStat.Status.Labels.IoKubernetesPodNamespace
	podName := mStat.Status.Labels.IoKubernetesPodName

	//fmt.Println("trying crio pod with URL: ", APIURL, podNs, podName)
	podInfo, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting Cri-o pod info:", err)
	}

	prepLog(podName, podNs, podInfo, nsInfo)
}

// PrepClamInfo gathers information about a user before logging the scan result.
func PrepClamInfo(scanResult models.ScanResult, mx *sync.Mutex) {
	podNs := scanResult.NameSpace
	podName := scanResult.PodName

	_, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting clam pod info:", err)
	}

	scanResult.UserName = nsInfo.Metadata.Annotations.OpenshiftIoRequester
	//go cloud.UploadScanLog(scanResult)

	mx.Lock()
	defer mx.Unlock()

	f, err := os.OpenFile(scanLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0655)
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

// getInfo GETs json data from the URL designated in the config package.
func getInfo(podNs, podName string) (apipod.APIPod, apinamespace.APINamespace, error) {
	var podDef apipod.APIPod
	var nsDef apinamespace.APINamespace

	nsURL := fmt.Sprintf("/api/v1/namespaces/%v", podNs)
	podURL := fmt.Sprintf("/api/v1/namespaces/%v/pods/%v/status", podNs, podName)

	//fmt.Println("provided:", podNs, podName)

	// Marshall the pod response from the API server into the podDef struct
	reqPod, err := http.NewRequest("GET", APIURL+podURL, nil)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error getting pod info: %v \n", err)
	}

	err = client.MakeClient(reqPod, &podDef)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error making pod request client: %v \n", err)
	}

	// Marshall the namespace response from the API server into the nsDef struct
	reqNs, err := http.NewRequest("GET", APIURL+nsURL, nil)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error getting pod info: %v \n", err)
	}

	err = client.MakeClient(reqNs, &nsDef)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error making NS request client: %v\n", err)
	}

	return podDef, nsDef, nil
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
}

// sendData POSTS a models.Log as JSON to logURL
func sendData(mlog models.Log) {
	jsonStr, err := json.Marshal(mlog)
	if err != nil {
		log.Println("Error marshalling json: ", err)
		return
	}

	req, err := http.NewRequest("POST", logURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("sendData: Error making request: ", err)
		return
	}
	defer resp.Body.Close()
}
