package k8s_client

import (
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type OutClusterClientConfig struct {
	KubeconfigFile string
}

type outClusterK8sClient struct {
	*baseK8sClient
}

var outClient IK8sClient

func newOutClusterClient(logger *log.Logger, config OutClusterClientConfig) IK8sClient {
	once.Do(func() {
		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", config.KubeconfigFile)
		if err != nil {
			panic(err.Error())
		}

		// creates the clientset
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
		outClient = &outClusterK8sClient{
			baseK8sClient: &baseK8sClient{client, logger},
		}
	})

	return outClient
}

func DefaultKubeConfigPath() string {
	if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	}

	return ""
}