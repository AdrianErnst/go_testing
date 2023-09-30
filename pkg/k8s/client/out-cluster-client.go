package k8s_client

import (
	"flag"
	"log"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type outClusterK8sClient struct {
	*baseK8sClient
}

var outClient IK8sClient

func newOutClusterClient(logger *log.Logger) IK8sClient {
	once.Do(func() {
		// Todo set kubeconfig-filepath from tiltfile via environment or cli
		// Use [[ ! -z "$KUBECONFIG" ]] && echo "$KUBECONFIG" || echo "$HOME/.kube/config"
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"),
				"(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
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
