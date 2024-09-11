package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"
)

func TestGetWeeklyRevenueData(t *testing.T) {
	fmt.Println("Test: TestGetWeeklyRevenueData")
	db, _ := PG_Make_connection()
	data, err := FetchRevenueDataToWeeklyRevenue(db, "log_year_rev")
	if err != nil {
		t.Error("Failed to get year by year data with error: ", err)
	}

	if len(data) < 2 {
		t.Error("Failed to get correct data")
	}

	for d, r := range data {
		fmt.Println(r.Name, r.Revenue2021, r.Revenue2022, r.Revenue2023, d)
	}
}

func TestNewFetchMyCache(t *testing.T) {
	fmt.Println("Test: TestNewFetchMyCache")

	data, newestWeek, err := NewFetchMyCache()
	if err != nil {
		t.Error("Failed to get cache with error: ", err)
	}

	startDate := helpers.FirstDayOfWeek(newestWeek)
	today := time.Now()

	for d, r := range data {
		fmt.Println(r.Name, r.Revenue2021, r.Revenue2022, r.Revenue2023, d)
	}

	fmt.Println("Newest week: ", newestWeek)
	fmt.Println("Start date: ", startDate)
	fmt.Println("Today: ", today)

}
func mapToWeeklyRevenue(data map[string]interface{}) models.WeeklyRevenue {
	wr := models.WeeklyRevenue{
		Name: data["Name"].(int),
	}

	years := []string{"2020 Revenue", "2021 Revenue", "2022 Revenue", "2023 Revenue"}
	for _, year := range years {
		if rev, ok := data[year]; ok {
			revFloat := rev.(float64)
			switch year {
			case "2021 Revenue":
				wr.Revenue2021 = &revFloat
			case "2022 Revenue":
				wr.Revenue2022 = &revFloat
			case "2023 Revenue":
				wr.Revenue2023 = &revFloat
			}
		}
	}

	return wr
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func TestFindMissingData(t *testing.T) {
	fmt.Println("Test: testfindMissingData")
	db, _ := PG_Make_connection()
	transportationData, err := FetchRevenueDataToWeeklyRevenue(db, "log_year_rev")
	if err != nil {
		t.Error("Failed to get year by year data with error: ", err)
	}

	// get the missing data

	missingData, err := FindMissingData(transportationData)
	if err != nil {
		t.Error("Failed to get missing data with error: ", err)
	}

	// check if the missing data is correct
	if len(missingData) < 2 {
		t.Error("Failed to get correct missing data")
	}

	// looking like its working

	data, err := GetYearByYearDataRefactored(transportationData, "logistics")
	if err != nil {
		t.Error("Failed to get year by year data with error: ", err)
	}

	if len(data) != 52 {
		t.Error("Failed to get correct data")
	}
}
