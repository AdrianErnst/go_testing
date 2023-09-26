package k8s_client

import (
	"context"
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type inClusterK8sClient struct {
	client *kubernetes.Clientset
	logger *log.Logger
}

var once sync.Once
var isClient K8sClient

func newInClusterClient(logger *log.Logger) K8sClient {
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
			client: client,
			logger: logger,
		}
	})

	return isClient
}

func (icClient *inClusterK8sClient) GetPodCount(ctx context.Context, namespace string) (int, error) {
	pods, err := icClient.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	return len(pods.Items), err
}

func (icClient *inClusterK8sClient) LogPodNames(ctx context.Context, namespace string) {
	list, err := icClient.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if errors.IsNotFound(err) {
		icClient.logger.Printf("Error while trying to list all pods from default service\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		icClient.logger.Printf("Error getting pods %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		icClient.logger.Printf("Found following pods in default namespace\n")
		for _, item := range list.Items {
			icClient.logger.Println(item.ObjectMeta.Namespace, "---", item.ObjectMeta.Name)
		}
	}
}
