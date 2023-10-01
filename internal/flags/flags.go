package flags

import (
	"errors"
	"flag"
	"go_testing/internal/config"
	k8s_client "go_testing/pkg/k8s/client"
	"os"
)

const (
	External       = k8s_client.External
	KubeconfigFile = k8s_client.KubeconfigFile
)

func ParseFlags() config.Config {
	// read flags
	var external bool
	flag.BoolVar(&external, External, false, "boot for use outside of a k8s cluster")
	var kubeconfigFile string
	flag.StringVar(&kubeconfigFile, KubeconfigFile, "", "path to the kubeconfig, if not set $HOMEDIR/.kube/config will be used")

	flag.Parse()

	if external {
		kubeconfigFile = k8s_client.DefaultKubeConfigPath()
		_, err := os.Stat(kubeconfigFile)
		if kubeconfigFile == "" || errors.Is(err, os.ErrNotExist) {
			panic("flag kubeconfigFile must be set to a valid path when setting external")
		}
	}

	return config.Config{
		External: external,
		K8sClient: k8s_client.K8sClientConfig{
			OutConfig: &k8s_client.OutClusterClientConfig{
				KubeconfigFile: kubeconfigFile,
			},
		},
	}
}
