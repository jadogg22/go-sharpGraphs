package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Make_connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "Production.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS daily_driver_data (
        dispatcher TEXT,
        deadhead_percent REAL,
        freight REAL,
        fuel_surcharge REAL,
        remain_chgs REAL,
        revenue REAL,
        total_rev_per_rev_miles REAL,
        total_rev_per_total_miles REAL,
        average_weekly_rev REAL,
        average_weekly_rev_miles REAL,
        average_rev_miles REAL,
        revenue_miles REAL,
        total_miles REAL,
		trucks INTEGER,
        date TEXT
    )`)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

// Add_DailyDriverData adds a DailyDriverData to the database
func Add_DailyDriverData(db *sql.DB, dailyData DailyDriverData) error {
	// Add the dailyDriverData to the db
	_, err := db.Exec(`INSERT INTO daily_driver_data (dispatcher, deadhead_percent, freight, fuel_surcharge, remain_chgs, revenue, total_rev_per_rev_miles, total_rev_per_total_miles, average_weekly_rev, average_weekly_rev_miles, average_rev_miles, revenue_miles, total_miles, trucks, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dailyData.Dispatcher, dailyData.Deadhead_percent, dailyData.Freight, dailyData.Fuel_Surcharge, dailyData.Remain_Chgs, dailyData.Revenue, dailyData.Total_Rev_per_rev_miles, dailyData.Total_Rev_per_Total_Miles, dailyData.Average_weekly_rev, dailyData.Average_weekly_Rev_Miles, dailyData.Average_rev_miles, dailyData.Revenue_Miles, dailyData.Total_Miles, dailyData.Trucks, dailyData.Date)
	if err != nil {
		return err
	}

	return nil
}

func GetDispacherDataFromDB(db *sql.DB) ([]DriverData, error) {
	query := `
        SELECT
            Dispatcher,
            Deadhead_percent,
            Freight,
            Fuel_Surcharge,
            Remain_Chgs,
            Revenue,
            Total_Rev_per_rev_miles,
            Total_Rev_per_Total_Miles,
            Average_weekly_rev,
            Average_weekly_Rev_Miles,
            Average_rev_miles,
            Revenue_Miles,
            Total_Miles,
		Trucks, 
            Date
        FROM (
            SELECT
                *,
                ROW_NUMBER() OVER (PARTITION BY Dispatcher ORDER BY Date DESC) AS RowNum
            FROM
                daily_driver_data
        ) AS RankedData
        WHERE
            RowNum <= 3
        ORDER BY
            Dispatcher,
            Date DESC;
    `

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer rows.Close()

	dispatcherData := make(map[string]DriverData)

	for rows.Next() {
		var dispatcher string
		var dateStr string
		var data DailyDriverData
		err := rows.Scan(
			&dispatcher,
			&data.Deadhead_percent,
			&data.Freight,
			&data.Fuel_Surcharge,
			&data.Remain_Chgs,
			&data.Revenue,
			&data.Total_Rev_per_rev_miles,
			&data.Total_Rev_per_Total_Miles,
			&data.Average_weekly_rev,
			&data.Average_weekly_Rev_Miles,
			&data.Average_rev_miles,
			&data.Revenue_Miles,
			&data.Total_Miles,
			&data.Trucks,
			&dateStr,
		)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}

		date, err := time.Parse("2006-01-02 00:00:00+00:00", dateStr)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}

		if existingData, ok := dispatcherData[dispatcher]; ok {
			// Append data to existing entry
			existingData.Deadhead_percent = append(existingData.Deadhead_percent, data.Deadhead_percent)
			existingData.Freight = append(existingData.Freight, data.Freight)
			existingData.Fuel_Surcharge = append(existingData.Fuel_Surcharge, data.Fuel_Surcharge)
			existingData.Remain_Chgs = append(existingData.Remain_Chgs, data.Remain_Chgs)
			existingData.Revenue = append(existingData.Revenue, data.Revenue)
			existingData.Total_Rev_per_rev_miles = append(existingData.Total_Rev_per_rev_miles, data.Total_Rev_per_rev_miles)
			existingData.Total_Rev_per_Total_Miles = append(existingData.Total_Rev_per_Total_Miles, data.Total_Rev_per_Total_Miles)
			existingData.Average_weekly_rev = append(existingData.Average_weekly_rev, data.Average_weekly_rev)
			existingData.Average_weekly_Rev_Miles = append(existingData.Average_weekly_Rev_Miles, data.Average_weekly_Rev_Miles)
			existingData.Average_rev_miles = append(existingData.Average_rev_miles, data.Average_rev_miles)
			existingData.Revenue_Miles = append(existingData.Revenue_Miles, data.Revenue_Miles)
			existingData.Total_Miles = append(existingData.Total_Miles, data.Total_Miles)
			existingData.Trucks = append(existingData.Trucks, data.Trucks)
			existingData.Date = append(existingData.Date, date)
			dispatcherData[dispatcher] = existingData
		} else {
			// Create new entry
			dispatcherData[dispatcher] = DriverData{
				Dispatcher:                dispatcher,
				Deadhead_percent:          []float64{data.Deadhead_percent},
				Freight:                   []float64{data.Freight},
				Fuel_Surcharge:            []float64{data.Fuel_Surcharge},
				Remain_Chgs:               []float64{data.Remain_Chgs},
				Revenue:                   []float64{data.Revenue},
				Total_Rev_per_rev_miles:   []float64{data.Total_Rev_per_rev_miles},
				Total_Rev_per_Total_Miles: []float64{data.Total_Rev_per_Total_Miles},
				Average_weekly_rev:        []float64{data.Average_weekly_rev},
				Average_weekly_Rev_Miles:  []float64{data.Average_weekly_Rev_Miles},
				Average_rev_miles:         []float64{data.Average_rev_miles},
				Revenue_Miles:             []float64{data.Revenue_Miles},
				Total_Miles:               []float64{data.Total_Miles},
				Trucks:                    []int64{data.Trucks},
				Date:                      []time.Time{date},
			}
		}
	}

	// Convert map to slice
	var result []DriverData
	for _, data := range dispatcherData {
		// Remove KEVIN BOYDSTUN and add the rest to the result
		if data.Dispatcher != "KEVIN BOYDSTUN" && data.Dispatcher != "ROCHELLE GENERA" && data.Dispatcher != "STEPHANIE BINGHAM" {
			result = append(result, data)
		}
	}

	return result, nil
}

