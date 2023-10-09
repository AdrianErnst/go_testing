/*
 * Go_Testing API
 *
 * Openapi specification for the cmd/web main
 *
 * API version: 1.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name		string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method		string
	// Pattern is the pattern of the URI.
	Pattern	 	string
	// HandlerFunc is the handler function of this route.
	HandlerFunc	gin.HandlerFunc
}

// NewRouter returns a new router.
func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
	router := gin.Default()
	for _, route := range getRoutes(handleFunctions) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Default handler for not yet implemented routes
func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

type ApiHandleFunctions struct {

	// Routes for the DocuAPI part of the API
	DocuAPI DocuAPI
	// Routes for the K8sAPI part of the API
	K8sAPI K8sAPI
}

func getRoutes(handleFunctions ApiHandleFunctions) []Route {
	return []Route{
	
		{
			"GetApi",
			http.MethodGet,
			"/api/",
			handleFunctions.DocuAPI.GetApi,
		},
		{
			"GetPodCountByNamespace",
			http.MethodGet,
			"/api/k8s/pods/count/:namespace",
			handleFunctions.K8sAPI.GetPodCountByNamespace,
		},
	}
}
