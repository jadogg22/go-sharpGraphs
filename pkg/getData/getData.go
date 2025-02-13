package getdata

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/database"
	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"
	"github.com/joho/godotenv"

	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/integratedauth/krb5"
)

var (
	SQL_USER     string
	SQL_PASSWORD string
	SQL_SERVER   string
	SQL_DB       string
	URL          string

	conn    *sql.DB                           // db connection
	limiter = time.NewTicker(1 * time.Second) // db rate limiter
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			err := godotenv.Load(filepath.Join(dir, ".env"))
			if err != nil {
				fmt.Println("Error loading .env file")
			}
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			fmt.Println("Error finding .env file")
			break
		}
		dir = parent
	}

	SQL_USER = os.Getenv("SQL_USER")
	SQL_PASSWORD = os.Getenv("SQL_PASSWORD")
	SQL_SERVER = os.Getenv("SQL_SERVER")
	SQL_DB = os.Getenv("SQL_DB")

	URL = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&connection+timeout=30", SQL_USER, SQL_PASSWORD, SQL_SERVER, SQL_DB)

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file: " + err.Error())
	}
	defer file.Close()

	conn, err = sql.Open("mssql", URL)
	if err != nil {
		fmt.Println("Error creating connection pool: " + err.Error())
		return
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connected to the database")
}

func getTransportationTractorRevenue(conn *sql.DB) {
	// first query the sql server to get the data from mcloud

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//coupleDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	query := MakeQuery(yesterday)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return
	}
	defer rows.Close()

	var dbData []*models.TractorRevenue
	for rows.Next() {
		var d models.TractorRevenue
		var moveidstr string
		var orderidstr string

		err := rows.Scan(
			&moveidstr,
			&d.MoveDistance,
			&d.Loaded,
			&orderidstr,
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
			fmt.Println("Error scanning row: " + err.Error())
			return
		}

		moveid, err := strconv.Atoi(strings.TrimSpace(moveidstr))
		if err != nil {
			fmt.Println("Error converting moveid to int: " + err.Error())
			return
		}

		orderid, err := strconv.Atoi(strings.TrimSpace(orderidstr))
		if err != nil {
			fmt.Println("Error converting orderid to int: " + err.Error())
			return
		}

		d.MoveID = moveid
		d.OrderID = orderid

		dbData = append(dbData, &d)
	}

	for _, d := range dbData {
		var user string
		if d.UserName.Valid {
			user = d.UserName.String
		} else {
			user = "NULL"
		}
		fmt.Printf("Move_ID: %d, username: %s\n", d.MoveID, user)
	}
}

func TransportationRevenue(startDate, endDate time.Time) (float64, error) {

	// format time objects to be "2024-08-11"
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")
	query := fmt.Sprintf("SELECT SUM(total_charge) from orders where company_id = 'tms' and bol_recv_date between %s and %s", startDateStr, endDateStr)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return 0, err
	}

	defer rows.Close()

	var totalRevenue float64
	for rows.Next() {
		err := rows.Scan(&totalRevenue)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return 0, err
		}
	}

	return totalRevenue, nil
}

func LogisiticsRevenue(startDate, endDate time.Time) (float64, error) {

	// format time objects to be "2024-08-11"
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")
	query := fmt.Sprintf("SELECT SUM(total_charge) from orders where company_id = 'tms2' and bol_recv_date between %s and %s", startDateStr, endDateStr)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return 0, err
	}

	defer rows.Close()

	var totalRevenue float64
	for rows.Next() {
		err := rows.Scan(&totalRevenue)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return 0, err
		}
	}

	return totalRevenue, nil
}