func GetYearByYearData(db *sql.DB, company string) ([]map[string]interface{}, error) {
	var query string
	var dbTable string

	if company != "transportation" && company != "logistics" {
		fmt.Println("this aint no company")
		return nil, fmt.Errorf("invalid company")
	}

	if company == "transportation" {
		dbTable = "trans_year_rev"
	} else {
		dbTable = "log_year_rev"
		//check if the table exists
	}

	query1 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", dbTable)
	query = query1 + `
        ID INTEGER PRIMARY KEY AUTOINCREMENT,
        Year INTEGER NOT NULL,
        Week INTEGER NOT NULL,
        TotalRevenue REAL NOT NULL
		);`

	_, err := db.Exec(query)
	if err != nil {
		// Handle error
		fmt.Println("we;re ducked")
		return nil, err
	}

	query = fmt.Sprintf("SELECT TotalRevenue FROM %s WHERE Year = ? AND Week = ?", dbTable)
	// Prepare the SQL statement for querying revenue da
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing SQL statement: %v", err)
	}
	defer stmt.Close()

	// Initialize the data slice to hold the results
	data := make([]map[string]interface{}, 0)

	// Get the current year
	currentYear := 2023

	// Iterate over the past 4 years

	// Iterate over the weeks
	for week := 1; week <= 52; week++ {
		mapData := make(map[string]interface{})
		mapData["Name"] = week

		for year := currentYear - 3; year <= currentYear; year++ {
			// Execute the query
			row := stmt.QueryRow(year, week)
			var rev float64
			err := row.Scan(&rev)
			// Only add the revenue if there was no error
			if err == nil {
				mapData[strconv.Itoa(year)+" Revenue"] = rev
			}

		}
		data = append(data, mapData) // Append inside the loop

	}

	return data, nil
}

// This function is pretty multi-faceted. So I'm going to break it down into steps
// 1. Calculate what data is needed from db
// 2. Get the newest week of data from the db
// 3. If the newest week is not the current week, update the db
// 4. return the newest data

