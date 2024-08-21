package helpers

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/models"
)

func GetYearAndWeek() (int, int) {
	// get current date
	date := time.Now()

	// get year and week from system
	year, week := date.ISOWeek()

	return year, week
}

func StartOfISOWeek(year int, week int) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, time.UTC)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year || (isoYear == year && isoWeek < week) {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}

// takes a slice of YearByYear data and returns the newest week and year
func GetNewestWeek(data []map[string]interface{}) int {
	newestYear := 0
	newestWeek := 0

	// Iterate through the first map to find the newest year
	for key := range data[0] {
		//check for keys that dont start with years and skip
		_, err := strconv.Atoi(key[:1])
		if err != nil {
			continue
		}

		// Keys are like "2023 Revenue"
		year, err := strconv.Atoi(key[:4])
		if err != nil {
			fmt.Println("Error: ", err)
			return 0
		}
		if year > newestYear {
			newestYear = year
		}
	}

	// Iterate through the data to find where the data stops
	for _, yearData := range data {
		// Check if yearData has the year + " Revenue" key
		_, ok := yearData[strconv.Itoa(newestYear)+" Revenue"]
		if ok {
			// Extract string from key
			nameStr, ok := yearData["Name"].(string)
			if !ok {
				fmt.Println("Error: ", ok)
				return 0
			}

			newerWeek := WeektoInt(nameStr)
			if newerWeek != 0 {
				newestWeek = newerWeek
			}
		} else {
			// If the year does not have the key then we can break
			break
		}
	}
	return newestWeek
}

// takes in a week and year and returns the quarter in this fomat: 2021 Q02
func GetQuarter(week, year int) string {
	// convert year to string
	yearStr := strconv.Itoa(year)
	if week <= 13 {
		return yearStr + " Q01"
	}
	if week <= 26 {
		return yearStr + " Q02"
	}
	if week <= 39 {
		return yearStr + " Q03"
	}
	return yearStr + " Q04"

}

// takes in a week and year and returns the month in this fomat: 2021 M09
func GetMonth(week, year int) string {
	// convert year to string
	yearStr := strconv.Itoa(year)
	if week <= 4 {
		return yearStr + " M01"
	}
	if week <= 8 {
		return yearStr + " M02"
	}
	if week <= 13 {
		return yearStr + " M03"
	}
	if week <= 17 {
		return yearStr + " M04"
	}
	if week <= 22 {
		return yearStr + " M05"
	}
	if week <= 26 {
		return yearStr + " M06"
	}
	if week <= 31 {
		return yearStr + " M07"
	}
	if week <= 35 {
		return yearStr + " M08"
	}
	if week <= 39 {
		return yearStr + " M09"
	}
	if week <= 44 {
		return yearStr + " M10"
	}
	if week <= 48 {
		return yearStr + " M11"
	}
	return yearStr + " M12"
}

// Function to count working days in the current month and return the current day's position
func CountWorkingDays() (totalWorkingDays, currentDay int) {
	// Get the current date
	today := time.Now()

	// Get the first and last day of the current month
	firstDay := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	lastDay := firstDay.AddDate(0, 1, -1)

	// Initialize count for working days
	workingDays := 0

	// Iterate through each day of the month
	for d := firstDay; d.Before(lastDay.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		// Check if the day is a weekday (Monday to Friday)
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			// Increment the count of working days
			workingDays++
			// Check if the current day is today
			if d.Day() == today.Day() {
				currentDay = workingDays
			}
		}
	}

	return workingDays, currentDay
}

func WeektoStr(week int) string {
	weekStr := strconv.Itoa(week)
	if len(weekStr) == 1 {
		weekStr = "0" + weekStr
	}
	return "W" + weekStr
}

// Helper function to convert week string to int
func WeektoInt(weekStr string) int {
	weekInt, err := strconv.Atoi(weekStr[1:])
	if err != nil {
		fmt.Println("Error converting week string to int:", err)
		return 0
	}
	return weekInt
}

// "2024W01" -> 2024, 1
func GetYearAndWeekFromStr(week string) (int, int, error) {
	year, err := strconv.Atoi(week[:4])
	if err != nil {
		fmt.Println("Error getting year from string: ", err)
		return 0, 0, err
	}
	weekInt, err := strconv.Atoi(week[6:])
	if err != nil {
		fmt.Println("Error getting week from string: ", err)
		return 0, 0, err
	}
	return year, weekInt, nil
}

