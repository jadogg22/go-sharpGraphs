package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	START_YEAR  int

	DB *sql.DB
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			err := godotenv.Load(filepath.Join(dir, ".env"))
			if err != nil {
				log.Fatalf("Error loading .env file: %v", err)
			}
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatalf(".env file not found in current or parent directories")
		}
		dir = parent
	}

	PG_HOST = os.Getenv("PG_HOST")
	PG_PORT = os.Getenv("PG_PORT")
	PG_USER = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	PG_DATABASE = os.Getenv("PG_DATABASE")

	startYearStr := os.Getenv("START_YEAR")
	if startYearStr == "" {
		START_YEAR = 2020 // Default value
	} else {
		parsedYear, err := strconv.Atoi(startYearStr)
		if err != nil {
			log.Fatalf("Error parsing START_YEAR from .env: %v", err)
		}
		START_YEAR = parsedYear
	}

	DB, err = PG_Make_connection()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database.")
}

/*
Now that we have acesss to the mcloud data base this database is going to serve as a cache for computationaly expensive data. ie. the year by year data
for the trans and logistics line graphs. This will allow us to store the data in the database and not have to recalculate it every time the user requests at
get the request much quicker.
*/

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
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

// Add_DailyDriverData adds a DailyDriverData to the database
func Add_DailyDriverData(dailyData models.DailyDriverData) error {
	// Add the dailyDriverData to the db
	_, err := DB.Exec(`INSERT INTO daily_driver_data (dispatcher, deadhead_percent, freight, fuel_surcharge, remain_chgs, revenue, total_rev_per_rev_miles, total_rev_per_total_miles, average_weekly_rev, average_weekly_rev_miles, average_rev_miles, revenue_miles, total_miles, trucks, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dailyData.Dispatcher, dailyData.Deadhead_percent, dailyData.Freight, dailyData.Fuel_Surcharge, dailyData.Remain_Chgs, dailyData.Revenue, dailyData.Total_Rev_per_rev_miles, dailyData.Total_Rev_per_Total_Miles, dailyData.Average_weekly_rev, dailyData.Average_weekly_Rev_Miles, dailyData.Average_rev_miles, dailyData.Revenue_Miles, dailyData.Total_Miles, dailyData.Trucks, dailyData.Date)
	if err != nil {
		return err
	}

	return nil
}

func AddCacheData(rev float64, week int, year int) error {
	query := `INSERT INTO transportation_weekly_revenue (totalrevenue, week, year) VALUES ($1, $2, $3) ON CONFLICT (week, year) DO UPDATE SET totalrevenue = EXCLUDED.totalrevenue`
	_, err := DB.Exec(query, rev, week, year)
	if err != nil {
		return fmt.Errorf("failed to insert data: %w", err)
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

func GetDispacherDataFromDB() ([]models.DriverData, error) {
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

	rows, err := DB.Query(query)
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

	return nil, nil
}

func GetCachedData(company string) ([]models.WeeklyRevenue, error) {
	var dbTable string
	switch company {
	case "transportation":
		dbTable = "trans_weekly_rev"
	case "logistics":
		dbTable = "log_year_rev"
	default:
		return nil, fmt.Errorf("invalid company: %s", company)
	}

	data, err := FetchRevenueDataToWeeklyRevenue(DB, dbTable)
	if err != nil {
		return nil, fmt.Errorf("error fetching revenue data: %w", err)
	}

	return data, nil
}

func NewFetchMyCache() ([]models.WeeklyRevenue, int, error) {
	currentYear := time.Now().Year()
	startYear := START_YEAR

	query := "SELECT week AS week_number"
	for year := startYear; year <= currentYear; year++ {
		query += fmt.Sprintf(", SUM(CASE WHEN year = %d THEN totalrevenue ELSE 0 END) AS \"%d_revenue\"", year, year)
	}
	query += " FROM transportation_weekly_revenue GROUP BY week ORDER BY week;"

	rows, err := DB.Query(query)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing query: %v", err)
	}

	defer rows.Close()

	// Use a map to store data by week number for easier population
	weeklyDataMap := make(map[int]models.WeeklyRevenue)

	for rows.Next() {
		var weekNumber int
		revenues := make([]sql.NullFloat64, currentYear-startYear+1)
		scanArgs := make([]interface{}, len(revenues)+1)
		scanArgs[0] = &weekNumber
		for i := range revenues {
			scanArgs[i+1] = &revenues[i]
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning row: %v", err)
		}

		revenueMap := make(map[string]*float64)
		for i, year := 0, startYear; year <= currentYear; i, year = i+1, year+1 {
			if revenues[i].Valid {
				val := revenues[i].Float64
				revenueMap[fmt.Sprintf("%d Revenue", year)] = &val
			} else {
				revenueMap[fmt.Sprintf("%d Revenue", year)] = nil
			}
		}

		theWeek := models.WeeklyRevenue{
			Name:     weekNumber,
			Revenues: revenueMap,
		}
		weeklyDataMap[weekNumber] = theWeek
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating rows: %v", err)
	}

	// Now, construct the final slice, only including weeks up to the current week
	var data []models.WeeklyRevenue
	newestWeek := 0 // Initialize newestWeek

	for week := 1; week <= 53; week++ {
		if entry, ok := weeklyDataMap[week]; ok {
			data = append(data, entry)
			// Update newestWeek if this week has data for the current year
			if entry.Revenues[fmt.Sprintf("%d Revenue", currentYear)] != nil {
				newestWeek = week
			}
		} else {
			// If no data for this week, create an empty entry for the current year
			emptyRevenueMap := make(map[string]*float64)
			for year := startYear; year <= currentYear; year++ {
				emptyRevenueMap[fmt.Sprintf("%d Revenue", year)] = nil
			}
			data = append(data, models.WeeklyRevenue{
				Name:     week,
				Revenues: emptyRevenueMap,
			})
		}
	}

	return data, newestWeek, nil
}

// refactored code
func FetchRevenueDataToWeeklyRevenue(db *sql.DB, dbTable string) ([]models.WeeklyRevenue, error) {
	query := fmt.Sprintf(`
        SELECT year, week, totalrevenue
        FROM %s
        ORDER BY Week, Year
    `, dbTable)

	rows, err := db.Query(query)
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
			key := fmt.Sprintf("%d Revenue", year)
			if data[week-1].Revenues == nil {
				data[week-1].Revenues = make(map[string]*float64)
			}
			data[week-1].Revenues[key] = &revValue
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return data, nil
}

func UpdateMyDatabase(data []*helpers.DateRange) error {
	// Create table if it doesn't exist
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS transportation_weekly_revenue (totalrevenue REAL, week INTEGER, year INTEGER, UNIQUE(week, year))`) // Ensure table schema is defined
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback on error

	stmt, err := tx.Prepare(`INSERT INTO transportation_weekly_revenue (totalrevenue, week, year) VALUES ($1, $2, $3) ON CONFLICT (week, year) DO UPDATE SET totalrevenue = EXCLUDED.totalrevenue`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert data into the table
	for _, dateRange := range data {
		if dateRange.Amount < 1 {
			continue
		}
		year, week := dateRange.StartDate.ISOWeek()
		_, err := stmt.Exec(dateRange.Amount, week, year)
		if err != nil {
			return fmt.Errorf("error executing statement for week %d, year %d: %w", week, year, err)
		}
	}

	return tx.Commit()
}

func GetYearByYearDataRefactored(data []models.WeeklyRevenue, company string) ([]models.WeeklyRevenue, error) {
	currentYear, currentWeek := time.Now().ISOWeek()
	if company != "transportation" && company != "logistics" {
		return nil, fmt.Errorf("invalid company")
	}

	missingData, err := FindMissingData(data)
	if err != nil {
		return nil, fmt.Errorf("error finding missing data: %w", err)
	}

	query := fmt.Sprintf("SELECT COALESCE(SUM(TotalRevenue), 0) FROM %s WHERE Week = $1", company)
	stmt, err := DB.Prepare(query)
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
				key := fmt.Sprintf("%d Revenue", missing.Year)
				if data[i].Revenues == nil {
					data[i].Revenues = make(map[string]*float64)
				}
				data[i].Revenues[key] = &revenue
				updated = true
				break
			}
		}

		if !updated {
			// Week doesn't exist in data slice, append new entry
			newEntry := models.WeeklyRevenue{Name: missing.Week}
			key := fmt.Sprintf("%d Revenue", missing.Year)
			newEntry.Revenues = make(map[string]*float64)
			newEntry.Revenues[key] = &revenue
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

	for year := 2020; year <= currentYear; year++ {
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
					key := fmt.Sprintf("%d Revenue", year)
					if _, ok := entry.Revenues[key]; !ok {
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

// This function should be called about once a week when we're getting more then a weeks worth of data from from the
//mcloud databse, this will add a new entry to our "cache" database so we dont have to access as much from mcloud

// TODO - change to the new yearly data db-table
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

// This function is curretly not being uesed and we're going to have to switch it over to the getData pkg
// TODO - implement coded revenue data

func GetCodedRevenueData(when string) ([]map[string]interface{}, float64, int64, error) {
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
	rows, err := DB.Query(query, when2, when)
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

func formatDate(dateStr string) (string, error) {
	// Parse the input date string
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %v", err)
	}

	// Format the date in a more human-readable way
	return t.Format("Jan 2"), nil
}

func GetMilesData(when, company string) ([]models.MilesData, error) {
	var query string

	if company != "transportation" && company != "logistics" {
		return nil, fmt.Errorf("this company doesn't exist")
	}

	switch when {
	case "week_to_date":
		query = fmt.Sprintf(`
			SELECT 
				deliverydate,
				deliverydate::date AS name,
				EXTRACT(DOW FROM deliverydate::date) AS NameStr,
				SUM(loadedmiles) AS total_loaded_miles,
				SUM(emptymiles) AS total_empty_miles,
				SUM(totalmiles) AS total_miles,
			CASE 
				WHEN SUM(totalmiles) > 0 THEN 
					(SUM(emptymiles) / SUM(totalmiles)) * 100 
					ELSE 0 
				END AS percent_empty
			FROM 
				%s
			GROUP BY 
				deliverydate
			ORDER BY 
				deliverydate::date DESC
				LIMIT 7;`, company)

	case "month_to_date":
		query = fmt.Sprintf(`
    WITH last_6_weeks AS (
        SELECT DISTINCT week
        FROM %s
        ORDER BY week DESC
        LIMIT 6
    )
    SELECT 
        t.week,
        MIN(t.deliverydate::date) AS week_start,
        MAX(t.deliverydate::date) AS week_end,
        SUM(t.loadedmiles) AS total_loaded_miles,
        SUM(t.emptymiles) AS total_empty_miles,
        SUM(t.totalmiles) AS total_miles,
        CASE 
            WHEN SUM(t.totalmiles) > 0 THEN 
                (SUM(t.emptymiles) / SUM(t.totalmiles)) * 100 
            ELSE 0 
        END AS percent_empty
    FROM 
        %s t
    JOIN 
        last_6_weeks l6w ON t.week = l6w.week
    GROUP BY 
        t.week
    ORDER BY 
        MIN(t.deliverydate::date) DESC;`, company, company)

	default:
		return nil, fmt.Errorf("unsupported time period")
	}

	// Execute the query
	rows, err := DB.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no data found")
		}
		fmt.Println("Error making query:", err)
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

		// Format the date in a more human-readable way
		md.NameStr, err = formatDate(md.Name)
		if err != nil {
			return nil, err
		}
		results = append(results, md)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no data found, Please check the DB")
	}

	return results, nil
}

func GetWeeklyRevenueData() ([]models.WeeklyRevenue, error) {
	currentYear := time.Now().Year()
	startYear := START_YEAR

	query := "SELECT week AS week_number"
	for year := startYear; year <= currentYear; year++ {
		query += fmt.Sprintf(", SUM(CASE WHEN year = %d THEN totalrevenue ELSE 0 END) AS \"%d_revenue\"", year, year)
	}
	query += " FROM transportation_weekly_revenue GROUP BY week ORDER BY week;"

	log.Printf("Attempting to execute query: %s\n", query)
	rows, err := DB.Query(query)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return nil, err
	}

	defer rows.Close()

	// Use a map to store data by week number for easier population
	weeklyDataMap := make(map[int]models.WeeklyRevenue)

	rowCount := 0
	for rows.Next() {
		rowCount++
		var weekNumber int
		revenues := make([]*float64, currentYear-startYear+1)
		scanArgs := make([]interface{}, len(revenues)+1)
		scanArgs[0] = &weekNumber
		for i := range revenues {
			scanArgs[i+1] = &revenues[i]
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Printf("Error scanning row: %v\n", err)
			return nil, err
		}

		revenueMap := make(map[string]*float64)
		for i, year := 0, startYear; year <= currentYear; i, year = i+1, year+1 {
			revenueMap[fmt.Sprintf("%d Revenue", year)] = revenues[i]
		}

		theWeek := models.WeeklyRevenue{
			Name:     weekNumber,
			Revenues: revenueMap,
		}
		weeklyDataMap[weekNumber] = theWeek
	}
	log.Printf("Finished processing %d rows from query.\n", rowCount)

	// Now, construct the final slice, only including weeks up to the current week
	var data []models.WeeklyRevenue
	for week := 1; week <= 53; week++ { // Iterate only up to the current week
		if entry, ok := weeklyDataMap[week]; ok {
			data = append(data, entry)
		} else {
			// If no data for this week, create an empty entry for the current year
			emptyRevenueMap := make(map[string]*float64)
			for year := startYear; year <= currentYear; year++ {
				emptyRevenueMap[fmt.Sprintf("%d Revenue", year)] = nil // Explicitly set to nil for no data
			}
			data = append(data, models.WeeklyRevenue{
				Name:     week,
				Revenues: emptyRevenueMap,
			})
		}
	}

	return data, nil
}