func GetNewestYearByYearData(db *sql.DB, data []map[string]interface{}, company string) ([]map[string]interface{}, error) {

	// now we can get the Week in the db that its format is 2024 W01
	currenYear, currentWeek := getYearAndWeek()
	newestDataWeek := getNewestWeek(data)

	newestDataYear := currenYear
	// if the newest data week is greter then the current week, its a different year!
	if newestDataWeek > currentWeek {
		newestDataYear = currenYear - 1
	}

	if newestDataWeek-currentWeek >= 2 {
		// Need to update the yearRevenue db
		fmt.Println("Need to update the yearRevenue db")
		fmt.Println("TODO: Update the yearRevenue db")
	}

	// convert the week and year to Week string thats in db
	newestDataWeekStr := strconv.Itoa(newestDataYear) + " W" + fmt.Sprintf("%02d", newestDataWeek)
	currentDataWeekStr := strconv.Itoa(currenYear) + " W" + fmt.Sprintf("%02d", currentWeek)

	fmt.Println("Getting data between: ", newestDataWeekStr, " and ", currentDataWeekStr)

	// queary to get the total revenue between the two weeks

	query := `SELECT Week, ROUND(SUM(TotalRevenue), 2) as TotalRevenue FROM transportation WHERE Week BETWEEN ? AND ? GROUP BY Week`

	//get all the rows between the two weeks
	rows, err := db.Query(query, newestDataWeekStr, currentDataWeekStr)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		// get the week and revenue
		var week string
		var revenue float64
		err := rows.Scan(&week, &revenue)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}

		// update the slice of the data with the new revenue numbers
		// convert the week into a year and week and then we can go to the specific index
		// "2024 W01" -> 2024, 1
		currenYear, weekint, err := GetYearAndWeekFromStr(week)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}

		// no outa bounds here.
		if weekint-1 <= len(data) {
			data[weekint-1][strconv.Itoa(currenYear)+" Revenue"] = revenue
		}
	}

	// return the updated data
	return data, nil
}

