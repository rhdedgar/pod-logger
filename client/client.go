/*
Copyright 2019 Doug Edgar.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	defaultTransport = http.DefaultTransport.(*http.Transport)
	// Create new Transport that accommodates self-signed SSL
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

// MakeClient takes an HTTP request, a destination struct to store the json results, and an access token.
// Uses a custom HTTP transport to accommodate OpenShift clusters with self-signed certificates.
func MakeClient(req *http.Request, ds interface{}, tok string) (int, error) {
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Transport: httpClientWithSelfSignedTLS}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("makeClient: Error making API request: %v", err)
	}
	defer resp.Body.Close()

	// TODO Prometheus to check header response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("makeClient: Error reading response body: %v", err)
	}

	//fmt.Println("response status: ", resp.Status)
	//fmt.Println("response Body:", string(body))

	err = json.Unmarshal(body, &ds)
	if err != nil {
		return 0, fmt.Errorf("makeClient: Error Unmarshalling json returned from API: %v\n", err)
	}
	return resp.StatusCode, nil
}
