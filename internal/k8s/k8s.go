package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

// Config is a callback that returns a Kubernetes Client API config for the given namespace; namespace = "" for all namespaces.
// see https://github.com/kubernetes/client-go/tree/v0.23.5/examples
type Config func(namespace string) (*rest.Config, error)

func GetKubernetesClientset(configCallback Config, namespace string) (*kubernetes.Clientset, error) {
	if config, err := configCallback(namespace); err != nil {
		return nil, fmt.Errorf("[ERROR] error calling k8s.Config: %s", err)
	} else {
		if clientset, err := kubernetes.NewForConfig(config); err != nil {
			return nil, fmt.Errorf("[ERROR] error getting Kubernetes Clientset: %s", err)
		} else {
			return clientset, nil
		}
	}
}

// getSecretsClient returns a Kubernetes Secrets client for the given namespace; namespace = "" for all namespaces.
func GetSecretsClient(configCallback Config, namespace string) (v1.SecretInterface, error) {
	if clientset, err := GetKubernetesClientset(configCallback, namespace); err != nil {
		return nil, fmt.Errorf("[ERROR] error getting Kubernetes Clientset: %s", err)
	} else {
		return clientset.CoreV1().Secrets(namespace), nil
	}
}
