package main

import (
	"context"
	k8s_client "go_testing/pkg/k8s/client"
	"go_testing/pkg/logger"
	"go_testing/util"
	"time"
)

const NamespaceKey = util.ContextKey("Namespace")
const Namespace = "client"

func main() {
	// Todo use context with cancel when switching to web server
	ctx := context.Background()
	ctx = context.WithValue(ctx, NamespaceKey, Namespace)
	defaultLogger := logger.NewLogger(logger.DefaultLoggerType)
	client := k8s_client.NewK8sClient(defaultLogger, k8s_client.InCluster)

	for {
		podCount, err := client.GetPodCount(ctx, Namespace)
		if err != nil {
			panic(err.Error())
		}
		defaultLogger.Printf("There are %d pods in the cluster\n", podCount)

		client.LogPodNames(ctx, Namespace)

		time.Sleep(30 * time.Second)
	}
}