func GetNewWeekFromWeek(weekStr string, amount int) string {
	year, week, err := GetYearAndWeekFromStr(weekStr)
	if err != nil {
		fmt.Println("error parsing weekstring")
	}
	fmt.Println("Week: ", week)
	fmt.Println("amount: ", amount)

	week += amount

	if week >= 0 {
		week += 52
		year -= 1
	}

	if week > 52 {
		week -= 52
		year += 1
	}

	// convert back into a string
	newWeek := strconv.Itoa(year) + " " + WeektoStr(week)
	return newWeek
}

// func sortByCount(revenueArray []map[string]interface{}) []map[string]interface{} {
// 	// Define the sorting function
// 	sortingFunc := func(i, j int) bool {
// 		return revenueArray[i]["Count"].(int64) > revenueArray[j]["Count"].(int64)
// 	}

// 	// Sort the revenueArray based on the count
// 	sort.Slice(revenueArray, sortingFunc)

// 	// Return the sorted revenueArray
// 	return revenueArray
// }

func GetStartDayOfWeek() time.Time {
	// Get the current date
	today := time.Now()

	// get the day of the week
	day := today.Weekday()

	// get monday as the start of the week
	startOfWeek := today.AddDate(0, 0, -int(day))

	return startOfWeek
}

func IsHoliday(date time.Time) bool {

	Holidays := []time.Time{
		// New Year's Day
		time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
		// Memorial Day
		time.Date(2024, time.May, 31, 0, 0, 0, 0, time.UTC),
		// Independence Day
		time.Date(2024, time.July, 4, 0, 0, 0, 0, time.UTC),
		// Labor Day
		time.Date(2024, time.September, 6, 0, 0, 0, 0, time.UTC),
		// Thanksgiving Day
		time.Date(2024, time.November, 25, 0, 0, 0, 0, time.UTC),
		// Christmas Day
		time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC),
	}
	for _, holiday := range Holidays {
		if date.Year() == holiday.Year() && date.Month() == holiday.Month() && date.Day() == holiday.Day() {
			return true
		}
	}
	return false
}

func CountWorkingDaysInTimeframe(start, end time.Time) (int, error) {
	// make sure start is before end
	if start.After(end) {
		err := fmt.Errorf("start date is after end date")
		return 0, err
	}

	// Initialize count for working days
	workingDays := 0

	// Iterate through each day of the month
	for d := start; d.Before(end.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		// Check if the day is a weekday (Monday to Friday)
		if d.Weekday() >= time.Monday && d.Weekday() <= time.Friday {
			if !IsHoliday(d) {
				// Increment the count of working days
				workingDays++
			}
		}
	}

	return workingDays, nil
}

func SortArraybyRevenue(revenueArray []map[string]interface{}) []map[string]interface{} {
	// Define the sorting function
	sortingFunc := func(i, j int) bool {
		return revenueArray[i]["Revenue"].(float64) > revenueArray[j]["Revenue"].(float64)
	}

	sort.Slice(revenueArray, sortingFunc)

	return revenueArray

}

func FilterRevenueArray(revenueArray []map[string]interface{}, totalrevenue float64, totalCount int64) ([]map[string]interface{}, error) {

	if totalrevenue <= 1 {
		// divide by zero error
		err := fmt.Errorf("totalrevenue is less then one, divide by zero")
		return revenueArray, err
	}

	if totalCount <= 1 {
		// divide by zero error
		err := fmt.Errorf("totalCount is less then one, divide by zero")
		return revenueArray, err
	}

	//need to create a new revenueArray to change the length
	newRevenueArray := []map[string]interface{}{}

	// create the map that will include all the small revenue
	otherMap := map[string]interface{}{
		"Count":    int64(0),
		"Reveneue": 0.0,
		"Code":     "Unique Items",

		"UniqueItems": 0,
		"Codes":       make([]string, 0),
	}

	for idx, dict := range revenueArray {
		if idx < 10 {
			Revenue := dict["Revenue"].(float64)
			Count := dict["Count"].(int64)
			revPercentage := ((Revenue / totalrevenue) * 100)
			countPercentage := (float64(Count) / float64(totalCount) * 100)

			dict["RevenuePercentage"] = revPercentage
			dict["CountPercentage"] = countPercentage

			newRevenueArray = append(newRevenueArray, dict)

		}

		// Only want the top ten now we're going to combine the rest into an other catagory

		//for now its only generic data, I need to come up with a way to show all the data probably with another array
		//TODO - create a nested array with additional info
		if idx >= 10 {
			code := dict["Code"].(string)
			revenue := dict["Revenue"].(float64)
			count := dict["Count"].(int64)

			otherMap["Count"] = otherMap["Count"].(int64) + count
			otherMap["Reveneue"] = otherMap["Reveneue"].(float64) + revenue
			otherMap["UniqueItems"] = otherMap["UniqueItems"].(int) + 1
			otherMap["Codes"] = append(otherMap["Codes"].([]string), code)

		}

	}

	// Calculate percentatges
	count := otherMap["Count"].(int64)
	revenue := otherMap["Reveneue"].(float64)

	otherMap["RevenuePercentage"] = (revenue / totalrevenue) * 100
	otherMap["CountPercentage"] = (float64(count) / float64(totalrevenue) * 100)

	newRevenueArray = append(newRevenueArray, otherMap)

	return newRevenueArray, nil
}

