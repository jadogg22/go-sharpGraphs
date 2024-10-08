package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jadogg22/go-sharpGraphs/pkg/handlers"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/Transportation/get_yearly_revenue", handlers.Trans_year_by_year)
	r.GET("/api/Transportation/Stacked_miles/:when", handlers.Trans_stacked_miles)
	r.GET("/api/Transportation/get_coded_revenue/:when", handlers.Trans_coded_revenue)
	r.GET("/api/Transportation/Daily_Ops", handlers.Daily_Ops)

	r.GET("/api/Logistics/get_yearly_revenue", handlers.Log_year_by_year)

	r.GET("/api/Logistics/MTD", handlers.LogisticsMTD)

	r.GET("/api/Dispatch/Week_to_date", handlers.Dispach_week_to_date)
	return r
}

func TestGetYearlyRevenue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := SetupRouter()

	req := httptest.NewRequest("GET", "/api/Transportation/get_yearly_revenue", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	fmt.Println("response for endpoint: /api/Transportation/get_yearly_revenue is: ", w.Body.String())
}

func TestEndpoints(t *testing.T) {
	// Set the gin mode to test - ? this doesn't really seem to be doing much tbh
	gin.SetMode(gin.TestMode)

	router := SetupRouter()

	endponts := []string{
		"/api/Transportation/get_yearly_revenue",
		"/api/Transportation/Stacked_miles/week",
		"/api/Transportation/get_coded_revenue/week",
		"/api/Transportation/Daily_Ops",
		"/api/Logistics/get_yearly_revenue",
		"/api/Logistics/MTD",
		"/api/Dispatch/Week_to_date",
	}

	for _, endpoint := range endponts {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest("GET", endpoint, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			fmt.Println("response for endpoint: ", endpoint, " is: ", tail(w.Body.String(), 100))
		})
	}
}

func tail(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if n >= len(s) {
		return s
	}
	return s[len(s)-n:]
}
