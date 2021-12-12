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

type ClusterVersion struct {
	APIVersion string   `json:"apiVersion"`
	Items      []Items  `json:"items"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
}
type NAMING_FAILED struct {
}
type FChannel struct {
}
type FClusterID struct {
}
type FSpec struct {
	NAMING_FAILED NAMING_FAILED `json:"."`
	FChannel      FChannel      `json:"f:channel"`
	FClusterID    FClusterID    `json:"f:clusterID"`
}

type FAvailableUpdates struct {
}
type FConditions struct {
}
type FChannels struct {
}
type FImage struct {
}
type FURL struct {
}
type FVersion struct {
}
type FDesired struct {
	NAMING_FAILED NAMING_FAILED `json:"."`
	FChannels     FChannels     `json:"f:channels"`
	FImage        FImage        `json:"f:image"`
	FURL          FURL          `json:"f:url"`
	FVersion      FVersion      `json:"f:version"`
}
type FHistory struct {
}
type FObservedGeneration struct {
}
type FVersionHash struct {
}
type FStatus struct {
	NAMING_FAILED       NAMING_FAILED       `json:"."`
	FAvailableUpdates   FAvailableUpdates   `json:"f:availableUpdates"`
	FConditions         FConditions         `json:"f:conditions"`
	FDesired            FDesired            `json:"f:desired"`
	FHistory            FHistory            `json:"f:history"`
	FObservedGeneration FObservedGeneration `json:"f:observedGeneration"`
	FVersionHash        FVersionHash        `json:"f:versionHash"`
}
type FieldsV1 struct {
	FStatus FStatus `json:"f:status"`
	FSpec   FSpec   `json:"f:spec"`
}
type ManagedFields struct {
	APIVersion  string    `json:"apiVersion"`
	FieldsType  string    `json:"fieldsType"`
	FieldsV1    FieldsV1  `json:"fieldsV1,omitempty"`
	Manager     string    `json:"manager"`
	Operation   string    `json:"operation"`
	Time        time.Time `json:"time"`
	Subresource string    `json:"subresource,omitempty"`
}
type Metadata struct {
	CreationTimestamp time.Time       `json:"creationTimestamp"`
	Generation        int             `json:"generation"`
	ManagedFields     []ManagedFields `json:"managedFields"`
	Name              string          `json:"name"`
	ResourceVersion   string          `json:"resourceVersion"`
	UID               string          `json:"uid"`
	Continue          string          `json:"continue"`
}
type Spec struct {
	Channel   string `json:"channel"`
	ClusterID string `json:"clusterID"`
}
type Conditions struct {
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Message            string    `json:"message,omitempty"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}
type Desired struct {
	Channels []string `json:"channels"`
	Image    string   `json:"image"`
	URL      string   `json:"url"`
	Version  string   `json:"version"`
}
type History struct {
	CompletionTime time.Time `json:"completionTime"`
	Image          string    `json:"image"`
	StartedTime    time.Time `json:"startedTime"`
	State          string    `json:"state"`
	Verified       bool      `json:"verified"`
	Version        string    `json:"version"`
}
type Status struct {
	AvailableUpdates   interface{}  `json:"availableUpdates"`
	Conditions         []Conditions `json:"conditions"`
	Desired            Desired      `json:"desired"`
	History            []History    `json:"history"`
	ObservedGeneration int          `json:"observedGeneration"`
	VersionHash        string       `json:"versionHash"`
}
type Items struct {
	APIVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
}
