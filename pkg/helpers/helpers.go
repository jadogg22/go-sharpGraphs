package helpers

import (
	"fmt"
	"slices"
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

func FirstDayOfWeek(week int) time.Time {
	if week < 1 || week > 53 {
		fmt.Println("Error: Week is less than 1")
		return time.Time{}
	}
	// Get the current date
	today := time.Now()

	// Define the first day of the year
	start := time.Date(today.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)

	// Find the first Monday of the year
	for start.Weekday() != time.Sunday {
		start = start.AddDate(0, 0, 1)
	}

	// Calculate the first day of the specified week
	firstDay := start.AddDate(0, 0, (week-1)*7)

	return firstDay
}

func GetWeekStartAndEndDates() (time.Time, time.Time) {
	// Get the current date
	today := time.Now()

	// Get the first day of the week
	startOfWeek := GetStartDayOfWeek()

	return startOfWeek, today
}

func GetMonthStartAndEndDates() (time.Time, time.Time) {
	// Get the current date
	today := time.Now()

	// Get the first day of the month
	firstDay := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	// Get the last day of the month
	lastDay := firstDay.AddDate(0, 1, -1)

	return firstDay, lastDay
}

func GetQuarterStartAndEndDates() (time.Time, time.Time) {
	// Get the current date
	today := time.Now()

	// Get the first day of the quarter
	quarter := monthToQuarter(int(today.Month()))
	quarterStartMonth := 3 * (quarter - 1) // 0, 3, 6, 9
	quarterStart := time.Date(today.Year(), time.Month(quarterStartMonth), 1, 0, 0, 0, 0, today.Location())

	// Get the last day of the quarter
	quarterEnd := quarterStart.AddDate(0, 3, -1)

	return quarterStart, quarterEnd
}

func monthToQuarter(month int) int {
	switch month {
	case 1, 2, 3:
		return 1
	case 4, 5, 6:
		return 2
	case 7, 8, 9:
		return 3
	case 10, 11, 12:
		return 4
	default:
		return 0
	}
}

func SortData(data map[string]models.CodedData) []models.CodedData {
	// Create a slice to hold the data

	var sortedData []models.CodedData

	// Iterate over the data and add it to the slice
	for _, v := range data {
		sortedData = append(sortedData, v)
	}

	// Sort the data by revenue
	sort.Slice(sortedData, func(i, j int) bool {
		return sortedData[i].Revenue > sortedData[j].Revenue
	})
	return sortedData
}

func CombineData(data []models.CodedData) []models.CodedData {
	// take the top 9 and combine the rest into one
	const NUM_ELEMENTS = 9

	// Create a slice to hold the data
	var combinedData []models.CodedData

	// Iterate over the data and add it to the slice
	for i, v := range data {
		if i <= NUM_ELEMENTS {
			combinedData = append(combinedData, v)
		} else {
			combinedData[NUM_ELEMENTS].Name += ", " + v.Name
			combinedData[NUM_ELEMENTS].Revenue += v.Revenue
			combinedData[NUM_ELEMENTS].Count += v.Count
		}
	}

	return combinedData
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

	// Start with today and iterate backwards until we hit Sunday
	myDay := today

	// get the day of the week
	day := myDay.Weekday()

	// get Sunday as the start of the week
	for day != time.Sunday {
		myDay = myDay.AddDate(0, 0, -1)
		day = myDay.Weekday()
	}

	return myDay
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

func SortVacationData(vacationData []models.VacationHours) []models.VacationHours {
	// Define the sorting function
	sortingFunc := func(i, j int) bool {
		amountI, okI := strconv.ParseFloat(vacationData[i].AmountDue, 64)
		amountJ, okJ := strconv.ParseFloat(vacationData[j].AmountDue, 64)

		if okI != nil || okJ != nil {
			return false
		}

		return amountI > amountJ
	}

	// Sort the vacationData based on the start date
	sort.Slice(vacationData, sortingFunc)

	// Return the sorted vacationData
	return vacationData
}

func FindLatestDateFromRevenueData(data []models.WeeklyRevenue) (time.Time, error) {
	// Define the latest date as the zero value
	today := time.Now()
	year, week := today.ISOWeek()

	weekAmt := len(data)

	if weekAmt == 0 {
		// No data
		return time.Time{}, fmt.Errorf("No data found")
	} else if weekAmt < 52 || weekAmt > 53 {
		// Only one week of data
		return time.Time{}, fmt.Errorf("Incorrect week count")
	}
	for i := year; i >= 2020; i-- {
		for j := week; j >= 1; j-- {
			revenueStruct := data[j]
			rev := revenueStruct.GetRevenue(i)

			if rev == nil {
				// no data found
				// continue with the a previous week
				continue
			} else if *rev > 0 {
				// data found
				// create a time object with the current week and year. returns a monday
				fmt.Printf("Week: %d Year: %d Revenue: %f\n", j, i, *rev)
				return FindSundayfromWeek(j, i), nil
			}

		}
		week = 52
	}
	return time.Time{}, fmt.Errorf("No data found")
}

func FindSundayfromWeek(week int, year int) time.Time {
	// Start with the first day of the year
	date := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	isoYear, isoWeek := date.ISOWeek()

	// Adjust to the first Monday of the specified week
	for isoYear < year || (isoYear == year && isoWeek < week) {
		date = date.AddDate(0, 0, 1) // Move forward one day
		isoYear, isoWeek = date.ISOWeek()
	}

	// Adjust to find the Sunday of that week
	return date.AddDate(0, 0, -1) // Go back one day to get Sunday
}

// DateRange represents a range of dates.
type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
	Week      int
	Amount    float64
}

// GenerateDateRanges generates a slice of DateRange structs from the start date (Sunday) to today.
func GenerateDateRanges(startDate time.Time) []*DateRange {
	var dateRanges []*DateRange
	today := time.Now()

	// Adjust the start date to be the Sunday of that week (if needed)
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	// Loop until we reach today, creating a DateRange for each week
	for startDate.Before(today) || startDate.Equal(today) {
		endDate := startDate.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second) // End date is the following Saturday at 23:59:59
		_, weekNumber := startDate.ISOWeek()

		dateRanges = append(dateRanges, &DateRange{
			StartDate: startDate,
			EndDate:   endDate,
			Week:      weekNumber,
		})
		startDate = startDate.AddDate(0, 0, 7) // Move to the next week
		weekNumber++
	}

	return dateRanges
}

