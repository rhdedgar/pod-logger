package apipod

import "time"

type Labels struct {
	Component              string `json:"component"`
	ControllerRevisionHash string `json:"controller-revision-hash"`
	LoggingInfra           string `json:"logging-infra"`
	PodTemplateGeneration  string `json:"pod-template-generation"`
	Provider               string `json:"provider"`
}

type Annotations struct {
	OpenshiftIoScc                        string `json:"openshift.io/scc"`
	SchedulerAlphaKubernetesIoCriticalPod string `json:"scheduler.alpha.kubernetes.io/critical-pod"`
}

type OwnerReferences struct {
	APIVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	UID                string `json:"uid"`
	Controller         bool   `json:"controller"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
}

type Metadata struct {
	Name              string            `json:"name"`
	GenerateName      string            `json:"generateName"`
	Namespace         string            `json:"namespace"`
	SelfLink          string            `json:"selfLink"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            Labels            `json:"labels"`
	Annotations       Annotations       `json:"annotations"`
	OwnerReferences   []OwnerReferences `json:"ownerReferences"`
}

type HostPath struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type ConfigMap struct {
	Name        string `json:"name"`
	DefaultMode int    `json:"defaultMode"`
}

type Secret struct {
	SecretName  string `json:"secretName"`
	DefaultMode int    `json:"defaultMode"`
}

type Volumes struct {
	Name      string    `json:"name"`
	HostPath  HostPath  `json:"hostPath,omitempty"`
	ConfigMap ConfigMap `json:"configMap,omitempty"`
	Secret    Secret    `json:"secret,omitempty"`
}

type ResourceFieldRef struct {
	ContainerName string `json:"containerName"`
	Resource      string `json:"resource"`
	Divisor       string `json:"divisor"`
}

type ValueFrom struct {
	ResourceFieldRef ResourceFieldRef `json:"resourceFieldRef"`
}

type Env struct {
	Name      string    `json:"name"`
	Value     string    `json:"value,omitempty"`
	ValueFrom ValueFrom `json:"valueFrom,omitempty"`
}

type Limits struct {
	Memory string `json:"memory"`
}

type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type Resources struct {
	Limits   Limits   `json:"limits"`
	Requests Requests `json:"requests"`
}

type VolumeMounts struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

type SecurityContext struct {
	Privileged bool `json:"privileged"`
}

type Containers []struct {
	Name                     string          `json:"name"`
	Image                    string          `json:"image"`
	Env                      []Env           `json:"env"`
	Resources                Resources       `json:"resources"`
	VolumeMounts             []VolumeMounts  `json:"volumeMounts"`
	TerminationMessagePath   string          `json:"terminationMessagePath"`
	TerminationMessagePolicy string          `json:"terminationMessagePolicy"`
	ImagePullPolicy          string          `json:"imagePullPolicy"`
	SecurityContext          SecurityContext `json:"securityContext"`
}

type NodeSelector struct {
	LoggingInfraFluentd string `json:"logging-infra-fluentd"`
}

type ImagePullSecrets struct {
	Name string `json:"name"`
}

type Tolerations struct {
	Key      string `json:"key"`
	Operator string `json:"operator"`
	Effect   string `json:"effect"`
}

type Spec struct {
	Volumes                       []Volumes    `json:"volumes"`
	Containers                    Containers   `json:"containers"`
	RestartPolicy                 string       `json:"restartPolicy"`
	TerminationGracePeriodSeconds int          `json:"terminationGracePeriodSeconds"`
	DNSPolicy                     string       `json:"dnsPolicy"`
	NodeSelector                  NodeSelector `json:"nodeSelector"`
	ServiceAccountName            string       `json:"serviceAccountName"`
	ServiceAccount                string       `json:"serviceAccount"`
	NodeName                      string       `json:"nodeName"`
	SecurityContext               struct {
	} `json:"securityContext"`
	ImagePullSecrets  []ImagePullSecrets `json:"imagePullSecrets"`
	SchedulerName     string             `json:"schedulerName"`
	Tolerations       []Tolerations      `json:"tolerations"`
	PriorityClassName string             `json:"priorityClassName"`
	Priority          int                `json:"priority"`
}

type Conditions struct {
	Type               string      `json:"type"`
	Status             string      `json:"status"`
	LastProbeTime      interface{} `json:"lastProbeTime"`
	LastTransitionTime time.Time   `json:"lastTransitionTime"`
}

type Running struct {
	StartedAt time.Time `json:"startedAt"`
}

type State struct {
	Running Running `json:"running"`
}

type ContainerStatuses struct {
	Name      string `json:"name"`
	State     State  `json:"state"`
	LastState struct {
	} `json:"lastState"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartCount"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	ContainerID  string `json:"containerID"`
}

type Status struct {
	Phase             string              `json:"phase"`
	Conditions        []Conditions        `json:"conditions"`
	HostIP            string              `json:"hostIP"`
	PodIP             string              `json:"podIP"`
	StartTime         time.Time           `json:"startTime"`
	ContainerStatuses []ContainerStatuses `json:"containerStatuses"`
	QosClass          string              `json:"qosClass"`
}

type APIPod struct {
	Kind       string   `json:"kind"`
	APIVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
}
