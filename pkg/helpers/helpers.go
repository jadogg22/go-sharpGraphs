package helpers

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

func GetYearAndWeek() (int, int) {
	// get current date
	date := time.Now()

	// get year and week from system
	year, week := date.ISOWeek()

	return year, week
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
