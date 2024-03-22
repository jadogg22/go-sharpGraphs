package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler function for	root route for debuging
func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "This came from the test handler",
	})
}

// ---------- Transportation Handlers ----------
func Trans_year_by_year(c *gin.Context) {
	// get date from system
	// conncet to database
	// pull all year by year data and and compair data
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

func Trans_stacked_miles(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

// ---------- Logisitics Handlers ----------
func Log_year_by_year(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

func Log_stacked_miles(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

// ---------- Dispatch Handlers ----------

func Dispach_week_to_date(c *gin.Context) {
	data := []DriverData{
		{
			Driver:               "felicia mcmichael",
			Day:                  []string{"2/26/24", "2/26/24", "2/26/24", "2/26/24", "2/26/24"},
			TotalTrucks:          []int{13, 22, 22, 23, 24},
			TotalMiles:           []int{16585, 29004, 34796, 47331, 54425},
			LoadedMiles:          []int{16022, 27841, 32806, 44972, 50546},
			RevenueWithoutFuel:   []float64{37767.65, 73531.66, 87396.90, 117991.35, 132408.13},
			MilesPerTruck:        []int{1276, 1318, 1582, 1893, 2268},
			RevenuePerTruck:      []float64{2905.20, 3342.35, 3972.59, 4719.65, 5517.01},
			Deadhead:             []float64{3.40, 4.0, 5.7, 5.0, 7.10},
			RevenuePerLoadedMile: []float64{2.36, 2.64, 2.66, 2.62, 2.62},
			OnTimeService:        []float64{100.00, 94.59, 95.45, 96.49, 97.57},
		},
		{
			Driver:               "trina sepulveda",
			Day:                  []string{"2/26/24", "2/26/24", "2/26/24", "2/26/24", "2/26/24"},
			TotalTrucks:          []int{14, 20, 22, 24, 26},
			TotalMiles:           []int{8177, 13890, 23318, 33791, 46455},
			LoadedMiles:          []int{7661, 12584, 21002, 31197, 42534},
			RevenueWithoutFuel:   []float64{2165.60, 40452.71, 66177.25, 90165.51, 114853.02},
			MilesPerTruck:        []int{584, 695, 1014, 1408, 1936},
			RevenuePerTruck:      []float64{154.69, 2022.64, 2877.27, 3756.90, 4785.54},
			Deadhead:             []float64{6.30, 9.3, 9.9, 7.7, 6.40},
			RevenuePerLoadedMile: []float64{0.28, 3.21, 3.15, 2.89, 2.70},
			OnTimeService:        []float64{100.00, 93.33, 97.83, 98.27, 98.57},
		},
	}

	c.JSON(http.StatusOK, data)
}
