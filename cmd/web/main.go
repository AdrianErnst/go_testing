package main

import (
	"context"
	"errors"
	"flag"
	k8s_client "go_testing/pkg/k8s/client"
	"go_testing/pkg/logger"
	"go_testing/util"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/util/homedir"
)

type Config struct {
	external  bool
	k8sClient k8s_client.K8sClientConfig
}

const NamespaceKey = util.ContextKey("Namespace")
const Namespace = "client"

func main() {
	config := parseFlags()

	// Todo use context with cancel when switching to web server
	ctx := context.Background()
	ctx = context.WithValue(ctx, NamespaceKey, Namespace)
	defaultLogger := logger.NewLogger(logger.DefaultLoggerType)
	var clientType k8s_client.ClientType
	if config.external {
		clientType = k8s_client.OutCluster
	} else {
		clientType = k8s_client.InCluster
	}
	client := k8s_client.NewK8sClient(defaultLogger, clientType, &config.k8sClient)

	for {
		podCount, err := client.GetPodCount(ctx, Namespace)
		if err != nil {
			panic(err.Error())
		}
		defaultLogger.Printf("There are %d pods in the cluster\n", podCount)

		client.LogPodNames(ctx, Namespace)

		time.Sleep(60 * time.Second)
	}
}

func parseFlags() Config {
	// read flags
	var external bool
	flag.BoolVar(&external, "external", false, "boot for use outside of a k8s cluster")
	var kubeconfigFile string
	flag.StringVar(&kubeconfigFile, "kubeconfigFile", "", "path to the kubeconfig, if not set $HOMEDIR/.kube/config will be used")

	flag.Parse()

	if external {
		if home := homedir.HomeDir(); home != "" {
			kubeconfigFile = filepath.Join(home, ".kube", "config")
		}
		_, err := os.Stat(kubeconfigFile)
		if kubeconfigFile == "" || errors.Is(err, os.ErrNotExist) {
			panic("flag kubeconfigFile must be set to a valid path when setting external")
		}
	}

	return Config{
		external: external,
		k8sClient: k8s_client.K8sClientConfig{
			OutConfig: &k8s_client.OutClusterClientConfig{
				KubeconfigFile: kubeconfigFile,
			},
		},
	}
}
