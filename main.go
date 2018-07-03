package main

import (
	"fmt"

	k8sMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sCoreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	k8sClientCmd "k8s.io/client-go/tools/clientcmd"
	"k8s.io/api/core/v1"
)

func main() {
	kubeConfig := k8sClientCmd.NewNonInteractiveDeferredLoadingClientConfig(
		k8sClientCmd.NewDefaultClientConfigLoadingRules(),
		&k8sClientCmd.ConfigOverrides{},
	)

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		panic(err)
	}

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	coreClient, err := k8sCoreV1.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	pods, err := coreClient.Pods(namespace).List(k8sMetaV1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("  %s\n", pod.Name)
	}

	watch, err := coreClient.Services(namespace).Watch(k8sMetaV1.ListOptions{})

	for {
		select {
		case event := <-watch.ResultChan():
			service := event.Object.(*v1.Service)

			fmt.Printf("Service %s has the following labels:", service.Name)

			for key, value := range service.Labels {
				fmt.Printf("Key, Value: %s %s\n", key, value)
			}
		}
	}
}
