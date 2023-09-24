package main

import (
	"context"
	"fmt"
	k8s_client "go_testing/pkg/k8s/client"
	"go_testing/util"
	"time"
)

const NamespaceKey = util.ContextKey("Namespace")
const Namespace = "client"

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, NamespaceKey, Namespace)
	// Todo use context with cancel when switching to web server

	for {
		podCount, err := k8s_client.GetPodCount(ctx, Namespace)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", podCount)

		k8s_client.LogPodNames(ctx, Namespace)

		time.Sleep(30 * time.Second)
	}
}