func UpdateDateRangeAmounts(dateRanges []*helpers.DateRange) error {

	for _, dateRange := range dateRanges {
		startDateStr := dateRange.StartDate.Format("2006-01-02 00:00:00")
		EndDateStr := dateRange.EndDate.Format("2006-01-02 00:00:00")

		query := fmt.Sprintf("SELECT SUM(total_charge) from orders where company_id = 'tms' and bol_recv_date between '%s' and '%s'", startDateStr, EndDateStr)

		var amount sql.NullFloat64

		err := conn.QueryRow(query).Scan(&amount)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("No rows returned from query")
				dateRange.Amount = 0
			} else {
				fmt.Println("Error querying database: " + err.Error())
				return err
			}
		}

		if amount.Valid {
			dateRange.Amount = amount.Float64
		} else {
			dateRange.Amount = 0
		}
	}
	return nil
}

func UpdateTransRevData(data []models.WeeklyRevenue) {

	latestDataDate, err := helpers.FindLatestDateFromRevenueData(data)
	if err != nil {
		err := errors.New("Error finding latest date from revenue data" + err.Error())
		fmt.Println(err)
	}

	var dateRanges []*helpers.DateRange
	dateRanges = helpers.GenerateDateRanges(latestDataDate)
	if len(dateRanges) < 1 {
		fmt.Println("No date ranges found")
	}

	// Update the date range amounts

	if err := UpdateDateRangeAmounts(dateRanges); err != nil {
		err := errors.New("Error updating date range amounts " + err.Error())
		fmt.Println(err)
	}

	if len(dateRanges) < 3 {
		// take the first weeks up the the last 3 weeks
		newDateRanges := dateRanges[:3]
		//update my database with the new data
		database.UpdateMyDatabase(newDateRanges)
	}

	helpers.UpdateWeeklyRevenue(data, dateRanges)
}

type revenueData struct {
	TotalRevenue float64
	year         int
	week         int
}

func GetLogisticsMTDData(startDate, endDate time.Time) []models.LogisticsMTDStats {

	var dispacherNames = map[string]string{
		"cami":     "Cami Hansen",
		"jerrami":  "Jerrami Marotz",
		"joylynn":  "Joy Lynn",
		"lenora":   "Lenora Smith",
		"liz":      "Liz Swenson",
		"mijken":   "Mijken Cassidy",
		"riki":     "Riki Marotz",
		"samswens": "Sam Swenson",
	}

	query := getLogisticsMTDQuery(startDate, endDate) // Helperfunction to get the long query string
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return nil
	}

	defer rows.Close()

	// Rows returns
	//dispatcher_user_id	revenue	override_pay_amt	truck_hire	total_stops	total_servicefail_count	orders_with_service_fail	total_orders
	//cami      	313240.22	251552.70	256503.35	307	63	51	135	69430.00
	var data []models.LogisticsMTDStats
	var myDispatcher sql.NullString
	var dispatcher string
	var revenue, overridePayAmt, truckHire, miles float64
	var totalStops, totalServiceFailCount, ordersWithServiceFail, totalOrders int

	for rows.Next() {
		err := rows.Scan(&myDispatcher, &revenue, &overridePayAmt,
			&truckHire, &totalStops, &totalServiceFailCount,
			&ordersWithServiceFail, &totalOrders, &miles)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return nil
		}

		if !myDispatcher.Valid {
			dispatcher = "Unknown"
		} else {
			dispatcher = myDispatcher.String
		}
		// remove all spaces and make the dispatcher name lowercase
		dispatcher = strings.ReplaceAll(strings.ToLower(dispatcher), " ", "")

		if name, exists := dispacherNames[dispatcher]; exists {
			dispatcher = name
		}

		myData := models.NewLogisticsMTDStats(dispatcher, truckHire, revenue, miles, totalStops, totalOrders, ordersWithServiceFail, totalServiceFailCount)

		data = append(data, *myData)
	}

	rows.Close()

	return data
}

