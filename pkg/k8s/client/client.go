package k8s_client

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var client *kubernetes.Clientset

func init() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func GetPodCount(ctx context.Context, namespace string) (int, error) {
	pods, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	return len(pods.Items), err
}

func LogPodNames(ctx context.Context, namespace string) {
	list, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	if errors.IsNotFound(err) {
		fmt.Printf("Error while trying to list all pods from default service\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pods %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found following pods in default namespace\n")
		for _, item := range list.Items {
			fmt.Println(item.ObjectMeta.Namespace, "---", item.ObjectMeta.Name)
		}
	}
}
