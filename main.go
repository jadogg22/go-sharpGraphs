package main

import (
	// my local packages
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jadogg22/go-sharpGraphs/pkg/handlers"
	"net/http"
	"strings"
)

// main function to display different endpoints for ease of use.
//	Handlers.go is where the logic is stored for each endpoint
//

func main() {
	r := gin.Default()

	// setup cors middleware
	r.Use(CORSMiddleware())
	r.Use(staticFileMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Serve static files from the embedded frontend/dist directory
	r.Use(static.Serve("/", static.LocalFile("./frontend/dist", true)))

	// // ---------- Transportation Handlers ----------
	r.GET("/api/Transportation/get_yearly_revenue/", handlers.Trans_year_by_year)
	r.GET("/api/Transportation/Stacked_miles/:when", handlers.Trans_stacked_miles)
	r.GET("/api/Transportation/get_coded_revenue/:when", handlers.Trans_coded_revenue)
	r.GET("/api/Transportation/Daily_Ops/", handlers.Daily_Ops)

	r.POST("/api/Transportation/add/", handlers.Transportation_post)

	// // ---------- Logisitics Handlers ----------
	r.GET("/api/Logistics/get_yearly_revenue", handlers.Log_year_by_year)

	// r.GET("/Logistics/Stacked_miles/", Log_stacked_miles)

	r.POST("/api/Logistics/add/", handlers.Logistics_post)

	// ---------- Dispatch Handlers ----------------
	r.GET("/api/Dispatch/Week_to_date/", handlers.Dispach_week_to_date)

	// ---------- receive data ----------
	// Define a POST endpoint to receive DispatcherStats data
	r.POST("/api/Dispatcher/add/", handlers.Dispatch_post)
	// run the server on port 5000
	r.Run(":5000")

}

// CORS middleware handler
// CORSMiddleware is a custom CORS middleware that takes a pointer to gin.Context
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set headers to allow all origins
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "True")

		if c.Request.Method == "OPTIONS" {
			// Handle preflight requests
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Continue processing request
		c.Next()
	}
}

func staticFileMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, ".js") {
			c.Writer.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(c.Request.URL.Path, ".css") {
			c.Writer.Header().Set("Content-Type", "text/css")
		}
		c.Next()
	}
}