func CheckTransWeeklyRevHealth() error {
	log.Println("Performing advanced trans_weekly_rev health check...")

	currentYear := time.Now().Year()
	yearsToCheck := []int{START_YEAR, currentYear}

	// Add a past year if START_YEAR is not the current year
	if START_YEAR < currentYear {
		yearsToCheck = append(yearsToCheck, currentYear-1)
	}

	weeksToCheck := []int{1, 20, 30, 40, 52} // Check various weeks

	for _, year := range yearsToCheck {
		for _, week := range weeksToCheck {
			// Ensure we don't check future weeks for the current year
			currentISOYear, currentISOWeek := time.Now().ISOWeek()
			if year == currentISOYear && week > currentISOWeek {
				continue
			}

			var totalRevenue sql.NullFloat64
			query := "SELECT totalrevenue FROM transportation_weekly_revenue WHERE year = $1 AND week = $2"
			err := DB.QueryRow(query, year, week).Scan(&totalRevenue)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Printf("Health check: No data found for Year=%d, Week=%d", year, week)
					return fmt.Errorf("health check failed: no data for Year=%d, Week=%d", year, week)
				} else {
					log.Printf("Health check: Error querying for Year=%d, Week=%d: %v", year, week, err)
					return fmt.Errorf("health check failed: query error for Year=%d, Week=%d: %w", year, week, err)
				}
			}

			if !totalRevenue.Valid || totalRevenue.Float64 <= 0 {
				log.Printf("Health check: Invalid or zero revenue for Year=%d, Week=%d (Value: %v)", year, week, totalRevenue)
				return fmt.Errorf("health check failed: invalid or zero revenue for Year=%d, Week=%d", year, week)
			}
			log.Printf("Health check: Data found for Year=%d, Week=%d, Revenue=%.2f", year, week, totalRevenue.Float64)
		}
	}

	log.Println("Advanced trans_weekly_rev health check passed.")
	return nil
}

