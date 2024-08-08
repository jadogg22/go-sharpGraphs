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

	DB *sql.DB
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

	DB, err = PG_Make_connection()
	if err != nil {
		fmt.Println("Error connecting to the database")
	}

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
func Add_DailyDriverData(dailyData models.DailyDriverData) error {
	// Add the dailyDriverData to the db
	_, err := DB.Exec(`INSERT INTO daily_driver_data (dispatcher, deadhead_percent, freight, fuel_surcharge, remain_chgs, revenue, total_rev_per_rev_miles, total_rev_per_total_miles, average_weekly_rev, average_weekly_rev_miles, average_rev_miles, revenue_miles, total_miles, trucks, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
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
		dbTable = "trans_year_rev"
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

func FindNewestMilesData(company string) (time.Time, error) {
	var query string
	if company != "transportation" && company != "logistics" {
		return time.Time{}, fmt.Errorf("this company doesn't exist")
	}

	query = fmt.Sprintf(`
		SELECT MAX(deliverydate)
		FROM %s
	`, company)

	var newestDate time.Time
	err := DB.QueryRow(query).Scan(&newestDate)
	if err != nil {
		return time.Time{}, err
	}
	return newestDate, nil
}

func AddOrderToDB(loadData *[]models.LoadData, company string) error {
	// Begin a transaction
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // This will be a no-op if the transaction is committed successfully

	if company != "transportation" && company != "logistics" {
		return fmt.Errorf("invalid company")
	}

	query := fmt.Sprintf(`
       INSERT INTO %s (
    revenuecode, ordernumber, ordertype, freight, fuelsurcharge, remainingcharges,
    totalrevenue, billmiles, loadedmiles, emptymiles, totalmiles, emptypercentage,
    revloadedmile, revtotalmile, deliverydate, origin, destination, customer,
    customercategory, operationsuser, billed, controllingparty, commodity,
    trailertype, originstate, destinationstate, week, month, quarter, brokered
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
    $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
) ON CONFLICT (ordernumber) DO NOTHING;`, company)

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

func AddTransprotationTractorRevenue(dbData []*models.TractorRevenue) error {
	// Begin a transaction

	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback() // This will be a no-op if the transaction is committed successfully

	// Prepare the INSERT statement
	stmt, err := tx.Prepare(`
		INSERT INTO Transportation_Tractor_Revenue (
			move_id, move_distance, loaded, order_id, charges, bill_distance, freight_charge, origin_city, origin_state, equip_id, actual_arrival, del_date, tractor, equipment_type_id, dispatcher, fleet_id, fleet_description, user_name, service_fail_count, has_service_fail, stop_count 
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)
		ON CONFLICT (move_id) DO UPDATE SET
			move_distance = EXCLUDED.move_distance,
			loaded = EXCLUDED.loaded,
			order_id = EXCLUDED.order_id,
			charges = EXCLUDED.charges,
			bill_distance = EXCLUDED.bill_distance,
			freight_charge = EXCLUDED.freight_charge,
			origin_city = EXCLUDED.origin_city,
			origin_state = EXCLUDED.origin_state,
			equip_id = EXCLUDED.equip_id,
			actual_arrival = EXCLUDED.actual_arrival,
			del_date = EXCLUDED.del_date,
			tractor = EXCLUDED.tractor,
			equipment_type_id = EXCLUDED.equipment_type_id,
			dispatcher = EXCLUDED.dispatcher,
			fleet_id = EXCLUDED.fleet_id,
			fleet_description = EXCLUDED.fleet_description,
			user_name = EXCLUDED.user_name,
			service_fail_count = EXCLUDED.service_fail_count,
			has_service_fail = EXCLUDED.has_service_fail,
			stop_count = EXCLUDED.stop_count;`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert each TractorRevenue object into the database
	for _, d := range dbData {
		_, err := stmt.Exec(
			&d.MoveID,
			&d.MoveDistance,
			&d.Loaded,
			&d.OrderID,
			&d.Charges,
			&d.BillDistance,
			&d.FreightCharge,
			&d.OriginCity,
			&d.OriginState,
			&d.EquipID,
			&d.ActualArrival,
			&d.DelDate,
			&d.Tractor,
			&d.EquipmentTypeID,
			&d.Dispatcher,
			&d.FleetID,
			&d.FleetDescription,
			&d.UserName,
			&d.ServiceFailCount,
			&d.HasServiceFail,
			&d.StopCount,
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
