package main

import (
	"go_testing/cmd/web/gin"
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
	router := gin.NewGinRouter(
		gin.Controllers{
			K8s: gin.DefaultK8sController{
				Client: client,
			},
			Docu: gin.NewDefaultSwaggerController(),
		},
	)
	server := &http.Server{
		Addr:         ":9292",
		Handler:      router.Handler,
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
