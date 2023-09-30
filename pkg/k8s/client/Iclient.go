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

// list implementations to quickly see which dont yet implement the interface fully
var _ IK8sClient = &inClusterK8sClient{}
var _ IK8sClient = &outClusterK8sClient{}

type IK8sClient interface {
	GetPodCount(ctx context.Context, namespace string) (int, error)
	LogPodNames(ctx context.Context, namespace string)
}

func NewK8sClient(logger *log.Logger, cType ClientType) IK8sClient {
	switch cType {
	case OutCluster:
		return newOutClusterClient(logger)
	case InCluster:
		fallthrough
	default:
		return newInClusterClient(logger)
	}
}
