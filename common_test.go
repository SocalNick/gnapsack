package main

import (
	"github.com/gin-gonic/gin"
)

// Helper function to create a router during testing
func getRouter(withTemplates bool) *gin.Engine {
	r := gin.Default()
	return r
}
