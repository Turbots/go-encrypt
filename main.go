package main

import (
	"log"

	k8sMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClientCmd "k8s.io/client-go/tools/clientcmd"
	routev1 "github.com/openshift/api/route/v1"
	routeClientv1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"time"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	kubeConfig := k8sClientCmd.NewNonInteractiveDeferredLoadingClientConfig(
		k8sClientCmd.NewDefaultClientConfigLoadingRules(),
		&k8sClientCmd.ConfigOverrides{},
	)

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		panic(err)
	}

	log.Println("Go Encrypt - Copyright Ordina JWorks 2018")

	go watchRoutes(namespace, kubeConfig)

	for true {
		time.Sleep(time.Hour)
	}
}

func watchRoutes(namespace string, kubeConfig k8sClientCmd.ClientConfig) {
	log.Println("Watching for Routes in namespace ", namespace)

	restConfig, err := kubeConfig.ClientConfig()

	if err != nil {
		panic(err)
	}

	routeClient, err := routeClientv1.NewForConfig(restConfig)

	if err != nil {
		panic(err)
	}

	routeWatch, err := routeClient.Routes(namespace).Watch(k8sMetav1.ListOptions{})

	if err != nil {
		panic(err)
	}

	for {
		select {
		case event := <-routeWatch.ResultChan():
			route := event.Object.(*routev1.Route)

			log.Println("Route Event Received")
			log.Println(route.Name, "\t", route.Spec.Host)
		}
	}
}
