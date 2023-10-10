package main

import (
	"go_testing/cmd/web/controllers"
	openapi_gin "go_testing/cmd/web/openapi-gin/go"
	"go_testing/internal/flags"
	k8s_client "go_testing/pkg/k8s/client"
	"go_testing/pkg/logger"
	"go_testing/util"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

const NamespaceKey = util.ContextKey("Namespace")
const Namespace = "client"
const (
	DocuHtmlPath = "/app/docs/generated/*.html"
)

var (
	g errgroup.Group
)

func main() {
	config := flags.ParseFlags()
	dLogger := logger.NewLogger(logger.DefaultLoggerType)
	var clientType k8s_client.ClientType
	if config.External {
		clientType = k8s_client.OutCluster
	} else {
		clientType = k8s_client.InCluster
	}
	client := k8s_client.NewK8sClient(dLogger, clientType, &config.K8sClient)

	handleFuncs := openapi_gin.ApiHandleFunctions{
		DocuAPI: controllers.NewDocuApiController(
			&controllers.DocuApiController{},
		),
		K8sAPI: controllers.NewK8sApiController(
			controllers.K8sController{
				Client: client,
			},
		),
	}

	router := openapi_gin.NewRouter(handleFuncs)
	router.LoadHTMLGlob(DocuHtmlPath)
	server := &http.Server{
		Addr:         ":9292",
		Handler:      router.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		dLogger.Fatal(err)
	}
}
