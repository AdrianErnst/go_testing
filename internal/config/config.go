package config

import k8s_client "go_testing/pkg/k8s/client"

type Config struct {
	External  bool
	K8sClient k8s_client.K8sClientConfig
}
