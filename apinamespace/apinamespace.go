package apinamespace

import "time"

type Spec struct {
	Finalizers []string `json:"finalizers"`
}

type Status struct {
	Phase string `json:"phase"`
}

type APINamespace struct {
	Kind       string   `json:"kind"`
	APIVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
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