func UpdateWeeklyRevenue(weeklyRevenues []models.WeeklyRevenue, dateRanges []*DateRange) {
	for _, dateRange := range dateRanges {
		year, weekNumber := dateRange.StartDate.ISOWeek() // Get the ISO week number

		// Update the revenue field based on the year
		switch year {
		case 2021:
			weeklyRevenues[weekNumber-1].Revenue2021 = new(float64)
			*weeklyRevenues[weekNumber-1].Revenue2021 = dateRange.Amount
		case 2022:
			weeklyRevenues[weekNumber-1].Revenue2022 = new(float64)
			*weeklyRevenues[weekNumber-1].Revenue2022 = dateRange.Amount
		case 2023:
			weeklyRevenues[weekNumber-1].Revenue2023 = new(float64)
			*weeklyRevenues[weekNumber-1].Revenue2023 = dateRange.Amount
		case 2024:
			weeklyRevenues[weekNumber-1].Revenue2024 = new(float64)
			*weeklyRevenues[weekNumber-1].Revenue2024 = dateRange.Amount
		default:
			fmt.Println("Year not in range")
			continue // Skip if the year is not in the range
		}
	}

	// if there are 53 weeks then we're going to remove it and apeend the results to the next year
	if len(weeklyRevenues) == 53 {
		// remove the last week
		endOfYear := weeklyRevenues[52]
		weeklyRevenues = weeklyRevenues[:52]

		// remove 2021 week one
		weeklyRevenues[0].Revenue2021 = nil

		if endOfYear.Revenue2021 != nil {
			*weeklyRevenues[0].Revenue2022 += *endOfYear.Revenue2021
		}

		if endOfYear.Revenue2022 != nil {
			*weeklyRevenues[0].Revenue2023 = *endOfYear.Revenue2022
		}

		if endOfYear.Revenue2023 != nil {
			*weeklyRevenues[0].Revenue2024 += *endOfYear.Revenue2023
		}

	}

	fmt.Println("Updated weekly revenues")

}

