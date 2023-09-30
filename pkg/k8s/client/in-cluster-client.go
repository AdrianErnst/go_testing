package k8s_client

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type inClusterK8sClient struct {
	*baseK8sClient
}

var isClient IK8sClient

func newInClusterClient(logger *log.Logger) IK8sClient {
	once.Do(func() {
		// creates the in-cluster config
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		// creates the clientset
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		isClient = &inClusterK8sClient{
			baseK8sClient: &baseK8sClient{client, logger},
		}
	})

	return isClient
}
