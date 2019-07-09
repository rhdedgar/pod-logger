package cloud

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	filename := time.Now().Format("2006-02-01") + "/" + config.ClusterName + "/" + sRes.UserName + "/" + sRes.PodName
	jsonStr, err := json.Marshal(sRes)
	if err != nil {
		fmt.Println("UploadScanLog: Error marshalling json: ", err)
	}

	file := bytes.NewReader(jsonStr)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AppSecrets.LogBucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		fmt.Printf("Unable to upload %q to %q, %v", filename, config.AppSecrets.LogBucketName, err)
	}
}
