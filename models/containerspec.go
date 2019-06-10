package models

import "time"

type Log struct {
	User      string
	Namespace string
	PodName   string
	HostIP    string
	PodIP     string
	StartTime time.Time
	UID       string
}

type Metadata struct {
	Attempt int    `json:"attempt"`
	Name    string `json:"name"`
}

type Image struct {
	Image string `json:"image"`
}

type Labels struct {
	IoKubernetesContainerName string `json:"io.kubernetes.container.name"`
	IoKubernetesPodName       string `json:"io.kubernetes.pod.name"`
	IoKubernetesPodNamespace  string `json:"io.kubernetes.pod.namespace"`
	IoKubernetesPodUID        string `json:"io.kubernetes.pod.uid"`
}

type Annotations struct {
	IoKubernetesContainerHash                     string `json:"io.kubernetes.container.hash"`
	IoKubernetesContainerPorts                    string `json:"io.kubernetes.container.ports"`
	IoKubernetesContainerRestartCount             string `json:"io.kubernetes.container.restartCount"`
	IoKubernetesContainerTerminationMessagePath   string `json:"io.kubernetes.container.terminationMessagePath"`
	IoKubernetesContainerTerminationMessagePolicy string `json:"io.kubernetes.container.terminationMessagePolicy"`
	IoKubernetesPodTerminationGracePeriod         string `json:"io.kubernetes.pod.terminationGracePeriod"`
}

type Mounts struct {
	ContainerPath  string `json:"containerPath"`
	HostPath       string `json:"hostPath"`
	Propagation    string `json:"propagation"`
	Readonly       bool   `json:"readonly"`
	SelinuxRelabel bool   `json:"selinuxRelabel"`
}

type Status struct {
	ID          string      `json:"id"`
	Metadata    Metadata    `json:"metadata"`
	State       string      `json:"state"`
	CreatedAt   time.Time   `json:"createdAt"`
	StartedAt   time.Time   `json:"startedAt"`
	FinishedAt  time.Time   `json:"finishedAt"`
	ExitCode    int         `json:"exitCode"`
	Image       Image       `json:"image"`
	ImageRef    string      `json:"imageRef"`
	Reason      string      `json:"reason"`
	Message     string      `json:"message"`
	Labels      Labels      `json:"labels"`
	Annotations Annotations `json:"annotations"`
	Mounts      []Mounts    `json:"mounts"`
	LogPath     string      `json:"logPath"`
}

type Container struct {
	Status Status `json:"status"`
}
