package main

import (
	k8sMetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sClientCmd "k8s.io/client-go/tools/clientcmd"
	routev1 "github.com/openshift/api/route/v1"
	routeClientv1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	"github.com/op/go-logging"
	"os"
	"k8s.io/apimachinery/pkg/watch"
)

var log = logging.MustGetLogger("go-encrypt")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc:15s} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	standardOut := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(standardOut, format)
	logging.SetBackend(formatter)

	kubeConfig := k8sClientCmd.NewNonInteractiveDeferredLoadingClientConfig(
		k8sClientCmd.NewDefaultClientConfigLoadingRules(),
		&k8sClientCmd.ConfigOverrides{},
	)

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		panic(err)
	}

	log.Info("Go Encrypt - Copyright Ordina JWorks 2018")

	go watchRoutes(namespace, kubeConfig)

	select {}
}

func watchRoutes(namespace string, kubeConfig k8sClientCmd.ClientConfig) {
	log.Info("Watching for Routes in namespace %v", namespace)

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

			log.Info("%v %v %v %v", route.Name, route.Spec.Host, "TLS:", route.Spec.TLS)

			if event.Type == watch.Added || event.Type == watch.Modified {
				if route.Spec.TLS == nil {
					go secureRoute(route)
				}
			}
		}
	}
}

func secureRoute(route *routev1.Route) {
	log.Info("Creating route for %v", route.Name)
}
