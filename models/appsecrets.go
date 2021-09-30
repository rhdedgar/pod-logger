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

package models

// AppSecrets represents the secret data from the json secrets file. It's needed for OpenShift and AWS API calls
type AppSecrets struct {
	OAPIToken       string            `json:"oapi_token"`
	OAPIURL         string            `json:"oapi_url"`
	TDAPIToken      string            `json:"td_api_token"`
	TDAPIURL        string            `json:"td_api_url"`
	TDAPIUser       string            `json:"td_api_user"`
	TDSigList       map[string]string `json:"td_sig_list"`
	LogBucketKeyID  string            `json:"log_bucket_key_id"`
	LogBucketKey    string            `json:"log_bucket_key"`
	LogBucketName   string            `json:"log_bucket_name"`
	LogBucketRegion string            `json:"log_bucket_region"`
	UserWhitelist   []string          `json:"user_whitelist"`
}
