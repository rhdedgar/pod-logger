package oapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/clam"
	"github.com/rhdedgar/pod-logger/cloud"
	"github.com/rhdedgar/pod-logger/docker"

	"io/ioutil"
	"net/http"

	"github.com/rhdedgar/pod-logger/config"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	// APIURL is the OpenShift API URL
	APIURL           = config.AppSecrets.OAPIURL
	token            = config.AppSecrets.OAPIToken
	logURL           = os.Getenv("LOG_WRITER_URL")
	defaultTransport = http.DefaultTransport.(*http.Transport)
	// Create new Transport that ignores self-signed SSL
	httpClientWithSelfSignedTLS = &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
)

func PrepDockerInfo(mStat docker.DockerContainer) {
	podNs := mStat[0].Config.Labels.IoKubernetesPodNamespace
	podName := mStat[0].Config.Labels.IoKubernetesPodName
	fmt.Println("trying docker pod: ", podNs, podName)
	podInfo, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		fmt.Println("Error getting Docker pod info:", err)
	}

	prepLog(podName, podNs, podInfo, nsInfo)
}

func PrepCrioInfo(mStat models.Container) {
	podNs := mStat.Status.Labels.IoKubernetesPodNamespace
	podName := mStat.Status.Labels.IoKubernetesPodName
	fmt.Println("trying crio pod with URL: ", APIURL, podNs, podName)
	podInfo, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		fmt.Println("Error getting Cri-o pod info:", err)
	}

	prepLog(podName, podNs, podInfo, nsInfo)
}

func PrepClamInfo(scanResult models.ScanResult) {
	podNs := scanResult.NameSpace
	podName := scanResult.PodName
	fmt.Println("trying clam pod: ", podNs, podName)
	_, nsInfo, err := getInfo(podNs, podName)

	if err != nil {
		fmt.Println("Error getting clam pod info:", err)
	}

	scanResult.UserName = nsInfo.Metadata.Annotations.OpenshiftIoRequester
	go clam.CheckScanResults(scanResult)
	cloud.UploadScanLog(scanResult)
}

//func uploadScanLog(sLog models.ScanResult) {
//	fmt.Println(sLog)
//}

// getInfo GETs json data from the URL designated in the config package.
func getInfo(podNs, podName string) (apipod.APIPod, apinamespace.APINamespace, error) {
	var podDef apipod.APIPod
	var nsDef apinamespace.APINamespace

	nsURL := fmt.Sprintf("/api/v1/namespaces/%v", podNs)
	podURL := fmt.Sprintf("/api/v1/namespaces/%v/pods/%v/status", podNs, podName)

	fmt.Println("provided:", podNs, podName)
	// Marshall the pod response from the API server into the podDef struct
	reqPod, err := http.NewRequest("GET", APIURL+podURL, nil)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error getting pod info: %v \n", err)
	}

	err = makeClient(reqPod, &podDef)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error making pod request client: %v \n", err)
	}

	// Marshall the namespace response from the API server into the nsDef struct
	reqNs, err := http.NewRequest("GET", APIURL+nsURL, nil)
	if err != nil {
		return podDef, nsDef, fmt.Errorf("Error getting pod info: %v \n", err)
	}

	err = makeClient(reqNs, &nsDef)
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
	//fmt.Printf("mLog: \n %+v", mLog)
	sendData(mLog)
}

func sendData(mlog models.Log) {
	jsonStr, err := json.Marshal(mlog)
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
		return
	}

	req, err := http.NewRequest("POST", logURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("sendData: Error making request: ", err)
		return
	}
	defer resp.Body.Close()
}

func makeClient(req *http.Request, ds interface{}) error {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Transport: httpClientWithSelfSignedTLS}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("makeClient: Error making API request: %v", err)
	}

	defer resp.Body.Close()

	// TODO Prometheus to check header response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("makeClient: Error reading response body: %v", err)
	}

	//fmt.Println("response Body:", string(body))

	err = json.Unmarshal(body, &ds)
	if err != nil {
		return fmt.Errorf("makeClient: Error Unmarshalling json returned from API: %v\n", err)
	}
	return nil
}
