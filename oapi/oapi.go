package oapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/docker"

	"io/ioutil"
	"net/http"

	"github.com/rhdedgar/pod-logger/config"

	"github.com/rhdedgar/pod-logger/models"
)

var (
	// APIURL is the OpenShift API URL
	APIURL           = config.APIURL
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
	getInfo(podNs, podName)
}

func PrepCrioInfo(mStat models.Container) {
	podNs := mStat.Status.Labels.IoKubernetesPodNamespace
	podName := mStat.Status.Labels.IoKubernetesPodName
	fmt.Println("trying crio pod: ", podNs, podName)
	getInfo(podNs, podName)
}

// GetInfo GETs json data from the URL designated in the config package.
func getInfo(podNs, podName string) {

	var podDef apipod.APIPod
	var nsDef apinamespace.APINamespace

	nsURL := fmt.Sprintf("/api/v1/namespaces/%v", podNs)
	podURL := fmt.Sprintf("/api/v1/namespaces/%v/pods/%v/status", podNs, podName)

	// Marshall the pod response from the API server into the podDef struct
	reqPod, err := http.NewRequest("GET", APIURL+podURL, nil)
	if err != nil {
		fmt.Println("Error getting pod info: ", err)
	}
	makeClient(reqPod, &podDef)

	// Marshall the namespace response from the API server into the nsDef struct
	reqNs, err := http.NewRequest("GET", APIURL+nsURL, nil)
	if err != nil {
		fmt.Println("Error getting pod info: ", err)
	}
	makeClient(reqNs, &nsDef)

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
	logURL := os.Getenv("LOG_WRITER_URL")

	jsonStr, err := json.Marshal(mlog)
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
	}

	req, err := http.NewRequest("POST", logURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("sendData: Error making request: ", err)
	}
	defer resp.Body.Close()

	// TODO Prometheus to check header response
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func makeClient(req *http.Request, ds interface{}) {
	token := config.Token

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Transport: httpClientWithSelfSignedTLS}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("makeClient: Error making API request: ", err)
	}

	defer resp.Body.Close()

	// TODO Prometheus to check header response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
	}

	fmt.Println("response Body:", string(body))

	err = json.Unmarshal(body, &ds)
	if err != nil {
		fmt.Println("Error Unmarshalling json returned from API: \n", err)
		fmt.Println(ds)
	}
}
