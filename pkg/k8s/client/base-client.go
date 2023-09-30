package k8s_client

import (
	"context"
	"log"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type baseK8sClient struct {
	client *kubernetes.Clientset
	logger *log.Logger
}

func (k8sClient *baseK8sClient) GetPodCount(ctx context.Context, namespace string) (int, error) {
	pods, err := k8sClient.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	return len(pods.Items), err
}

func (k8sClient *baseK8sClient) LogPodNames(ctx context.Context, namespace string) {
	list, err := k8sClient.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if errors.IsNotFound(err) {
		k8sClient.logger.Printf("Error while trying to list all pods from default service\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		k8sClient.logger.Printf("Error getting pods %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		k8sClient.logger.Printf("Found following pods in default namespace\n")
		for _, item := range list.Items {
			k8sClient.logger.Println(item.ObjectMeta.Namespace, "---", item.ObjectMeta.Name)
		}
	}
}