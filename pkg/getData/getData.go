package getdata

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

	SQL_USER = os.Getenv("SQL_USER")
	SQL_PASSWORD = os.Getenv("SQL_PASSWORD")
	SQL_SERVER = os.Getenv("SQL_SERVER")
	SQL_DB = os.Getenv("SQL_DB")

	URL = fmt.Sprintf("sqlserver://%s:%s@%s?database=%s&connection+timeout=30", SQL_USER, SQL_PASSWORD, SQL_SERVER, SQL_DB)

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer file.Close()

	// Set log output to the file
	log.SetOutput(file)

}

func getTransportationTractorRevenue(conn *sql.DB) {
	// first query the sql server to get the data from mcloud

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//coupleDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	query := MakeQuery(yesterday)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
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
			log.Println("Error scanning row: " + err.Error())
			return
		}

		moveid, err := strconv.Atoi(strings.TrimSpace(moveidstr))
		if err != nil {
			fmt.Println("Error converting moveid to int: " + err.Error())
			log.Println("Error converting moveid to int: " + err.Error())
			return
		}

		orderid, err := strconv.Atoi(strings.TrimSpace(orderidstr))
		if err != nil {
			fmt.Println("Error converting orderid to int: " + err.Error())
			log.Println("Error converting orderid to int: " + err.Error())
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

/*
// the newest iteration of getting orders data - now in real time!
func GetTransportationOrdersData(conn *sql.DB, startDate, endDate time.Time) ([]models.OrderDetails, error) {

	// make start and end dates into strings
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	// get query string from helper function
	query := MakeTransportationOrdersQuery(startDateStr, endDateStr)

	// query the database
	rows, err := conn.Query(query)
	if err != nil {
		err := fmt.Errorf("Error querying database: %v", err)
		return nil, err
	}

	defer rows.Close()

	var data []models.OrderDetails
	for rows.Next() {
		var d models.OrderDetails
		var orderIDStr string
		var OperationsUser sql.NullString
		var RevenueCodeID string
		var FreightCharge float64
		var BillMiles float64
		var BillDate sql.NullTime
		var OrderTrailerType string
		var OriginValue string
		var DestinationValue string
		var CustomerID sql.NullString
		var CustomerName sql.NullString
		var CustomerCategory sql.NullString
		var CategoryDescr sql.NullString
		var MovementID string
		var Loaded string
		var MoveDistance float64
		var Brokerage string
		var OriginCity string
		var OriginState string
		var DestCity string
		var DestState string
		var ReportDate sql.NullTime
		var ActualDate sql.NullTime
		var EmptyMiles sql.NullFloat64
		var LoadedMiles sql.NullFloat64
		var TotalMiles sql.NullFloat64
		var TotalRevenue float64
		var WeekValue sql.NullString
		var MonthValue sql.NullString
		var QuarterValue sql.NullString
		var DetailID string

		err := rows.Scan(
			&orderIDStr,
			&OperationsUser,
			&RevenueCodeID,
			&FreightCharge,
			&BillMiles,
			&BillDate,
			&OrderTrailerType,
			&OriginValue,
			&DestinationValue,
			&CustomerID,
			&CustomerName,
			&CustomerCategory,
			&CategoryDescr,
			&MovementID,
			&Loaded,
			&MoveDistance,
			&Brokerage,
			&OriginCity,
			&OriginState,
			&DestCity,
			&DestState,
			&ReportDate,
			&ActualDate,
			&EmptyMiles,
			&LoadedMiles,
			&TotalMiles,
			&TotalRevenue,
			&WeekValue,
			&MonthValue,
			&QuarterValue,
			&DetailID)

		if err != nil {
			err := fmt.Errorf("Error scanning row: %v", err)
			return nil, err
		}

		d.OrderID = orderIDStr
		if OperationsUser.Valid {
			d.OperationsUser = OperationsUser.String
		} else {
			d.OperationsUser = "Loadmaster"
		}
		d.RevenueCodeID = RevenueCodeID
		d.FreightCharge = FreightCharge
		d.BillMiles = BillMiles

		if BillDate.Valid {
			d.BillDate = BillDate.Time
		}
		d.OrderTrailerType = OrderTrailerType
		d.OriginValue = OriginValue
		d.DestinationValue = DestinationValue
		if CustomerID.Valid {
			d.CustomerID = CustomerID.String
		} else {
			d.CustomerID = "Unknown"
		}
		if CustomerName.Valid {
			d.CustomerName = CustomerName.String
		} else {
			d.CustomerName = "Unknown"
		}
		if CustomerCategory.Valid {
			d.CustomerCategory = CustomerCategory.String
		} else {
			d.CustomerCategory = "Unknown"
		}
		if CategoryDescr.Valid {
			d.CategoryDescr = CategoryDescr.String
		} else {
			d.CategoryDescr = "Unknown"
		}
		d.MovementID = MovementID
		d.Loaded = Loaded
		d.MoveDistance = MoveDistance
		d.Brokerage = Brokerage
		d.OriginCity = OriginCity
		d.OriginState = OriginState
		d.DestCity = DestCity
		d.DestState = DestState
		if ReportDate.Valid {
			d.ReportDate = ReportDate.Time
		}
		if ActualDate.Valid {
			d.ActualDate = ActualDate.Time
		}
		if EmptyMiles.Valid {
			d.EmptyMiles = EmptyMiles.Float64
		}
		if LoadedMiles.Valid {
			d.LoadedMiles = LoadedMiles.Float64
		}
		if TotalMiles.Valid {
			d.TotalMiles = TotalMiles.Float64
		}
		d.TotalRevenue = TotalRevenue
		if WeekValue.Valid {
			d.WeekValue = WeekValue.String
		}
		if MonthValue.Valid {
			d.MonthValue = MonthValue.String
		}
		if QuarterValue.Valid {
			d.QuarterValue = QuarterValue.String
		}
		d.DetailID = DetailID

		data = append(data, d)

	}

	return data, nil
}
*/

// With access to the db, we can now grab the data from the production database and dont need
// to replicate the data in our own database. This will save time and space and provide greater accuracy
func NoDBDailyOpsCacheGrab() {
	// make a connection to the database
	conn, err := sql.Open("mssql", URL)
	if err != nil {
		fmt.Println("Error creating connection pool: " + err.Error())
		log.Println("Error creating connection pool: " + err.Error())
		return
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
		log.Println("Error pinging database: " + err.Error())
		return
	}

	// get all the data from the database
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

	conn, err := sql.Open("mssql", URL)
	if err != nil {
		fmt.Println("Error creating connection pool: " + err.Error())
		return nil
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
		return nil
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
	conn.Close()

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
	}

	conn, err := sql.Open("mssql", URL)
	if err != nil {
		fmt.Println("Error creating connection pool: " + err.Error())
		return nil, err
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
		return nil, err
	}

	myData := make([]*models.DailyOpsData, 0)

	query := MakeTransportationDailyOpsQuery(startDateStr, endDateStr)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
		return nil, fmt.Errorf("Error querying database: %v", err)
	}

	defer rows.Close()

	var dispacher_user_id string
	var total_bill_distance, total_move_distance sql.NullFloat64
	var total_stops, total_servicefail_count, orders_with_service_fail, total_orders, total_unique_trucks int
	//unpack the rows from the query into the data struct

	for rows.Next() {
		if rows.Err() != nil {
			fmt.Println("Error scanning row: " + rows.Err().Error())
			log.Println("Error scanning row: " + rows.Err().Error())
			continue
		}
		// scan the row into the variables
		rows.Scan(&dispacher_user_id, &total_stops, &total_servicefail_count, &orders_with_service_fail, &total_orders, &total_bill_distance, &total_move_distance, &total_unique_trucks)

		// sanitize the dispacher_user_id
		dispacher_user_id = strings.ToLower(strings.ReplaceAll(dispacher_user_id, " ", ""))

		// look for name in the map
		if name, exists := dispacherNames[dispacher_user_id]; !exists {
			continue
		} else {
			dispacher_user_id = name
		}

		// create a new data struct
		myDispacherStats := models.NewDailyOpsDataFromDB(dispacher_user_id, total_bill_distance, total_move_distance, total_stops, total_servicefail_count, orders_with_service_fail, total_orders, total_unique_trucks)

		myData = append(myData, myDispacherStats)
	}

	if len(myData) < 1 {
		fmt.Println("No data returned from the query")
		log.Println("No data returned from the query")
	}

	rows.Close()
	conn.Close()
	return myData, nil
}

func GetVacationFromDB(companyId string) ([]models.VacationHours, error) {
	if companyId != "tms" && companyId != "tms2" && companyId != "tms3" {
		return nil, fmt.Errorf("Invalid companyID")
	}
	if companyId == "drivers" || companyId == "all" {
		error := fmt.Errorf("Server error, Unimplmented companyID")
		return nil, error
	}

	// helper function to grab the sql query string
	query := GetVacationHoursByCompanyQuery(companyId)
	conn, err := sql.Open("mssql", URL)
	if err != nil {
		fmt.Println("Error creating connection pool: " + err.Error())
		return nil, err
	}

	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
		return nil, err
	}

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
			if vacationHoursRate.Float64 > 45.0 {
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
