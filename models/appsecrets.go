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

// AppSecrets represents the secret data from the json secrets file. It's needed for various OpenShift and Cloud Provider API calls
type AppSecrets struct {
	// OAPIToken is the token used with the OpenShift API to gather metadata on a pod.
	// Usually obtained from /var/run/secrets/kubernetes.io/serviceaccount/token
	OAPIToken string `json:"oapi_token"`
	// OAPIURL is the OpenShift API URL used for the cluster.
	OAPIURL string `json:"oapi_url"`
	// TDAPIToken is the token to be used with the optional TakeDown API.
	TDAPIToken string `json:"td_api_token"`
	// TDAPIURL is the API URL endpoint of the optional TakeDown API.
	TDAPIURL string `json:"td_api_url"`
	// TDAPIUser is the username to be used with the optional TakeDown API.
	TDAPIUser string `json:"td_api_user"`
	// TDSigList is a curated list of signatures that warrant immediate takedown.
	// Usually used with custom signatures in which we have a high degree of confidence to only match malicious code.
	TDSigList map[string]string `json:"td_sig_list"`
	// LogBucketKeyID is the cloud provider key ID for the signature storage medium.
	LogBucketKeyID string `json:"log_bucket_key_id"`
	// LogBucketKey is the cloud provider secret key for the signature storage medium.
	LogBucketKey string `json:"log_bucket_key"`
	// LogBucketName is the name of the cloud provider signature storage medium.
	LogBucketName string `json:"log_bucket_name"`
	// LogBucketRegion is the region of the cloud provider signature storage medium.
	LogBucketRegion string `json:"log_bucket_region"`
	// UserWhitelist is a list of users whose pods are exempted from scans.
	UserWhitelist []string `json:"user_whitelist"`
	// ClusterID matches the clusterID field of our clusterversion object. Useful for identification in multi-cluster setups.
	// oc get clusterversion -o jsonpath='{.items[].spec.clusterID}{"\n"}'
	ClusterUUID string `json:"cluster_uuid"`
	// DynamoDBKeyID is the IAM API key ID with write access to our log storage table.
	DynamoDBKeyID string `json:"dynamodb_user"`
	// DynamoDBUser is the IAM API secret key with write access to our log storage table.
	DynamoDBKey string `json:"dynamodb_key"`
	// DynamoDBTable is the storage table for pod creation logs and scan logs.
	DynamoDBTable string `json:"dynamodb_table"`
}
