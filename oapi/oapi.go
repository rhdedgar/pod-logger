package oapi

import (
	"encoding/json"
	"fmt"
	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
	"io/ioutil"
	"net/http"
)

// SendData Marshals and POSTs json data to the URL designated in the config file.
func SendData(mStat models.Status) {
	i
	var mLog models.Log
	url := config.URL
	token := config.Token

	fmt.Println("API URL: ", url)

	jsonStr, err := json.Marshal(mStat)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest("GET", url, nil)
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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
