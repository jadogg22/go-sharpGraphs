package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jadogg22/go-sharpGraphs/pkg/handlers"
	"time"
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
	r.GET("/api2/Sportsman", handlers.Sportsman)
	// Select the spacific sportsman by dates for invoice
	r.GET("/api2/Sportsman/:date1/:date2", handlers.SportsmanWithDates)

	apiGroup := r.Group("/api")
	apiGroup.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.0.62"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	{
		// ---------- Transportation Handlers ----------
		TransportationGroup := apiGroup.Group("/Transportation")
		{
			TransportationGroup.GET("/dashboard", handlers.Dashboard)
			TransportationGroup.GET("/get_yearly_revenue", handlers.Trans_year_by_year)
			TransportationGroup.GET("/Stacked_miles/:when", handlers.Trans_stacked_miles)
			TransportationGroup.GET("/get_coded_revenue/:when", handlers.Trans_coded_revenue)
			TransportationGroup.GET("/Daily_Ops", handlers.Daily_Ops)
			TransportationGroup.Any("/Generate_Sportsmans", handlers.Generate_Sportsmans)
			TransportationGroup.GET("/DriverManager", handlers.DriverManager)
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
