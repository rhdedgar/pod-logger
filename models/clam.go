package models

import "time"

type ScanResult struct {
	APIVersion  string    `json:"apiVersion"`
	ContainerID string    `json:"containerID"`
	ImageID     string    `json:"imageID"`
	ImageName   string    `json:"imageName"`
	NameSpace   string    `json:"nameSpace"`
	PodName     string    `json:"podName"`
	Results     []Results `json:"results"`
	UserName    string    `json:"userName"`
}

type Results struct {
	Description    string    `json:"description"`
	Name           string    `json:"name"`
	Reference      string    `json:"reference"`
	ScannerVersion string    `json:"scannerVersion"`
	Timestamp      time.Time `json:"timestamp"`
}
