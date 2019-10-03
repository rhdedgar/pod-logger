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
