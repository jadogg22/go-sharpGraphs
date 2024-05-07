package main

import (
	// import gin cors middleware

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main function to display different endpoints for ease of use.
//	Handlers.go is where the logic is stored for each endpoint

func main() {
	r := gin.Default()
	// setup cors middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin, Cobtent-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Message": "server is running"})
	})

	// // ---------- Transportation Handlers ----------
	r.GET("/Transportation/Year_by_year/", Trans_year_by_year)
	//r.GET("/Transportation/Stacked_miles/", Trans_stacked_miles)
	r.GET("/Transportation/Coded_revenue/:when", Trans_coded_revenue)

	// // ---------- Logisitics Handlers ----------
	// r.GET("/Logistics/Year_by_year/", Log_year_by_year)
	// r.GET("/Logistics/Stacked_miles/", Log_stacked_miles)

	// ---------- Dispatch Handlers ----------------
	r.GET("/Dispatch/Week_to_date/", Dispach_week_to_date)

	// ---------- receive data ----------
	// Define a POST endpoint to receive DispatcherStats data
	r.POST("/Dispatch/add/", Dispatch_post)
	// run the server on port 5000
	r.Run(":5000")

}
