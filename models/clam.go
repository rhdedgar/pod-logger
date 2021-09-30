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

type BanAPICall struct {
	AuthUser     string `json:"authorization_username"`
	IsBanned     string `json:"is_banned"`
	TakedownCode string `json:"takedown_code"`
}
