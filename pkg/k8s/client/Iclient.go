package k8s_client

import (
	"context"
	"log"
	"sync"
)

var once sync.Once

type ClientType int64

const (
	Start ClientType = iota
	InCluster
	OutCluster
	End // should always be last element in enum
)

// flags for testing
const (
	External       = "external"
	KubeconfigFile = "kubeconfigFile"
)

// list implementations to quickly see which dont yet implement the interface fully
var _ IK8sClient = &inClusterK8sClient{}
var _ IK8sClient = &outClusterK8sClient{}

type IK8sClient interface {
	GetPodCount(ctx context.Context, namespace string) (int, error)
	LogPodNames(ctx context.Context, namespace string)
}

type K8sClientConfig struct {
	OutConfig *OutClusterClientConfig
}

func NewK8sClient(logger *log.Logger, cType ClientType, config *K8sClientConfig) IK8sClient {
	switch cType {
	case OutCluster:
		if config.OutConfig == nil {
			panic("trying to start external cluster without setting the needed config")
		}
		return newOutClusterClient(logger, *config.OutConfig)
	case InCluster:
		fallthrough
	default:
		return newInClusterClient(logger)
	}
}
