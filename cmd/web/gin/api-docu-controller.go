package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ContentTypeHTML = "text/html; charset=utf-8"
)

type IDocuController interface {
	Get(*gin.Context)
}

type DefaultDocuController struct{}

func NewDefaultSwaggerController() *DefaultDocuController {
	return &DefaultDocuController{}
}

func (ctrl *DefaultDocuController) Get(c *gin.Context) {
	c.HTML(http.StatusOK, DocuHTMLFileName, nil)
}
