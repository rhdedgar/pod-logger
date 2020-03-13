package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rhdedgar/pod-logger/config"
)

var (
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

// MakeClient takes an HTTP request and a destination struct to store the json results.
// Uses a custom HTTP transport to accommodate OpenShift clusters with self-signed certificates.
func MakeClient(req *http.Request, ds interface{}) error {
	req.Header.Set("Authorization", "Bearer "+config.AppSecrets.OAPIToken)
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

	//fmt.Println("response status: ", resp.Status)
	//fmt.Println("response Body:", string(body))

	err = json.Unmarshal(body, &ds)
	if err != nil {
		return fmt.Errorf("makeClient: Error Unmarshalling json returned from API: %v\n", err)
	}
	return nil
}
