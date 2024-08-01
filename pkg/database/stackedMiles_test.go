package database

import (
	"fmt"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"
	"testing"
)

func TestGetMilesData(t *testing.T) {

	fmt.Println("Test: TestGetMilesData")

	var mileData []models.MilesData
	var mileData2 []models.MilesData
	var err error

	mileData, err = GetMilesData("week_to_date", "transportation")
	if err != nil {
		t.Error("Failed to get miles data with error: ", err)
	}

	if mileData == nil {
		t.Error("Failed to get any data from the database")
	}
	if len(mileData) == 0 {
		t.Error("Failed to get any data from the database")
	}

	mileData2, err = GetMilesData("month_to_date", "transportation")
	if err != nil {
		t.Error("Failed to get miles data with error: ", err)
	}

	if mileData2 == nil {
		t.Error("Failed to get any data from the database")
	}
	if len(mileData2) == 0 {
		t.Error("Failed to get any data from the database")
	}

	for _, data := range mileData {

		fmt.Printf("Name: %s, TotalLoadedMiles: %f, TotalEmptyMiles: %f, TotalMiles: %f, PercentEmpty: %f\n", data.Name, data.TotalLoadedMiles, data.TotalEmptyMiles, data.TotalMiles, data.PercentEmpty)
	}

	for _, data := range mileData2 {

		fmt.Printf("Name: %s, TotalLoadedMiles: %f, TotalEmptyMiles: %f, TotalMiles: %f, PercentEmpty: %f\n", data.Name, data.TotalLoadedMiles, data.TotalEmptyMiles, data.TotalMiles, data.PercentEmpty)
	}
	t.Errorf("Test: TestGetMilesData - FAILED")

	fmt.Println("Test: TestGetMilesData - PASSED")
}
