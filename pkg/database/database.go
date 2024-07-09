package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"
	_ "github.com/lib/pq"

	_ "github.com/mattn/go-sqlite3"
	// get the env variables
	"github.com/joho/godotenv"
)

var (
	PG_HOST     string
	PG_PORT     string
	PG_USER     string
	PG_PASSWORD string
	PG_DATABASE string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PG_HOST = os.Getenv("PG_HOST")
	PG_PORT = os.Getenv("PG_PORT")
	PG_USER = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	PG_DATABASE = os.Getenv("PG_DATABASE")
}

func Make_connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "Data/Production.db")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func PG_Make_connection() (*sql.DB, error) {

	// grab env variables for db connection
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DATABASE)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return db, nil
}

func test_make_connection(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return db, nil
}

func isThereData(db *sql.DB) bool {
	// Query to check if there is data in the db
	query := `SELECT COUNT(*) FROM transportation;`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	defer rows.Close()

	// Iterate through the rows and check if there is data
	for rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			fmt.Println("Error: ", err)
			return false
		}
		if count > 0 {
			return true
		}
	}
	return false
}

// Add_DailyDriverData adds a DailyDriverData to the database
func Add_DailyDriverData(db *sql.DB, dailyData models.DailyDriverData) error {
	// Add the dailyDriverData to the db
	_, err := db.Exec(`INSERT INTO daily_driver_data (dispatcher, deadhead_percent, freight, fuel_surcharge, remain_chgs, revenue, total_rev_per_rev_miles, total_rev_per_total_miles, average_weekly_rev, average_weekly_rev_miles, average_rev_miles, revenue_miles, total_miles, trucks, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dailyData.Dispatcher, dailyData.Deadhead_percent, dailyData.Freight, dailyData.Fuel_Surcharge, dailyData.Remain_Chgs, dailyData.Revenue, dailyData.Total_Rev_per_rev_miles, dailyData.Total_Rev_per_Total_Miles, dailyData.Average_weekly_rev, dailyData.Average_weekly_Rev_Miles, dailyData.Average_rev_miles, dailyData.Revenue_Miles, dailyData.Total_Miles, dailyData.Trucks, dailyData.Date)
	if err != nil {
		return err
	}

	return nil
}

// Wea are going to show 2 graphs on the front end
// first graph
// Manager | # of trucks | Miles | Deadhead | Order | Stop
// secound graph
// Manager | Average MPTPD | Average RPTPD | DH% | Order OTP | Stop OTP | AVG MPTPD Needed to Make Goal
// I would like to have the data color coded so that we can show the user if they are doing well or not.
// Manager | Average MPTPD | Average RPTPD | DH% | Order OTP | Stop OTP | AVG MPTPD Needed to Make Goal
// then for each of them we can have a color code like item.AverageMPTPDCOlOR: "Green"
// Then on the front end we can have a function take the color and return the correct color.

func GetDispacherDataFromDB(db *sql.DB) ([]models.DriverData, error) {
	// Query to retrieve this weeks data for each dispatcher

	// TODO - combine Freight and Fueil_Surcharge into a single field called truckHire
	// TODO - subtract revenue from truckhire to get net revenue
	// TODO - add a field for net revenue called margin %

	// Get the date of the start of the work week
	startDay := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	startDayStr := startDay.Format("2006-01-02")

	fmt.Println("Getting data from: ", startDayStr)

	// for each dispatcher sum the data for the week # trucks miles deadead order and stop
	query := fmt.Sprintf(`
		SELECT dispatcher,
			SUM(trucks) as trucks,
			SUM(total_miles) as miles,
			SUM(deadhead_percent) as deadhead,
		FROM daily_driver_data
		WHERE date >= ?
		GROUP BY dispatcher
		ORDER BY dispatcher;
	`)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer rows.Close()

	// Map to store the data

	// Iterate through the rows and populate the dispatcherData map
	for rows.Next() {
		var dispatcher string
		var miles float64
		var deadhead float64
		var trucks int64
		err := rows.Scan(
			&dispatcher,
			&trucks,
			&miles,
			&deadhead,
		)
		if err != nil {
			fmt.Println("Error: ", err)
			return nil, err
		}
	}

	// for rows.Next() {
	// 	var dispatcher string
	// 	var dateStr string
	// 	var data models.DailyDriverData
	// 	err := rows.Scan(
	// 		&dispatcher,
	// 		&data.Deadhead_percent,
	// 		&data.Freight,
	// 		&data.Fuel_Surcharge,
	// 		&data.Remain_Chgs,
	// 		&data.Revenue,
	// 		&data.Total_Rev_per_rev_miles,
	// 		&data.Total_Rev_per_Total_Miles,
	// 		&data.Average_weekly_rev,
	// 		&data.Average_weekly_Rev_Miles,
	// 		&data.Average_rev_miles,
	// 		&data.Revenue_Miles,
	// 		&data.Total_Miles,
	// 		&data.Trucks,
	// 		&dateStr,
	// 	)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return nil, err
	// 	}

	// 	date, err := time.Parse("2006-01-02 00:00:00+00:00", dateStr)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return nil, err
	// 	}

	// 	if existingData, ok := dispatcherData[dispatcher]; ok {
	// 		// Append data to existing entry
	// 		existingData.Deadhead_percent = append(existingData.Deadhead_percent, data.Deadhead_percent)
	// 		existingData.Freight = append(existingData.Freight, data.Freight)
	// 		existingData.Fuel_Surcharge = append(existingData.Fuel_Surcharge, data.Fuel_Surcharge)
	// 		existingData.Remain_Chgs = append(existingData.Remain_Chgs, data.Remain_Chgs)
	// 		existingData.Revenue = append(existingData.Revenue, data.Revenue)
	// 		existingData.Total_Rev_per_rev_miles = append(existingData.Total_Rev_per_rev_miles, data.Total_Rev_per_rev_miles)
	// 		existingData.Total_Rev_per_Total_Miles = append(existingData.Total_Rev_per_Total_Miles, data.Total_Rev_per_Total_Miles)
	// 		existingData.Average_weekly_rev = append(existingData.Average_weekly_rev, data.Average_weekly_rev)
	// 		existingData.Average_weekly_Rev_Miles = append(existingData.Average_weekly_Rev_Miles, data.Average_weekly_Rev_Miles)
	// 		existingData.Average_rev_miles = append(existingData.Average_rev_miles, data.Average_rev_miles)
	// 		existingData.Revenue_Miles = append(existingData.Revenue_Miles, data.Revenue_Miles)
	// 		existingData.Total_Miles = append(existingData.Total_Miles, data.Total_Miles)
	// 		existingData.Trucks = append(existingData.Trucks, data.Trucks)
	// 		existingData.Date = append(existingData.Date, date)
	// 		dispatcherData[dispatcher] = existingData
	// 	} else {
	// 		// Create new entry
	// 		dispatcherData[dispatcher] = models.DriverData{
	// 			Dispatcher:                dispatcher,
	// 			Deadhead_percent:          []float64{data.Deadhead_percent},
	// 			Freight:                   []float64{data.Freight},
	// 			Fuel_Surcharge:            []float64{data.Fuel_Surcharge},
	// 			Remain_Chgs:               []float64{data.Remain_Chgs},
	// 			Revenue:                   []float64{data.Revenue},
	// 			Total_Rev_per_rev_miles:   []float64{data.Total_Rev_per_rev_miles},
	// 			Total_Rev_per_Total_Miles: []float64{data.Total_Rev_per_Total_Miles},
	// 			Average_weekly_rev:        []float64{data.Average_weekly_rev},
	// 			Average_weekly_Rev_Miles:  []float64{data.Average_weekly_Rev_Miles},
	// 			Average_rev_miles:         []float64{data.Average_rev_miles},
	// 			Revenue_Miles:             []float64{data.Revenue_Miles},
	// 			Total_Miles:               []float64{data.Total_Miles},
	// 			Trucks:                    []int64{data.Trucks},
	// 			Date:                      []time.Time{date},
	// 		}
	// 	}
	// }

	// // Convert map to slice
	// var result []models.DriverData
	// for _, data := range dispatcherData {
	// 	// Remove KEVIN BOYDSTUN and add the rest to the result
	// 	if data.Dispatcher != "KEVIN BOYDSTUN" && data.Dispatcher != "ROCHELLE GENERA" && data.Dispatcher != "STEPHANIE BINGHAM" {
	// 		result = append(result, data)
	// 	}
	// }

	// return result, nil

	return nil, nil
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
	currenYear, currentWeek := helpers.GetYearAndWeek()
	newestDataWeek := helpers.GetNewestWeek(data)

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
	if company != "transportation" && company != "logistics" {
		return nil, fmt.Errorf("company is not correct")
	}

	query := fmt.Sprintf(`SELECT Week, ROUND(SUM(TotalRevenue), 2) as TotalRevenue FROM %s WHERE Week BETWEEN ? AND ? GROUP BY Week`, company)

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
		currenYear, weekint, err := helpers.GetYearAndWeekFromStr(week)
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

	quarter := helpers.GetQuarter(week, year)
	month := helpers.GetMonth(week, year)

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
	when2 := helpers.GetNewWeekFromWeek(when, 4)

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
	revenueArray = helpers.SortArraybyRevenue(revenueArray)

	//filter data - make the "other" group
	filteredRevenueArray, err := helpers.FilterRevenueArray(revenueArray, totalRevenue, totalCount)
	if err != nil {
		return revenueArray, totalRevenue, totalCount, err
	}

	return filteredRevenueArray, totalRevenue, totalCount, nil
}

func GetMilesData(conn *sql.DB, when, company string) ([]models.MilesData, error) {

	var query string
	var startDate string
	var endDate string

	if company != "transportation" && company != "logistics" {
		return nil, fmt.Errorf("this company doesnt exist")
	}

	switch when {
	case "week_to_date":
		today := time.Now()
		year, week := today.ISOWeek()

		startDate = fmt.Sprintf("%d W%02d", year, week-1)
		endDate = fmt.Sprintf("%d W%02d", year, week)

		fmt.Println("Getting data between: ", startDate, " and ", endDate)

		query = fmt.Sprintf(`
			SELECT DeliveryDate,
				DeliveryDate as Name,
				strftime('%%w', DeliveryDate) as NameStr,
				SUM(LoadedMiles) AS TotalLoadedMiles,
				SUM(EmptyMiles) AS TotalEmptyMiles,
				SUM(TotalMiles) AS TotalMiles,
				SUM(EmptyMiles) / SUM(TotalMiles) * 100 AS PercentEmpty
			FROM %s
			WHERE Week BETWEEN ? AND ?
			GROUP BY DeliveryDate;`, company)

	case "month_to_date":
		today := time.Now()
		year, month, _ := today.Date()

		startDate = fmt.Sprintf("%d M%02d", year, int(month)-1)
		endDate = fmt.Sprintf("%d M%02d", year, int(month))

		query = fmt.Sprintf(`
        SELECT DeliveryDate,
            Week as Name,
            strftime('%%m', DeliveryDate) as NameStr,
            SUM(LoadedMiles) AS TotalLoadedMiles,
            SUM(EmptyMiles) AS TotalEmptyMiles,
            SUM(TotalMiles) AS TotalMiles,
            SUM(EmptyMiles) / SUM(TotalMiles) * 100 AS PercentEmpty
        FROM %s
        WHERE Month BETWEEN ? AND ?
        GROUP BY Name;`, company)

	default:
		return nil, fmt.Errorf("unsupported time period")

	}

	fmt.Println("Getting data between ", startDate, " and ", endDate)

	// Execute the query
	rows, err := conn.Query(query, startDate, endDate)
	if err != nil {
		println("Error making query: ", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.MilesData
	for rows.Next() {
		var md models.MilesData
		err := rows.Scan(
			&md.DeliveryDate,
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

func AddOrderToDB(conn *sql.DB, loadData *[]models.LoadData, company string) error {
	// Begin a transaction
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction")
	}
	if company != "transportation" && company != "logistics" {
		return fmt.Errorf("invalid company")
	}

	query := fmt.Sprintf("INSERT OR REPLACE INTO %s (RevenueCode, OrderNumber, OrderType, Freight, FuelSurcharge, RemainingCharges, TotalRevenue, BillMiles, LoadedMiles, EmptyMiles, TotalMiles, EmptyPercentage, RevLoadedMile, RevTotalMile, DeliveryDate, Origin, Destination, Customer, CustomerCategory, OperationsUser, Billed, ControllingParty, Commodity, TrailerType, OriginState, DestinationState, Week, Month, Quarter, Brokered) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (OrderNumber) DO NOTHING;", company)

	// Prepare the INSERT OR REPLACE statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statment")
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
			return fmt.Errorf("faild to do transaction")
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("faild to commit the transaction")
	}
	return nil
}
