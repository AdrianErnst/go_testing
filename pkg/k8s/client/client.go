package k8s_client

import (
	"context"
	"log"
)

type ClientType int64

const (
	InCluster ClientType = iota
)

// list implementations to quickly see which dont yet implement the interface fully
var _ K8sClient = &inClusterK8sClient{}

type K8sClient interface {
	GetPodCount(ctx context.Context, namespace string) (int, error)
	LogPodNames(ctx context.Context, namespace string)
}

func NewK8sClient(logger *log.Logger, cType ClientType) K8sClient {
	switch cType {
	case InCluster:
		fallthrough
	default:
		return newInClusterClient(logger)
	}
}
