package controllers

import (
	openapi_gin "go_testing/cmd/web/openapi-gin/go"
	k8s_client "go_testing/pkg/k8s/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Param string

const (
	Namespace Param = "namespace"
)

type IK8sController interface {
	GetPodCount(*gin.Context)
}

type K8sController struct {
	Client k8s_client.IK8sClient
}

var _ IK8sController = &K8sController{}

func NewK8sApiController(ctrl IK8sController) openapi_gin.K8sAPI {
	return openapi_gin.K8sAPI{
		GetPodCountByNamespace: ctrl.GetPodCount,
	}
}

func (ctrl K8sController) GetPodCount(c *gin.Context) {
	namespace := c.Param(string(Namespace))
	count, err := ctrl.Client.GetPodCount(c, namespace)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, openapi_gin.Count{
		Count: int32(count),
	})
}
