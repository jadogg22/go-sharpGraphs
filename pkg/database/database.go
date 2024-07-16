package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			err := godotenv.Load(filepath.Join(dir, ".env"))
			if err != nil {
				log.Fatal("Error loading .env file")
			}
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal(".env file not found")
		}
		dir = parent
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

func GetCachedData(db *sql.DB, company string) ([]models.WeeklyRevenue, error) {
	var dbTable string
	switch company {
	case "transportation":
		dbTable = "trans_year_rev"
	case "logistics":
		dbTable = "log_year_rev"
	default:
		return nil, fmt.Errorf("invalid company: %s", company)
	}

	data, err := FetchRevenueDataToWeeklyRevenue(db, dbTable)
	if err != nil {
		return nil, fmt.Errorf("error fetching revenue data: %w", err)
	}

	return data, nil
}

// refactored code
func FetchRevenueDataToWeeklyRevenue(db *sql.DB, dbTable string) ([]models.WeeklyRevenue, error) {
	query := fmt.Sprintf(`
        SELECT Year, Week, TotalRevenue
        FROM %s
        WHERE Year BETWEEN $1 AND $2
        ORDER BY Week, Year
    `, dbTable)

	currentYear := 2024
	startYear := currentYear - 4

	rows, err := db.Query(query, startYear, currentYear)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	data := make([]models.WeeklyRevenue, 52)
	for i := range data {
		data[i] = models.WeeklyRevenue{Name: i + 1}
	}

	for rows.Next() {
		var year, week int
		var rev sql.NullFloat64
		if err := rows.Scan(&year, &week, &rev); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		if rev.Valid {
			revValue := rev.Float64
			switch year {
			case 2021:
				data[week-1].Revenue2021 = &revValue
			case 2022:
				data[week-1].Revenue2022 = &revValue
			case 2023:
				data[week-1].Revenue2023 = &revValue
			case 2024:
				data[week-1].Revenue2024 = &revValue
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return data, nil
}

func GetYearByYearDataRefactored(db *sql.DB, data []models.WeeklyRevenue, company string) ([]models.WeeklyRevenue, error) {
	currentYear, currentWeek := time.Now().ISOWeek()
	if company != "transportation" && company != "logistics" {
		return nil, fmt.Errorf("invalid company")
	}

	missingData, err := FindMissingData(data)
	if err != nil {
		return nil, fmt.Errorf("error finding missing data: %w", err)
	}

	query := fmt.Sprintf("SELECT COALESCE(SUM(TotalRevenue), 0) FROM %s WHERE Week = $1", company)
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, missing := range missingData {
		// Check if we're not trying to update future weeks
		if missing.Year > currentYear || (missing.Year == currentYear && missing.Week > currentWeek) {
			continue
		}

		var totalRevenue sql.NullFloat64
		err := stmt.QueryRow(missing.WeekID).Scan(&totalRevenue)
		if err != nil {
			if err == sql.ErrNoRows {
				// No data for this week, skip it
				continue
			}
			return nil, fmt.Errorf("error querying data for week %s: %w", missing.WeekID, err)
		}

		// Skip weeks with no revenue (NULL or 0)
		if !totalRevenue.Valid || totalRevenue.Float64 == 0 {
			continue
		}

		revenue := totalRevenue.Float64

		updated := false
		for i, entry := range data {
			if entry.Name == missing.Week {
				switch missing.Year {
				case 2021:
					data[i].Revenue2021 = &revenue
				case 2022:
					data[i].Revenue2022 = &revenue
				case 2023:
					data[i].Revenue2023 = &revenue
				case 2024:
					data[i].Revenue2024 = &revenue
				}
				updated = true
				break
			}
		}

		if !updated {
			// Week doesn't exist in data slice, append new entry
			newEntry := models.WeeklyRevenue{Name: missing.Week}
			switch missing.Year {
			case 2021:
				newEntry.Revenue2021 = &revenue
			case 2022:
				newEntry.Revenue2022 = &revenue
			case 2023:
				newEntry.Revenue2023 = &revenue
			case 2024:
				newEntry.Revenue2024 = &revenue
			}
			data = append(data, newEntry)
		}
	}

	return data, nil
}

type MissingDataPoint struct {
	Year   int    `json:"year"`
	Week   int    `json:"week"`
	WeekID string `json:"weekID"`
}

func FindMissingData(data []models.WeeklyRevenue) ([]MissingDataPoint, error) {
	currentYear, currentWeek := time.Now().ISOWeek()
	var missingData []MissingDataPoint

	for year := 2021; year <= currentYear; year++ {
		endWeek := 52
		if year == currentYear {
			endWeek = currentWeek
		}

		for week := 1; week <= endWeek; week++ {
			weekID := fmt.Sprintf("%d W%02d", year, week)

			found := false
			for _, entry := range data {
				if entry.Name == week { // Note: Comparing with week, not weekID
					found = true
					isMissing := false
					switch year {
					case 2021:
						isMissing = entry.Revenue2021 == nil
					case 2022:
						isMissing = entry.Revenue2022 == nil
					case 2023:
						isMissing = entry.Revenue2023 == nil
					case 2024:
						isMissing = entry.Revenue2024 == nil
					}
					if isMissing {
						missingData = append(missingData, MissingDataPoint{
							Year:   year,
							Week:   week,
							WeekID: weekID,
						})
					}
					break
				}
			}

			if !found {
				missingData = append(missingData, MissingDataPoint{
					Year:   year,
					Week:   week,
					WeekID: weekID,
				})
			}
		}
	}

	return missingData, nil
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
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // This will be a no-op if the transaction is committed successfully

	if company != "transportation" && company != "logistics" {
		return fmt.Errorf("invalid company")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (
            RevenueCode, OrderNumber, OrderType, Freight, FuelSurcharge, RemainingCharges,
            TotalRevenue, BillMiles, LoadedMiles, EmptyMiles, TotalMiles, EmptyPercentage,
            RevLoadedMile, RevTotalMile, DeliveryDate, Origin, Destination, Customer,
            CustomerCategory, OperationsUser, Billed, ControllingParty, Commodity,
            TrailerType, OriginState, DestinationState, Week, Month, Quarter, Brokered
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
            $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
        ) ON CONFLICT (OrderNumber) DO NOTHING;`, company)

	// Prepare the INSERT statement
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert or ignore each LoadData object into the database
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
			return fmt.Errorf("failed to execute statement: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %w", err)
	}

	return nil
}
