package main

import (
	"context"
	"go_testing/internal/flags"
	k8s_client "go_testing/pkg/k8s/client"
	"go_testing/pkg/logger"
	"go_testing/util"
	"time"
)

const NamespaceKey = util.ContextKey("Namespace")
const Namespace = "client"

func main() {
	config := flags.ParseFlags()

	// Todo use context with cancel when switching to web server
	ctx := context.Background()
	ctx = context.WithValue(ctx, NamespaceKey, Namespace)
	dLogger := logger.NewLogger(logger.DefaultLoggerType)
	var clientType k8s_client.ClientType
	if config.External {
		clientType = k8s_client.OutCluster
	} else {
		clientType = k8s_client.InCluster
	}
	client := k8s_client.NewK8sClient(dLogger, clientType, &config.K8sClient)

	for {
		podCount, err := client.GetPodCount(ctx, Namespace)
		if err != nil {
			panic(err.Error())
		}
		dLogger.Printf("There are %d pods in the cluster\n", podCount)

		client.LogPodNames(ctx, Namespace)

		time.Sleep(30 * time.Second)
	}
}
