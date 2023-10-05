package gin

import (
	k8s_client "go_testing/pkg/k8s/client"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Param string

const (
	Namespace Param = "namespace"
)

type IK8sController interface {
	GetPodCount(c *gin.Context)
}

type DefaultK8sController struct {
	Client k8s_client.IK8sClient
}

func (ctrl DefaultK8sController) GetPodCount(c *gin.Context) {
	namespace := c.Param(string(Namespace))
	count, err := ctrl.Client.GetPodCount(c, namespace)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}