func CreateYearlyRevenueRecord(db *sql.DB, year int, week int, revenue float64) error {
	// create table if it doesn't exist
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS year_by_year (
	Name TEXT,
	Week TEXT,
	Month TEXT,
	Quarter TEXT,
	Year INTEGER,
	Revenue Real
)`)
	if err != nil {
		return err
	}

	quarter := GetQuarter(week, year)
	month := GetMonth(week, year)

	// Add the revenue to the db
	_, err = db.Exec(`INSERT INTO year_by_year (Name, Week, Month, Quarter, Year, Revenue) VALUES (?, ?, ?, ?, ?, ?)`, week, week, month, quarter, year, revenue)
	if err != nil {
		return err
	}

	return nil

}

func GetCodedRevenueData(conn *sql.DB, when string) ([]map[string]interface{}, float64, int64, error) {
	// Query to retrieve revenue data grouped by RevenueCode
	query := `
    SELECT RevenueCode,
           COUNT(*) AS OrderCount,
           SUM(TotalRevenue) AS TotalRevenue
    FROM transportation
    WHERE Week BETWEEN ? AND ?
    GROUP BY RevenueCode

`
	when2 := GetNewWeekFromWeek(when, 4)

	fmt.Println("calling database fron", when, when2)

	//when2 is actually 4 weeks before
	rows, err := conn.Query(query, when2, when)
	if err != nil {
		return nil, 0.0, 0, err
	}
	defer rows.Close()

	// Map to store array of data
	revenueArray := []map[string]interface{}{}
	var totalRevenue float64 = 0.0 // this is the total revenue used for calculating percentages
	var totalCount int64 = 0       // this is the total orders

	// Iterate through the rows and populate the revenueMap
	for rows.Next() {
		revenueMap := make(map[string]interface{})
		var revenueCode string
		var revenue float64
		var count int64
		if err := rows.Scan(&revenueCode, &count, &revenue); err != nil {
			return nil, 0.0, 0, err
		}

		totalRevenue += revenue
		totalCount += count

		revenueMap["Code"] = revenueCode
		revenueMap["Revenue"] = revenue
		revenueMap["Count"] = count

		revenueArray = append(revenueArray, revenueMap)

	}
	if err := rows.Err(); err != nil {
		return nil, 0.0, 0, err
	}

	// sort the data
	revenueArray = SortArraybyRevenue(revenueArray)

	//filter data - make the "other" group
	filteredRevenueArray, err := FilterRevenueArray(revenueArray, totalRevenue, totalCount)
	if err != nil {
		return revenueArray, totalRevenue, totalCount, err
	}

	return filteredRevenueArray, totalRevenue, totalCount, nil
}

func GetMilesData(conn *sql.DB, when, company string) ([]MilesData, error) {

	var query string
	if company != "transportation" && company != "logistics" {
		return nil, fmt.Errorf("This company doesnt exist")
	}

	switch when {
	case "week_to_data":
		today := time.Now()
		year, week := today.ISOWeek()

		startOfWeek := fmt.Sprintf("%d W%02d", year, week-1)
		endOfWeek := fmt.Sprintf("%d W%02d", year, week)

		query = fmt.Sprintf(`
			SELECT DeliveryDate,
				DeliveryDate as Name,
				strftime('%w', DeliveryDate) as NameStr,
				SUM(LoadedMiles) AS TotalLoadedMiles,
				SUM(EmptyMiles) AS TotalEmptyMiles,
				SUM(TotalMiles) AS TotalMiles,
				SUM(EmptyMiles) / SUM(TotalMiles) * 100 AS PercentEmpty
			FROM %s
			WHERE Week BETWEEN '%s' AND '%s'
			GROUP BY DeliveryDate;`, company, startOfWeek, endOfWeek)

	case "month_to_date":
		today := time.Now()
		year, month, _ := today.Date()

		startOfMonth := fmt.Sprintf("%d M%02d", year, int(month)-1)
		endOfMonth := fmt.Sprintf("%d M%02d", year, int(month))

		query = fmt.Sprintf(`
			SELECT DeliveryDate,
				Week as Name,
				strftime('%m', DeliveryDate) as NameStr,
				SUM(LoadedMiles) AS TotalLoadedMiles,
				SUM(EmptyMiles) AS TotalEmptyMiles,
				SUM(TotalMiles) AS TotalMiles,
				SUM(EmptyMiles) / SUM(TotalMiles) * 100 AS PercentEmpty
			FROM %s
			WHERE Month BETWEEN '%s' AND '%s'
			GROUP BY Week;`, company, startOfMonth, endOfMonth)

	default:
		return nil, fmt.Errorf("unsupported time period")
	}

	// Execute the query
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []MilesData
	for rows.Next() {
		var md MilesData
		err := rows.Scan(
			&md.Name,
			&md.NameStr,
			&md.TotalLoadedMiles,
			&md.TotalEmptyMiles,
			&md.TotalMiles,
			&md.PercentEmpty,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, md)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return results, nil
}

func AddOrderToDB(conn *sql.DB, loadData *[]LoadData, company string) error {
	// Begin a transaction
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("Failed to start transaction")
	}
	if company != "transportation" && company != "logistics" {
		return fmt.Errorf("Invalid company")
	}

	query := fmt.Sprintf("INSERT OR REPLACE INTO %s (RevenueCode, OrderNumber, OrderType, Freight, FuelSurcharge, RemainingCharges, TotalRevenue, BillMiles, LoadedMiles, EmptyMiles, TotalMiles, EmptyPercentage, RevLoadedMile, RevTotalMile, DeliveryDate, Origin, Destination, Customer, CustomerCategory, OperationsUser, Billed, ControllingParty, Commodity, TrailerType, OriginState, DestinationState, Week, Month, Quarter, Brokered) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", company)

	// Prepare the INSERT OR REPLACE statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare statment")
	}
	defer stmt.Close()

	// Insert or replace each LoadData object into the database
	for _, data := range *loadData {
		_, err := stmt.Exec(
			data.RevenueCode,
			data.Order,
			data.OrderType,
			data.Freight,
			data.FuelSurcharge,
			data.RemainingCharges,
			data.TotalRevenue,
			data.BillMiles,
			data.LoadedMiles,
			data.EmptyMiles,
			data.TotalMiles,
			data.EmptyPct,
			data.RevPerLoadedMile,
			data.RevPerTotalMile,
			data.DeliveryDate,
			data.Origin,
			data.Destination,
			data.Customer,
			data.CustomerCategory,
			data.OperationsUser,
			data.Billed,
			data.ControllingParty,
			data.Commodity,
			data.TrailerType,
			data.OriginState,
			data.DestinationState,
			data.Week,
			data.Month,
			data.Quarter,
			data.Brokered,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("Faild to do transaction")
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Faild to commit the transaction")
	}
	return nil
}