func GetTransportationDailyOps(startDate, endDate time.Time) ([]*models.DailyOpsData, error) {
	// format start and endDates to be "2024-08-11"
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	dispacherNames := map[string]string{
		"tracy":    "Tracy",
		"sheridan": "Sheridan",
		"rochelle": "Rochelle",
		"patrick":  "Patrick",
		"lindsay":  "Lindsay",
		"katrina":  "Katrina",
		"kenjr":    "Ken Jr",
		"todd":     "Todd",
		"amber":    "Amber",
	}

	err := conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
		return nil, err
	}

	myData := make([]*models.DailyOpsData, 0)

	query := MakeTransportationDailyOpsQuery(startDateStr, endDateStr)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return nil, fmt.Errorf("Error querying database: %v", err)
	}

	defer rows.Close()

	var dispacher_user_id string
	var total_empty_distance, total_loaded_distance sql.NullFloat64
	var total_stops, total_servicefail_count, orders_with_service_fail, total_orders, total_unique_trucks int
	//unpack the rows from the query into the data struct

	for rows.Next() {
		if rows.Err() != nil {
			fmt.Println("Error scanning row: " + rows.Err().Error())
			continue
		}
		// scan the row into the variables
		rows.Scan(&dispacher_user_id, &total_stops, &total_servicefail_count, &orders_with_service_fail, &total_orders, &total_empty_distance, &total_loaded_distance, &total_unique_trucks)

		// sanitize the dispacher_user_id
		dispacher_user_id = strings.ToLower(strings.ReplaceAll(dispacher_user_id, " ", ""))

		// look for name in the map
		if name, exists := dispacherNames[dispacher_user_id]; !exists {
			continue
		} else {
			dispacher_user_id = name
		}

		// create a new data struct
		myDispacherStats := models.NewDailyOpsDataFromDB(dispacher_user_id, total_loaded_distance, total_empty_distance, total_stops, total_servicefail_count, orders_with_service_fail, total_orders, total_unique_trucks)

		myData = append(myData, myDispacherStats)
	}

	if len(myData) < 1 {
		fmt.Println("No data returned from the query")
	}

	rows.Close()
	return myData, nil
}

func GetVacationFromDB(companyId string) ([]models.VacationHours, error) {
	if companyId != "tms" && companyId != "tms2" && companyId != "tms3" && companyId != "drivers" {
		return nil, fmt.Errorf("Invalid companyID")
	}

	// helper function to grab the sql query string
	query := GetVacationHoursByCompanyQuery(companyId)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	var data []models.VacationHours
	var employeeID, employeeName string
	var vacationHoursDue, vacationHoursRate sql.NullFloat64
	for rows.Next() {
		err := rows.Scan(&employeeID, &employeeName, &vacationHoursRate, &vacationHoursDue)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return nil, err
		}

		// clean up the data
		employeeID = strings.TrimSpace(employeeID)
		employeeName = strings.TrimSpace(employeeName)

		var vacationHoursDueStr string
		var vacationHoursRateStr string
		var amount float64

		if vacationHoursDue.Valid && vacationHoursRate.Valid {
			// some rates are salary based for the week
			if vacationHoursRate.Float64 > 100.0 {
				amount = (vacationHoursRate.Float64 / 80.0) * vacationHoursDue.Float64
			} else {
				amount = vacationHoursRate.Float64 * vacationHoursDue.Float64
			}
		}

		if !vacationHoursDue.Valid {
			vacationHoursDueStr = "Not entered"
		} else {
			vacationHoursDueStr = fmt.Sprintf("%.2f", vacationHoursDue.Float64)
		}

		if !vacationHoursRate.Valid {
			vacationHoursRateStr = "Not entered"
		} else {
			vacationHoursRateStr = fmt.Sprintf("%.2f", vacationHoursRate.Float64)
		}

		myData := models.VacationHours{
			EmployeeID:        employeeID,
			EmployeeName:      employeeName,
			VacationHoursDue:  vacationHoursDueStr,
			VacationHoursRate: vacationHoursRateStr,
			AmountDue:         fmt.Sprintf("%.2f", amount),
		}

		data = append(data, myData)
	}
	if len(data) < 1 {
		err := fmt.Errorf("Server error, No data returned from the query")
		return nil, err
	}
	return data, nil
}

