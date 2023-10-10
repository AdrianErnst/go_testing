package controllers

import (
	openapi_gin "go_testing/cmd/web/openapi-gin/go"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	DocuHTMLFileName = "index.html"
)

type IDocuApiController interface {
	GetApi(c *gin.Context)
}

type DocuApiController struct{}

var _ IDocuApiController = &DocuApiController{}

func NewDocuApiController(ctrl IDocuApiController) openapi_gin.DocuAPI {
	return openapi_gin.DocuAPI{
		GetApi: ctrl.GetApi,
	}
}

func (ctrl *DocuApiController) GetApi(c *gin.Context) {
	c.HTML(http.StatusOK, DocuHTMLFileName, nil)
}