func CombineStackedMilesData(when string, data []models.StackedMilesData) []models.StackedMilesData {
	// Create a map to hold the data
	if when == "week_to_date" {
		// Create a map to hold the data
		weekData := make(map[int]models.StackedMilesData)
		timeLayout := "2006-01-02T15:04:05Z"

		// Iterate over the data and add it to the map
		for _, v := range data {
			// get the day number from the date
			parsedTime, err := time.Parse(timeLayout, v.Date)
			if err != nil {
				fmt.Println("Error parsing time: ", err)
				continue
			}

			day := parsedTime.Day()

			// Check if the day is already in the map
			if d, ok := weekData[day]; ok {
				// Add the miles to the existing data

				d.EmptyMiles += v.EmptyMiles
				d.LoadedMiles += v.LoadedMiles
				weekData[day] = d
			} else {
				// Add the data to the map
				weekData[day] = v
			}
		}

		// place the data in ascending order
		keys := make([]int, 0, len(weekData))
		for k := range weekData {
			keys = append(keys, k)
		}

		sort.Ints(keys)

		// Create a slice to hold the data
		sortedData := make([]models.StackedMilesData, 0, len(keys))
		for _, k := range keys {
			sortedData = append(sortedData, weekData[k])
		}
		return sortedData
	}

	if when == "month_to_date" {
		// Create a map to hold the data
		weekData := make(map[int]models.StackedMilesData)
		timeLayout := "2006-01-02T15:04:05Z"

		// Iterate over the data and add it to the map
		for _, v := range data {
			// get the day number from the date
			parsedTime, err := time.Parse(timeLayout, v.Date)
			if err != nil {
				fmt.Println("Error parsing time: ", err)
				continue
			}

			_, week := parsedTime.ISOWeek()

			// Check if the day is already in the map
			if d, ok := weekData[week]; ok {
				// Add the miles to the existing data
				d.EmptyMiles += v.EmptyMiles
				d.LoadedMiles += v.LoadedMiles
				weekData[week] = d

			} else {
				// Add the data to the map
				weekData[week] = v
			}
		}

		// place the data in ascending order
		keys := make([]int, 0, len(weekData))
		for k := range weekData {
			keys = append(keys, k)
		}

		slices.Sort(keys)

		// create another final slice to hold the data sorted by the weeks
		sortedData := make([]models.StackedMilesData, 0, len(keys))
		// Create a slice to hold the data
		for _, k := range keys {
			sortedData = append(sortedData, weekData[k])
		}

		return sortedData
	}
	if when == "quarter" {
		// Create a map to hold the data
		fmt.Println("unimplemented")
		return nil
	}

	fmt.Println("Error: Invalid when parameter")
	return nil
}

func StackedToMilesData(timeframe string, data []models.StackedMilesData) []models.MilesData {
	var milesData []models.MilesData
	var timeframestr, section string

	for i, v := range data {

		if v.LoadedMiles < 1 {
			// skip if there are no loaded miles
			continue
		}

		if timeframe == "week_to_date" {
			timeframe = "week"
			// get the day number from the date
			parsedTime, err := time.Parse("2006-01-02T15:04:05Z", v.Date)
			if err != nil {
				fmt.Println("Error parsing time: ", err)
				continue
			}
			timeframestr = strconv.Itoa(parsedTime.Day())
			section = "Day"
		}

		if timeframe == "month_to_date" {
			timeframe = "month"
			// we are doing weekly data
			time, error := time.Parse("2006-01-02T15:04:05Z", v.Date)
			if error != nil {
				fmt.Println("Error parsing time: ", error)
				continue
			}
			_, week := time.ISOWeek()

			timeframestr = strconv.Itoa(week)
			section = "Week"
		}

		totalmiles := v.EmptyMiles + v.LoadedMiles
		percentEmpty := ((v.LoadedMiles - v.EmptyMiles) / v.LoadedMiles) * 100
		milesData = append(milesData, models.MilesData{
			Name:             "time " + timeframestr,
			NameStr:          section + strconv.Itoa(i+1),
			DeliveryDate:     v.Date,
			TotalEmptyMiles:  v.EmptyMiles,
			TotalLoadedMiles: v.LoadedMiles,
			TotalMiles:       totalmiles,
			PercentEmpty:     percentEmpty,
		})
	}

	return milesData
}