func GetCodedRevenueData(when string) ([]models.CodedData, error) {
	// get the dates ranges
	fmt.Println("Getting the start and end dates")
	var startDate, endDate time.Time
	switch when {
	case "week":
		startDate, endDate = helpers.GetWeekStartAndEndDates()
	case "month":
		startDate, endDate = helpers.GetMonthStartAndEndDates()
	case "quarter":
		startDate, endDate = helpers.GetQuarterStartAndEndDates()
	default:
		return nil, fmt.Errorf("Invalid time period")
	}

	if startDate.IsZero() || endDate.IsZero() {
		return nil, fmt.Errorf("Error getting start and end dates")
	}

	// get the data from the database
	query := MakeCodedRevenueQuery(startDate, endDate)
	fmt.Println("Query: ", query)
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	myData := make(map[string]models.CodedData)

	fmt.Println("Getting data from the database")
	// scan the rows into the data struct
	for rows.Next() {
		var name string
		var revenue float64

		err := rows.Scan(&name, &revenue)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return nil, err
		}

		info, ok := myData[name]
		if !ok {
			myData[name] = models.CodedData{Name: name, Revenue: revenue, Count: 1}
		} else {
			info.Revenue += revenue
			info.Count++
			myData[name] = info
		}
	}

	sortedData := helpers.SortData(myData)          // changes the map to a slice of structs and sorts it uses sort package
	combinedData := helpers.CombineData(sortedData) // combines the small data into a single struct
	return combinedData, nil
}

// Stacked miles data endpont

func GetStackedMilesData(when string) ([]models.StackedMilesData, error) {
	// get the dates ranges
	fmt.Println("Getting the start and end dates")
	var startDate, endDate time.Time
	switch when {
	case "week_to_date":
		endDate = time.Now()
		startDate = endDate.AddDate(0, 0, -10)
		fmt.Println("Start Date: ", startDate)
		fmt.Println("End Date: ", endDate)
	case "month_to_date":
		endDate = time.Now()
		// 6 weeks worth
		startDate = endDate.AddDate(0, -1, -14)
	case "quarter":
		// IDK how to split this one up exactly yet
		startDate, endDate = helpers.GetQuarterStartAndEndDates()
	default:
		return nil, fmt.Errorf("Invalid time period")
	}

	if startDate.IsZero() || endDate.IsZero() {
		return nil, fmt.Errorf("Error getting start and end dates")
	}

	// get the data from the database
	query := MakeStackedMilesQuery(startDate, endDate)
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return nil, err
	}

	defer rows.Close()

	myData := make([]models.StackedMilesData, 0)

	// scan the rows into the data struct
	for rows.Next() {
		var id, date string
		var emptyMiles, loadedMiles float64

		err := rows.Scan(&id, &date, &loadedMiles, &emptyMiles)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return nil, err
		}

		myData = append(myData, models.StackedMilesData{ID: id, Date: date, EmptyMiles: emptyMiles, LoadedMiles: loadedMiles})
	}

	aggregateData := helpers.CombineStackedMilesData(when, myData)

	return aggregateData, nil
}