// Because its easier to query twice and get the totals and then combine them into a different struct here instead of a giant query that would be hard to read
// harder to debug im going to do it this way for now.

func AgregateLogisticMTDStats(ordersData []models.LogisticsOrdersData, stopOrdersData []models.LogisticsStopOrdersData) []models.LogisticsMTDStats {

	ordersMap := make(map[string]models.LogisticsOrdersData)

	for _, order := range ordersData {
		if exist, ok := ordersMap[order.Dispacher]; ok {
			exist.Charges += order.Charges
			exist.Miles += order.Miles
			exist.Truck_hire += order.Truck_hire
			ordersMap[order.Dispacher] = exist
		} else {
			ordersMap[order.Dispacher] = order
		}
	}

	var (
		// totals for the summation of all the data

		TotalOrders      int
		TotalStops       int
		TotalOrderFaults int
		TotalStopFaults  int
		TotalMiles       float64
		TotalCharges     float64
		TotalTruckHire   float64

		// variables to hold the data for each dispacher
		thisDispacher   string
		thisTotalStops  int
		thisTotalOrders int
		thisOrderFaults int
		thisStopFaults  int
		thisTruck_hire  float64
		thisCharges     float64
		thisMiles       float64

		// slice to hold the data
		agregateData = make([]models.LogisticsMTDStats, 0)
	)

	// outer loop to go through the stop data
	for _, stop := range stopOrdersData {
		if stop.Dispacher.String == "NULL" {
			// just add the stop data to the totals
			TotalStops += stop.Total_stops
			TotalOrders += stop.Total_orders
			TotalOrderFaults += stop.Order_faults
			TotalStopFaults += stop.Stop_faults
			break
		}

		thisDispacher = stop.Dispacher.String
		thisTotalStops = stop.Total_stops
		thisTotalOrders = stop.Total_orders
		thisOrderFaults = stop.Order_faults
		thisStopFaults = stop.Stop_faults

		TotalOrders += thisTotalOrders
		TotalStops += thisTotalStops
		TotalOrderFaults += thisOrderFaults
		TotalStopFaults += thisStopFaults

		thisTruck_hire = 0.0
		thisCharges = 0.0
		thisMiles = 0.0

		// inner loop to go through the orders data
		if order, ok := ordersMap[stop.Dispacher.String]; ok {
			thisDispacher = order.Dispacher
			thisTruck_hire = order.Truck_hire
			thisCharges = order.Charges
			thisMiles = order.Miles

			// Add for the totals data
			TotalTruckHire += thisTruck_hire
			TotalCharges += thisCharges
			TotalMiles += thisMiles

			// create a new data struct
			newData := models.NewLogisticsMTDStats(thisDispacher, thisTruck_hire, thisCharges, thisMiles, thisTotalStops, thisTotalOrders, thisOrderFaults, thisStopFaults)

			agregateData = append(agregateData, *newData)
		}
	}

	// Add the totals to the agregate data

	totalsData := models.NewLogisticsMTDStats("Total", TotalTruckHire, TotalCharges, TotalMiles, TotalStops, TotalOrders, TotalOrderFaults, TotalStopFaults)

	agregateData = append(agregateData, *totalsData)
	return agregateData

}

func removeSpacesToLower(s string) string {
	// Remove all spaces
	noSpaces := strings.ReplaceAll(s, " ", "")
	// Convert to lowercase
	lowerCase := strings.ToLower(noSpaces)
	return lowerCase
}

func SameName(s1, s2 string) bool {
	// Remove spaces and convert to lowercase
	s1 = removeSpacesToLower(s1)
	s2 = removeSpacesToLower(s2)
	// Check if the strings are the same
	return s1 == s2
}
