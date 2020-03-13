package oapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/clam"
	"github.com/rhdedgar/pod-logger/client"
	"github.com/rhdedgar/pod-logger/cloud"
	"github.com/rhdedgar/pod-logger/docker"

	"net/http"

	"github.com/rhdedgar/pod-logger/config"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	// APIURL is the OpenShift API URL
	APIURL = config.AppSecrets.OAPIURL
	logURL = os.Getenv("LOG_WRITER_URL")
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

// PrepClamInfo gathers information about a user before sending off to the scan result bucket.
func PrepClamInfo(scanResult models.ScanResult) {
	podNs := scanResult.NameSpace
	podName := scanResult.PodName

	//fmt.Println("trying clam pod: ", podNs, podName)
	_, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		log.Println("Error getting clam pod info:", err)
	}

	scanResult.UserName = nsInfo.Metadata.Annotations.OpenshiftIoRequester
	go cloud.UploadScanLog(scanResult)
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
	sendData(mLog)
}

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