func GetSportsmanFromDB(date1, date2 string) ([]models.SportsmanData, error) {

	// helper function to get the query string in msQuerys.go
	myQuery := MakeSportsmansQuery(date1, date2)

	dbData := make([]models.SportsmanData, 0)

	rows, err := conn.Query(myQuery)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return dbData, err
	}

	defer rows.Close()

	var StartCity, StartState, StartZip, EndCity, EndState, EndZip sql.NullString

	var OrderID, OrderedDate, DelDate, BillDate, Consignee, Miles, BolNumber, Commodity, Weight sql.NullString
	var MovementSequence, PalletsDropped, PalletsPickedUp, TotalPallets sql.NullInt64
	var FreightCharge, fuel_surcharge, extra_drops, extra_pickup, other_charge, OtherChargeTotal, TotalCharge, per_pallet_fuel, per_pallet_freight sql.NullFloat64

	var TrailerNumber sql.NullString

	totalPalletsPickedUp := make(map[string]int64)

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Error fetching columns: " + err.Error())
		return dbData, err
	}

	fmt.Println("Columns: ", cols)

	for rows.Next() {
		err := rows.Scan(&OrderID, &OrderedDate,
			&DelDate, &BillDate,
			&EndCity, &EndState,
			&EndZip, &Consignee,
			&Miles, &BolNumber,
			&Commodity, &Weight,
			&MovementSequence,
			&PalletsDropped, &PalletsPickedUp,
			&FreightCharge, &OtherChargeTotal,
			&TotalCharge, &fuel_surcharge,
			&extra_drops, &extra_pickup,
			&other_charge, &per_pallet_fuel,
			&per_pallet_freight, &TrailerNumber)

		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			return dbData, err
		}

		if MovementSequence.Valid && MovementSequence.Int64 == 1 {
			TotalPallets = PalletsPickedUp
			StartCity = EndCity
			StartState = EndState
			StartZip = EndZip
		}

		// if pallets is not null add to the total pallets for the order
		if PalletsPickedUp.Valid && OrderID.Valid {
			myOrder := OrderID.String
			//strip the whitespace
			myOrder = strings.TrimSpace(myOrder)
			myPallets := PalletsPickedUp.Int64
			totalPalletsPickedUp[myOrder] += myPallets
		}

		//

		myData := models.NewSportsmanData(OrderID, OrderedDate, DelDate, BillDate, StartCity, StartState, StartZip, EndCity, EndState, EndZip, Consignee, Miles, BolNumber, Commodity, Weight, MovementSequence, PalletsDropped, PalletsPickedUp, TotalPallets, FreightCharge, fuel_surcharge, extra_drops, extra_pickup, other_charge, OtherChargeTotal, TotalCharge, per_pallet_fuel, per_pallet_freight, TrailerNumber)

		dbData = append(dbData, *myData)
	}

	// Second pass to update totals and per-pallet charges
	// Probally should have done a pointer to sructs buttt idc
	for i, d := range dbData {

		// Retrieve the total pickups for this order
		totalPickups := totalPalletsPickedUp[d.Order_id]

		// Avoid division by zero
		var pickupPercentage float64

		if totalPickups > 0 {
			pickupPercentage = float64(d.Pallets_Droped) / float64(totalPickups)
		} else {
			pickupPercentage = 0
		}

		// Calculate the per-pallet charges
		freightPerPallet := d.Freight_Charges * pickupPercentage
		fuelPerPallet := d.Fuel_Surcharge * pickupPercentage

		// Update the struct with the calculated values
		d.Freight_per_pallet = freightPerPallet
		d.Fuel_per_pallet = fuelPerPallet
		d.Total_Pallets = totalPickups

		dbData[i] = d
	}

	return dbData, nil
}

type FreightRecord struct {
	FreightCharge    float64   `db:"freight_charge"`     // Freight charge amount
	OtherChargeTotal float64   `db:"otherchargetotal"`   // Other charges
	TotalCharge      float64   `db:"total_charge"`       // Total charge
	Xferred2Billing  *float64  `db:"xferred2billing"`    // Transfer to billing (nullable)
	ID               string    `db:"id"`                 // ID of the order
	CityName         string    `db:"city_name"`          // City name
	State            string    `db:"state"`              // State code
	ActualArrival    time.Time `db:"actual_arrival"`     // Actual arrival timestamp
	SchedArriveEarly time.Time `db:"sched_arrive_early"` // Scheduled early arrival timestamp
	RevenueCodeID    string    `db:"revenue_code_id"`    // Revenue code ID
	BillDistance     float64   `db:"bill_distance"`      // Billable distance
}