func FlushTransWeeklyRevTable() error {
	// Ensure the table exists before truncating
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS transportation_weekly_revenue (totalrevenue REAL, week INTEGER, year INTEGER, UNIQUE(week, year))`)
	if err != nil {
		return fmt.Errorf("error ensuring table exists before truncation: %w", err)
	}

	_, err = DB.Exec("TRUNCATE TABLE transportation_weekly_revenue;")
	if err != nil {
		return fmt.Errorf("error truncating transportation_weekly_revenue table: %w", err)
	}
	fmt.Println("transportation_weekly_revenue table truncated successfully.")
	return nil
}

// DumpTransWeeklyRevTable temporarily dumps the contents of the table to logs.
func DumpTransWeeklyRevTable() {
	log.Println("Dumping trans_weekly_rev table contents...")
	rows, err := DB.Query("SELECT year, week, totalrevenue FROM trans_weekly_rev ORDER BY year, week")
	if err != nil {
		log.Printf("Error dumping trans_weekly_rev table: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var year, week int
		var totalRevenue sql.NullFloat64
		if err := rows.Scan(&year, &week, &totalRevenue); err != nil {
			log.Printf("Error scanning row while dumping table: %v", err)
			continue
		}
		var revStr string
		if totalRevenue.Valid {
			revStr = fmt.Sprintf("%.2f", totalRevenue.Float64)
		} else {
			revStr = "NULL"
		}
		log.Printf("DB Row: Year=%d, Week=%d, Revenue=%s", year, week, revStr)
	}
	log.Println("Finished dumping trans_weekly_rev table contents.")
}
