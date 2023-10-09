package gin

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	DocuHTMLFileName = "index.html"
	DocuHtmlPath     = "/app/docs/generated/*.html"
)

var once sync.Once
var router *gin.Engine

type Router struct {
	Handler http.Handler
}

type Controllers struct {
	K8s  IK8sController
	Docu IDocuController
}

func NewGinRouter(ctrls Controllers, globalMiddleware ...gin.HandlerFunc) *Router {
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
		router.LoadHTMLGlob(DocuHtmlPath)

		api := router.Group("/api")
		{
			k8s := api.Group("/k8s")
			{
				pods := k8s.Group("/pods")
				pods.GET(fmt.Sprintf("/count/:%s", string(Namespace)),
					ctrls.K8s.GetPodCount)
			}

			api.GET("/", ctrls.Docu.Get)
		}
	})

	return &Router{
		Handler: router.Handler(),
	}
}
