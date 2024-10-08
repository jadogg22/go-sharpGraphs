package getdata

import (
	"testing"
)

func TestEnpoints(t *testing.T) {
	// Set the gin mode to test - ? this doesn't really seem to be doing much tbh
	// gin.SetMode(gin.TestMode)

	// router := SetupRouter()

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
			// req := httptest.NewRequest("GET", endpoint, nil)
		})
	}
}
