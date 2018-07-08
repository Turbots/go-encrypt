package main

import (
	"fmt"

	k8sMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	k8sClientCmd "k8s.io/client-go/tools/clientcmd"
	routev1 "github.com/openshift/api/route/v1"
	routeClientv1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
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

	coreClient, err := k8sCorev1.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	pods, err := coreClient.Pods(namespace).List(k8sMetav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Pods in namespace %s:\n", namespace)
	for _, pod := range pods.Items {
		fmt.Printf("  %s\n", pod.Name)
	}

	routeClient, err := routeClientv1.NewForConfig(restConfig)
	routeWatch, err := routeClient.Routes(namespace).Watch(k8sMetav1.ListOptions{})

	for {
		select {
		case event := <-routeWatch.ResultChan():
			route := event.Object.(*routev1.Route)

			fmt.Println("Route\t", route.Name, "\t", route.Spec.Host)
		}
	}
}
