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

type K8sAPI struct {
	// Get /api/k8s/pods/count/:namespace
	// Get Pod count for a specific namespace
	GetPodCountByNamespace gin.HandlerFunc
}