func GetDashboardStats() (float64, float64) {
	date1 := "2024-12-18 00:00:00"
	date2 := "2024-12-18 23:59:59"

	query := DashboardQuery(date1, date2)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return 0, 0
	}

	// Initialize an empty slice to hold FreightRecords
	var freightRecords []FreightRecord
	var totalCharge float64
	var totalDistance float64

	// Loop through the rows returned from the query
	for rows.Next() {
		var rec FreightRecord
		var totalChargeStr, billDistanceStr, xferred2BillingStr string
		var actualArrivalStr, schedArriveEarlyStr string

		// Scan the row data into variables
		err := rows.Scan(
			&rec.ID,
			&totalChargeStr,
			&billDistanceStr,
			&xferred2BillingStr,
			&rec.CityName,
			&rec.State,
			&actualArrivalStr,
			&schedArriveEarlyStr,
			&rec.RevenueCodeID,
			&totalChargeStr,
		)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			continue
		}

		// Convert string fields to the appropriate types
		rec.TotalCharge, err = strconv.ParseFloat(totalChargeStr, 64)
		if err != nil {
			rec.TotalCharge = 0
		}

		rec.BillDistance, err = strconv.ParseFloat(billDistanceStr, 64)
		if err != nil {
			rec.BillDistance = 0
		}

		// If Xferred2Billing is not empty, parse it
		if xferred2BillingStr != "" {
			xferred2Billing, err := strconv.ParseFloat(xferred2BillingStr, 64)
			if err != nil {
				rec.Xferred2Billing = nil
			} else {
				rec.Xferred2Billing = &xferred2Billing
			}
		}

		// Convert the arrival dates into time
		rec.ActualArrival, err = time.Parse("2006-01-02 15:04:05", actualArrivalStr)
		if err != nil {
			rec.ActualArrival = time.Time{}
		}

		rec.SchedArriveEarly, err = time.Parse("2006-01-02 15:04:05", schedArriveEarlyStr)
		if err != nil {
			rec.SchedArriveEarly = time.Time{}
		}

		// Add the FreightRecord to the array
		freightRecords = append(freightRecords, rec)

		// Accumulate totals (if needed)
		totalCharge += rec.FreightCharge
		totalDistance += rec.BillDistance
	}

	// Handle any errors that occurred while iterating the rows
	if err = rows.Err(); err != nil {
		fmt.Println("Error iterating over rows: " + err.Error())
		return 0, 0
	}

	// Return the total charge and total distance
	return totalCharge, totalDistance
}

func GetDriverManagerData() ([]models.Driver, error) {
	// Get the query
	query := GetDriverManagerQuery()

	// Execute the query
	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	// Map to store drivers and their summaries
	driversMap := make(map[string]*models.Driver)

	// get the time
	now := time.Now()
	_, week := now.ISOWeek()

	// Iterate through rows
	for rows.Next() {
		var driverID, fleetManager, monthName string
		var weekNumber, monthOrder int
		var totalMoveDistance float64

		// Scan the row
		if err := rows.Scan(&driverID, &fleetManager, &weekNumber, &monthName, &monthOrder, &totalMoveDistance); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Check if driver exists in the map
		driver, exists := driversMap[driverID]
		if !exists {
			// Create a new driver if not already in the map
			driver = &models.Driver{
				DriverID:     driverID,
				FleetManager: fleetManager,
				Miles:        make([]int, (week - 1)), // Initialize an empty array for miles
			}
			driversMap[driverID] = driver
		}

		// Expand the array size if needed
		for len(driver.Miles) <= week {
			driver.Miles = append(driver.Miles, 0) // Fill missing weeks with 0 miles
		}

		// Store the miles in the correct position
		driver.Miles[weekNumber-1] = int(totalMoveDistance)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Convert map to a slice
	drivers := make([]models.Driver, 0, len(driversMap))
	for _, driver := range driversMap {
		drivers = append(drivers, *driver)
	}

	return drivers, nil
}
