package main

import (
	// my local packages
	"github.com/jadogg22/go-sharpGraphs/pkg/handlers"

	"github.com/gin-gonic/gin"
	"net/http"
)

// main function to display different endpoints for ease of use.
//	Handlers.go is where the logic is stored for each endpoint
//

func main() {
	r := gin.Default()

	// setup cors middleware
	r.Use(CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Message": "server is running"})
	})

	// // ---------- Transportation Handlers ----------
	r.GET("/Transportation/get_yearly_revenue/", handlers.Trans_year_by_year)
	r.GET("/Transportation/Stacked_miles/:when", handlers.Trans_stacked_miles)
	r.GET("/Transportation/get_coded_revenue/:when", handlers.Trans_coded_revenue)

	r.POST("/Transportation/add/", handlers.Transportation_post)

	// // ---------- Logisitics Handlers ----------
	r.GET("/Logistics/get_yearly_revenue", handlers.Log_year_by_year)

	// r.GET("/Logistics/Stacked_miles/", Log_stacked_miles)

	r.POST("/Logistics/add/", handlers.Logistics_post)

	// ---------- Dispatch Handlers ----------------
	r.GET("/Dispatch/Week_to_date/", handlers.Dispach_week_to_date)

	// ---------- receive data ----------
	// Define a POST endpoint to receive DispatcherStats data
	r.POST("/Dispatch/add/", handlers.Dispatch_post)
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
