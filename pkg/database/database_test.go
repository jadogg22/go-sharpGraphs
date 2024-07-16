package database

import (
	"fmt"
	"testing"

	"github.com/jadogg22/go-sharpGraphs/pkg/models"
)

func TestPG_make_connection(t *testing.T) {
	fmt.Println("Test: testPG_make_connection")
	db, err := PG_Make_connection()
	if db == nil {
		t.Error("Failed to connect to database")
	}
	if err != nil {
		t.Error("Failed to connect to database with error: ", err)
	}

	// make a test query to check if the connection is working
	rows, err := db.Query("SELECT count(*) FROM trans_year_rev")
	if err != nil {
		t.Error("Failed to query database with error: ", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			t.Error("Failed to scan row with error: ", err)
		}
	}

	if count == 0 {
		t.Error("Failed to get any data from the database trans_year_rev")
	}

}

func TestPG_tables_exist(t *testing.T) {
	fmt.Println("Test: testPG_tables_exist")

	// add tables as we create more of them
	tables := []string{"daily_driver_data", "log_year_rev", "logistics", "trans_year_rev", "transportation"}
	db, _ := PG_Make_connection()

	for _, table := range tables {
		rows, err := db.Query("SELECT count(*) FROM " + table)
		if err != nil {
			t.Error("Failed to query database with error: ", err)
		}
		defer rows.Close()

		var count int
		for rows.Next() {
			err = rows.Scan(&count)
			if err != nil {
				t.Error("Failed to scan row with error: ", err)
			}
		}

		if count == 0 {
			t.Error("Failed to get any data from the database " + table)
		}
	}

	// check for table that doesnt exist
	_, err := db.Query("SELECT count(*) FROM no_table")
	if err == nil {
		t.Error("Found a table that doesnt exist")
	}
	fmt.Println("Test: testPG_tables_exist passed - all tables exist")
}

// func TestRefactoredYearbyYearData(t *testing.T) {
// fmt.Println("Test: testRefactoredYearbyYearData")
// db, _ := PG_Make_connection()
// transportationData, err := GetYearByYearData(db, "transportation")
// if err != nil {
// t.Error("Failed to get year by year data with error: ", err)
// }
//
// fmt.Println("This should not have ran")
//
// refactoredData, err := fetchRevenueDataToWeeklyRevenue(db, "trans_year_rev")
// if err != nil {
// t.Error("Failed to fetch revenue data with error: ", err)
// }
//
//	converTransportationData
// convertedData := make([]models.WeeklyRevenue, 52)
// for i, data := range transportationData {
// convertedData[i] = mapToWeeklyRevenue(data)
// }
//
// hndlPointer := func(rev *float64) float64 {
// if rev == nil {
// return 0
// }
// return *rev
// }
//
// for i, data := range refactoredData {
// assertEqual(t, convertedData[i].Name, data.Name, "Name")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2021), hndlPointer(data.Revenue2021), "Revenue2021 missmatch")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2022), hndlPointer(data.Revenue2022), "Revenue2022 missmatch")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2023), hndlPointer(data.Revenue2023), "Revenue2023 missmatch")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2024), hndlPointer(data.Revenue2024), "Revenue2024 missmatch")
// }
//
//	now test logistics data
// logisticsData, err := GetYearByYearData(db, "logistics")
// if err != nil {
// t.Error("Failed to get year by year data with error: ", err)
// }
//
// refactoredData, err = fetchRevenueDataToWeeklyRevenue(db, "log_year_rev")
// if err != nil {
// t.Error("Failed to fetch revenue data with error: ", err)
// }
//
//	convert logistics data to weekly revenue
// convertedData = make([]models.WeeklyRevenue, 52)
// for i, data := range logisticsData {
// convertedData[i] = mapToWeeklyRevenue(data)
// }
//
// for i, data := range refactoredData {
// assertEqual(t, convertedData[i].Name, data.Name, "Name")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2021), hndlPointer(data.Revenue2021), "Revenue2021 missmatch")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2022), hndlPointer(data.Revenue2022), "Revenue2022 missmatch")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2023), hndlPointer(data.Revenue2023), "Revenue2023 missmatch")
// assertEqual(t, hndlPointer(convertedData[i].Revenue2024), hndlPointer(data.Revenue2024), "Revenue2024 missmatch")
// }
//
// for i, data := range refactoredData {
// fmt.Println(printCompareWeeklyData(data, convertedData[i]))
// }
//
//	error to see data
// }
//
//assertEqual compares two values and fails the test if they're not equal
// func assertEqual(t *testing.T, got, want interface{}, msg string) {
// t.Helper()
// if got != want {
// t.Errorf("%s: got %v, want %v", msg, got, want)
// }
// }
//
//assertFloatEqual compares two float64 values with a small tolerance
// func assertFloatEqual(t *testing.T, got, want float64, msg string) {
// t.Helper()
// if math.Abs(got-want) > 0.0001 { // You can adjust the tolerance as needed
// t.Errorf("%s: got %v, want %v", msg, got, want)
// }
// }
//
// func printCompareWeeklyData(wd1, wd2 models.WeeklyRevenue) string {
// formatRevenue := func(rev *float64) string {
// if rev == nil {
// return "nil"
// }
// return fmt.Sprintf("%.2f", *rev)
// }
//
// msg := fmt.Sprintf(`
// Name: %d | %d
// 2021: %s | %s
// 2022: %s | %s
// 2023: %s | %s
// 2024: %s | %s`,
// wd1.Name, wd2.Name,
// formatRevenue(wd1.Revenue2021), formatRevenue(wd2.Revenue2021),
// formatRevenue(wd1.Revenue2022), formatRevenue(wd2.Revenue2022),
// formatRevenue(wd1.Revenue2023), formatRevenue(wd2.Revenue2023),
// formatRevenue(wd1.Revenue2024), formatRevenue(wd2.Revenue2024))
//
// return msg
// }

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

	fmt.Println("Missing data:" + fmt.Sprint(len(missingData)))

	for _, data := range missingData {
		fmt.Println(data)
	}
	// looking like its working

	data, err := GetYearByYearDataRefactored(db, transportationData, "logistics")
	if err != nil {
		t.Error("Failed to get year by year data with error: ", err)
	}

	fmt.Println("Data:" + fmt.Sprint(len(data)))

	for _, data := range data {
		fmt.Println(data)
	}

	t.Error("This test is not complete")

}
