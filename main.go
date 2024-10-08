package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jadogg22/go-sharpGraphs/pkg/handlers"
	"net/http"
)

// main function to display different endpoints for ease of use.
//	Handlers.go is where the logic is stored for each endpoint
//

func main() {
	r := gin.Default()

	// setup cors middleware

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api2/vacation/:type", handlers.Vacation)

	apiGroup := r.Group("/api")
	apiGroup.Use(CORSMiddleware())
	{
		// ---------- Transportation Handlers ----------
		TransportationGroup := apiGroup.Group("/Transportation")
		{
			TransportationGroup.GET("/get_yearly_revenue", handlers.Trans_year_by_year)
			TransportationGroup.GET("/Stacked_miles/:when", handlers.Trans_stacked_miles)
			TransportationGroup.GET("/get_coded_revenue/:when", handlers.Trans_coded_revenue)
			TransportationGroup.GET("/Daily_Ops", handlers.Daily_Ops)
		}
		// ---------- Logisitics Handlers ----------
		LogisticsGroup := apiGroup.Group("/Logistics")
		{
			LogisticsGroup.GET("/get_yearly_revenue", handlers.Log_year_by_year)
			LogisticsGroup.GET("/MTD", handlers.LogisticsMTD)
		}
		// ---------- Dispatch Handlers ----------------
		DispatchGroup := apiGroup.Group("/Dispatch")
		{
			DispatchGroup.GET("/Week_to_date", handlers.Dispach_week_to_date)
		}
	}

	r.Run(":5000")

}

// CORS middleware handler
// CORSMiddleware is a middleware handler that adds CORS headers to requests
// make sure that the request is coming from the correct origin of my proxy server.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow cors from my proxy
		allowedOrigins := "http://192.168.0.62"

		origin := c.Request.Header.Get("Origin")

		if origin != allowedOrigins {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Origin", allowedOrigins)
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
