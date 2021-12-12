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

package cloud

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rhdedgar/pod-logger/config"
	"github.com/rhdedgar/pod-logger/models"
)

func UploadScanLog(sRes models.ScanResult) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.AppSecrets.LogBucketRegion),
		Credentials: credentials.NewStaticCredentials(config.AppSecrets.LogBucketKeyID, config.AppSecrets.LogBucketKey, ""),
	})

	uploader := s3manager.NewUploader(sess)
	filename := time.Now().Format("2006-02-01") + "/" + config.ClusterUUID + "/" + sRes.UserName + "/" + sRes.PodName
	jsonStr, err := json.Marshal(sRes)
	if err != nil {
		log.Println("UploadScanLog: Error marshalling json: ", err)
	}

	file := bytes.NewReader(jsonStr)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AppSecrets.LogBucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		log.Printf("Unable to upload %q to %q, %v", filename, config.AppSecrets.LogBucketName, err)
	}
}
