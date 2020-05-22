package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gusto/gnapsack/build"
	"github.com/gusto/gnapsack/db"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var gormDb *gorm.DB

func buildDistributionsSubset(c *gin.Context) {
	var inputDistributionsSubset build.InputDistributionsSubset
	if err := c.ShouldBindJSON(&inputDistributionsSubset); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	//	_, err := db.FindBuildForDistributionSubset(gormDb, inputDistributionsSubset)
	//	if err != nil {
	//		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	//		return
	//	}

	buildDistribution, err := inputDistributionsSubset.Distribution()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, buildDistribution)
}

func createBuildSubset(c *gin.Context) {
	var buildSubset build.Subset

	if err := c.ShouldBindJSON(&buildSubset); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := db.CreateBuildSubset(gormDb, buildSubset); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var empty struct{}
	c.JSON(http.StatusCreated, empty)
}

func listBuildSubsets(c *gin.Context) {
	buildSubsets, err := db.ListBuildSubsets(gormDb)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, buildSubsets)
}

func main() {
	gormDb = initializeDb()
	defer gormDb.Close()

	router := gin.Default()
	router.Use(formatRequest())
	router.Use(requestId())

	initializeRoutes(router)

	_ = router.Run(":8080")
}

func initializeDb() *gorm.DB {
	db, err := gorm.Open("postgres", "dbname=gnapsack_development sslmode=disable")
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	// Migrate the schema
	db.AutoMigrate(&build.Subset{}, &build.DistributionTestFile{})

	return db
}

func initializeRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.POST("/build_distributions/subset", buildDistributionsSubset)
		v1.POST("/build_subsets", createBuildSubset)
		v1.GET("/build_subsets", listBuildSubsets)
	}
}

// formatRequest generates ascii representation of a request
func formatRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		// Create return string
		var request []string
		// Add the request string
		url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
		request = append(request, url)
		// Add the host
		request = append(request, fmt.Sprintf("Host: %v", r.Host))
		// Loop through headers
		request = append(request, "\nHEADERS:")
		for name, headers := range r.Header {
			name = strings.ToLower(name)
			for _, h := range headers {
				request = append(request, fmt.Sprintf("%v: %v", name, h))
			}
		}

		// If this is a POST, add post data
		if r.Method == "POST" {
			request = append(request, "\nPOST DATA")
			_ = r.ParseForm()
			request = append(request, "\n")
			request = append(request, r.PostForm.Encode())
		}
		// Return the request as a string
		log.Print(strings.Join(request, "\n"))
		c.Next()
	}
}

func requestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		requestID := c.Request.Header.Get("X-Request-Id")

		// Create request id with UUID4
		if requestID == "" {
			uuid4, _ := uuid.NewRandom()
			requestID = uuid4.String()
		}

		// Expose it for use in the application
		c.Set("RequestId", requestID)

		// Set X-Request-Id header
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}
