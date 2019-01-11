package oapi

import (
	"encoding/json"
	"fmt"
	"github.com/rhdedgar/pod-logger/apinamespace"
	"github.com/rhdedgar/pod-logger/apipod"
	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
	"io/ioutil"
	"net/http"
)

// GetInfo GETs json data from the URL designated in the config file.
func GetInfo(mStat models.Container) {

	url := config.URL
	podNs := mStat.Status.Labels.IoKubernetesPodNamespace
	podName := mStat.Status.Labels.IoKubernetesPodName

	var podDef apipod.APIPod
	var nsDef apinamespace.APINamespace

	nsUrl := fmt.Sprintf("/api/v1/namespaces/%v", podNs)
	podUrl := fmt.Sprintf("/api/v1/namespaces/%v/pods/%v/status", nsUrl, podName)

	fmt.Println("API URL: ", url)

	reqPod, err := http.NewRequest("GET", url+podUrl, nil)
	if err != nil {
		fmt.Println(err)
	}
	makeClient(reqPod, &podDef)

	reqNs, err := http.NewRequest("GET", url+nsUrl, nil)
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
	fmt.Println("mLog: \n", mLog)
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
	//	jsonStr, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
}
