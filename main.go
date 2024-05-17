package main

import (
	// import gin cors middleware

	"net/http"

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
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Message": "server is running"})
	})

	// // ---------- Transportation Handlers ----------
	r.GET("/Transportation/get_yearly_revenue", Trans_year_by_year)
	r.GET("/Transportation/Stacked_miles/:when", Trans_stacked_miles)
	r.GET("/Transportation/get_coded_revenue/:when", Trans_coded_revenue)

	r.POST("/Transportation/add/", Transportation_post)

	// // ---------- Logisitics Handlers ----------
	r.GET("/Logistics/get_yearly_revenue/", Log_year_by_year)

	// r.GET("/Logistics/Stacked_miles/", Log_stacked_miles)

	r.POST("/Logisics/add/", Logistics_post)

	// ---------- Dispatch Handlers ----------------
	r.GET("/Dispatch/Week_to_date/", Dispach_week_to_date)

	// ---------- receive data ----------
	// Define a POST endpoint to receive DispatcherStats data
	r.POST("/Dispatch/add/", Dispatch_post)
	// run the server on port 5000
	r.Run(":5000")

}

// CORS middleware handler
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers to allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			// Handle preflight requests
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue processing request
		next.ServeHTTP(w, r)
	})
}
