/*
 * Go_Testing API
 *
 * Openapi specification for the cmd/web main
 *
 * API version: 1.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"github.com/gin-gonic/gin"
)

type DocuAPI struct {
	// Get /api/
	// Get this html
	GetApi gin.HandlerFunc
}
