package gin

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var once sync.Once
var router *gin.Engine

type Router struct {
	Handler http.Handler
}

func NewGinRouter(ctrl IK8sController, globalMiddleware ...gin.HandlerFunc) *Router {
	once.Do(func() {
		router = gin.New()
		if len(globalMiddleware) == 0 {
			router.Use(gin.Logger())
			router.Use(gin.Recovery())
		} else {
			for _, middleware := range globalMiddleware {
				router.Use(middleware)
			}
		}

		api := router.Group("/api")
		{
			k8s := api.Group("/k8s")
			{
				pods := k8s.Group("/pods")
				pods.GET(fmt.Sprintf("/count/:%s", string(Namespace)),
					ctrl.GetPodCount)
			}
		}
	})

	return &Router{
		Handler: router.Handler(),
	}
}
