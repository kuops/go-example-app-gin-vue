package kubernetes

import (
	"github.com/kuops/go-example-app/server/pkg/config"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"

)

func NewClientSet(cfg *config.KubernetesConfig) (*kubernetes.Clientset, error) {
	var kubeconfig *rest.Config
	var err error
	if cfg.Type == "in" {
		kubeconfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil,err
	}
	return clientset,nil
}

func NewDynamic(cfg *config.KubernetesConfig)  (dynamic.Interface,error)  {
	var kubeconfig *rest.Config
	var err error
	if cfg.Type == "in" {
		kubeconfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	client, err := dynamic.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	return client,nil
}
