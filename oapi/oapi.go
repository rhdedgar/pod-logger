package oapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"

	"io/ioutil"
	"net/http"

	"github.com/rhdedgar/pod-logger/config"

	"github.com/rhdedgar/pod-logger/models"
)

// GetInfo GETs json data from the URL designated in the config file.
func GetInfo(mStat models.Container) {
	url := config.URL
	podNs := mStat.Status.Labels.IoKubernetesPodNamespace
	podName := mStat.Status.Labels.IoKubernetesPodName

	fmt.Println("trying: ", podNs, podName)

	var podDef apipod.APIPod
	var nsDef apinamespace.APINamespace

	nsURL := fmt.Sprintf("/api/v1/namespaces/%v", podNs)
	podURL := fmt.Sprintf("/api/v1/namespaces/%v/pods/%v/status", podNs, podName)

	fmt.Println("API URL: ", url)

	reqPod, err := http.NewRequest("GET", url+podURL, nil)
	if err != nil {
		fmt.Println(err)
	}
	makeClient(reqPod, &podDef)

	reqNs, err := http.NewRequest("GET", url+nsURL, nil)
	if err != nil {
		fmt.Println(err)
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
	url := os.Getenv("LOG_WRITER_URL")

	jsonStr, err := json.Marshal(mlog)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// TODO Prometheus to check header response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Body:", string(body))

	err = json.Unmarshal(body, &ds)
	if err != nil {
		fmt.Println("Error Unmarshalling json returned from API: \n", err)
		fmt.Println(ds)
	}
}
