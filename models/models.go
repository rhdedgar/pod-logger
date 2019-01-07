package models

type Label struct {
	IoKubernetesContainerName string `json:"io.kubernetes.container.name" yaml:"io.kubernetes.container.name"`
	IoKubernetesPodName       string `json:"io.kubernetes.pod.name" yaml:"io.kubernetes.pod.name"`
	IoKubernetesPodNamespace  string `json:"io.kubernetes.pod.namespace" yaml:"io.kubernetes.pod.namespace"`
	IoKubernetesPodUID        string `json:"io.kubernetes.pod.uid" yaml:"io.kubernetes.pod.uid"`
}

type Pod struct {
	Labels Label `json:"Labels,omitempty" yaml:"Labels,omitempty"`
}
