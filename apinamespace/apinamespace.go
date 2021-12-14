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

package apinamespace

import "time"

type Spec struct {
	Finalizers []string `json:"finalizers"`
}

type PhaseStatus struct {
	Phase string `json:"phase"`
}

type APINamespace struct {
	Kind       string      `json:"kind"`
	APIVersion string      `json:"apiVersion"`
	Metadata   Metadata    `json:"metadata"`
	Spec       Spec        `json:"spec"`
	Status     PhaseStatus `json:"status"`
}

type Labels struct {
	OpenshiftIoHibernateInclude string `json:"openshift.io/hibernate-include"`
}

type Annotations struct {
	OpenshiftIoDescription             string `json:"openshift.io/description"`
	OpenshiftIoDisplayName             string `json:"openshift.io/display-name"`
	OpenshiftIoRequester               string `json:"openshift.io/requester"`
	OpenshiftIoSaSccMcs                string `json:"openshift.io/sa.scc.mcs"`
	OpenshiftIoSaSccSupplementalGroups string `json:"openshift.io/sa.scc.supplemental-groups"`
	OpenshiftIoSaSccUIDRange           string `json:"openshift.io/sa.scc.uid-range"`
}

type Metadata struct {
	Name              string      `json:"name"`
	SelfLink          string      `json:"selfLink"`
	UID               string      `json:"uid"`
	ResourceVersion   string      `json:"resourceVersion"`
	CreationTimestamp time.Time   `json:"creationTimestamp"`
	Labels            Labels      `json:"labels"`
	Annotations       Annotations `json:"annotations"`
}
