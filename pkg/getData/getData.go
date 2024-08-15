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

	"github.com/jadogg22/go-sharpGraphs/pkg/database"
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

func RunUpdater() {

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

	// helper functions to grab data from the database
	getTransportationTractorRevenue(conn)
	addTransportationBadStops(conn)

	conn.Close()

	return
}

func TestConnection() (string, error) {
	conn, err := sql.Open("mssql", URL)
	if err != nil {
		return "Error creating connection pool: " + err.Error(), err
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		return "Error pinging database: " + err.Error(), err
	}

	queary := "select top 10 customer_id from orders"
	rows, err := conn.Query(queary)
	if err != nil {
		return "Error querying database: " + err.Error(), err
	}
	defer rows.Close()

	rows.Next()
	var customerID string
	err = rows.Scan(&customerID)
	if err != nil {
		return "Error scanning row: " + err.Error(), err
	}

	conn.Close()
	return fmt.Sprintf("We did it boys: %s", customerID), nil
}

func getTransportationTractorRevenue(conn *sql.DB) {
	// first query the sql server to get the data from mcloud

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//coupleDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	query := fmt.Sprintf(`SELECT 
    movement.id AS move_id,
    movement.move_distance AS move_distance,
    movement.loaded AS loaded,
    orders.id AS order_id,
    orders.total_charge AS charges,
    orders.bill_distance AS bill_distance,
    orders.freight_charge AS freight_charge,
    origin.city_name AS origin_city,
    origin.state AS origin_state,
    continuity.equipment_id AS equip_id,
    continuity.actual_arrival AS actual_arrival,
    continuity.dest_actualarrival AS del_date,
    continuity.equipment_id AS tractor,
    continuity.equipment_type_id,
    tractor.dispatcher AS dispatcher,
    tractor.fleet_id AS fleet_id,
    fleet.description AS fleet_description,
    users.name AS user_name,
    COUNT(servicefail.id) AS servicefail_count,
    CASE WHEN COUNT(servicefail.id) > 0 THEN 1 ELSE 0 END AS has_servicefail,
	stop_count.stop_count AS stop_count
FROM 
    movement
    JOIN movement_order ON movement.id = movement_order.movement_id AND movement_order.company_id = 'TMS'
    JOIN orders ON orders.id = movement_order.order_id AND orders.company_id = 'TMS'
    JOIN stop origin ON origin.movement_id = movement.id AND origin.movement_sequence = 1 AND origin.company_id = 'TMS'
    JOIN stop dest ON dest.id = movement.dest_stop_id AND dest.company_id = 'TMS'
    JOIN continuity ON movement.id = continuity.movement_id AND continuity.equipment_type_id = 'T' AND continuity.company_id = 'TMS'
    JOIN tractor ON tractor.id = continuity.equipment_id AND tractor.company_id = 'TMS'
    LEFT JOIN fleet ON fleet.id = tractor.fleet_id AND fleet.company_id = 'TMS'
    LEFT JOIN users ON users.id = tractor.dispatcher AND users.company_id = 'TMS'
    LEFT JOIN servicefail ON servicefail.order_id = orders.id AND servicefail.company_id = 'TMS'
	LEFT JOIN (
        SELECT order_id, COUNT(*) AS stop_count
        FROM stop
        WHERE company_id = 'TMS'
        GROUP BY order_id
    ) AS stop_count ON stop_count.order_id = orders.id

WHERE 
    movement.company_id = 'TMS' 
    AND continuity.dest_actualarrival >= {ts '%s 00:00:00'}
    AND continuity.dest_actualarrival <= {ts '%s 23:59:59'}
    AND movement.status <> 'V'
GROUP BY 
    movement.id, 
    movement.move_distance, 
    movement.loaded, 
    orders.id,
    orders.total_charge, 
    orders.bill_distance, 
    orders.freight_charge, 
    origin.city_name, 
    origin.state,
    continuity.equipment_id,
    continuity.actual_arrival, 
    continuity.dest_actualarrival, 
    continuity.equipment_id, 
    continuity.equipment_type_id,
    tractor.dispatcher, 
    tractor.fleet_id, 
    fleet.description, 
    users.name, 
	stop_count.stop_count
ORDER BY 
    dispatcher, 
    tractor, 
    continuity.dest_actualarrival;`, yesterday, yesterday)

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

	// add dbData to the database
	err = database.AddTransprotationTractorRevenue(dbData)
	if err != nil {
		fmt.Println("Error adding data to database: " + err.Error())
		log.Println("Error adding data to database: " + err.Error())
		return
	}
	log.Printf("Added %d rows to the database\n", len(dbData))
}

func addTransportationBadStops(conn *sql.DB) {
	// first query the sql server to get the data from mcloud

	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//coupleDaysAgo := time.Now().AddDate(0, 0, -3).Format("2006-01-02")

	query := fmt.Sprintf(`select stop.id, stop.order_id, stop.movement_id, stop.actual_arrival,
	stop.sched_arrive_early, stop.sched_arrive_late, stop.movement_sequence,
	movement.equipment_group_id equipment_group_id, movement.dispatcher_user_id dispatcher_user_id,
	equipment_item.equipment_id equipment_id, equipment_item.equipment_type_id equipment_type_id,
	driver.fleet_manager fleet_manager, driver.id driver_id, servicefail.stop_id stop_id, servicefail.minutes_late minutes_late,
	servicefail.appt_required appt_required, servicefail.stop_type stop_type, servicefail.entered_user_id entered_user_id,
	servicefail.entered_date entered_date, servicefail.edi_standard_code edi_standard_code, servicefail.dsp_comment dsp_comment,
	servicefail.fault_of_carrier_or_driver sf_fault_of_carrier_or_driver, orders.customer_id customer_id,
	orders.operations_user operations_user, orders.status order_status 
from stop left outer join servicefail on servicefail.stop_id = stop.id and servicefail.status != 'V' and servicefail.company_id = 'TMS'  
left outer join orders on orders.id = stop.order_id  and orders.company_id = 'TMS'  ,movement 
left outer join equipment_item on equipment_item.equipment_group_id = movement.equipment_group_id and equipment_item.equipment_type_id = 'D' and equipment_item.type_sequence = 0 and equipment_item.company_id = 'TMS'  
left outer join driver on driver.id = equipment_item.equipment_id and driver.company_id = 'TMS'  
where stop.company_id = 'TMS' and stop.sched_arrive_early between {ts '%s 00:00:00'} and {ts '%s 23:59:59'} and stop.stop_type in ('PU', 'SO') and movement.loaded = 'L' and movement.id = stop.movement_id and movement.company_id = 'TMS' and driver.fleet_manager != 'NULL' and stop_id != 'NULL'
order by driver.fleet_manager, servicefail.minutes_late`, yesterday, yesterday)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
		return
	}
	defer rows.Close()

	var dbData []*models.BadStop
	for rows.Next() {
		row := models.BadStop{}
		err := rows.Scan(
			&row.ID,
			&row.OrderID,
			&row.MovementID,
			&row.ActualArrival,
			&row.SchedArriveEarly,
			&row.SchedArriveLate,
			&row.MovementSequence,
			&row.EquipmentGroupID,
			&row.DispatcherUserID,
			&row.EquipmentID,
			&row.EquipmentTypeID,
			&row.FleetManager,
			&row.DriverID,
			&row.StopID,
			&row.MinutesLate,
			&row.ApptRequired,
			&row.StopType,
			&row.EnteredUserID,
			&row.EnteredDate,
			&row.EDIStandardCode,
			&row.DSPComment,
			&row.SFFaultOfCarrierOrDriver,
			&row.CustomerID,
			&row.OperationsUser,
			&row.OrderStatus,
		)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			log.Println("Error scanning row: " + err.Error())
			return
		}

		dbData = append(dbData, &row)
	}

	// add dbData to the database
	err = database.AddBadStops(dbData)
	if err != nil {
		fmt.Println("Error adding data to database: " + err.Error())
		log.Println("Error adding data to database: " + err.Error())
		return
	}

	log.Printf("Added %d rows to the database\n", len(dbData))
}

func getLogisticsData(conn *sql.DB, startDate time.Time, endDate time.Time) {
	// format the dates to be used in the query "2024-08-11 00:00:00"
	startDateStr := startDate.Format("2006-01-02 15:04:05")
	endDateStr := endDate.Format("2006-01-02 15:04:05")

	fmt.Println(startDateStr, endDateStr)

	query := `WITH stop_data AS (
    SELECT 
        stop.id,
        stop.order_id,
        stop.movement_id,
        stop.actual_arrival,
        stop.sched_arrive_early,
        stop.sched_arrive_late,
        stop.movement_sequence,
        movement.dispatcher_user_id AS dispatcher_user_id,

        driver.fleet_manager AS fleet_manager,
        driver.id AS driver_id,
        servicefail.stop_id AS servicefail_stop_id,
        servicefail.minutes_late AS minutes_late,
        servicefail.appt_required AS appt_required,
        servicefail.stop_type AS stop_type,
        servicefail.entered_user_id AS entered_user_id,
        servicefail.entered_date AS entered_date,
        servicefail.edi_standard_code AS edi_standard_code,
        servicefail.dsp_comment AS dsp_comment,
        servicefail.fault_of_carrier_or_driver AS sf_fault_of_carrier_or_driver,
        movement.override_payee_id AS override_payee_id,
        orders.customer_id AS customer_id,
        movement.override_pay_amt AS pay_amt,
        orders.id AS orders_id,
        orders.operations_user AS operations_user,
        orders.total_charge AS total_charge,
        orders.status AS order_status,
		orders.bill_distance as bill_distance
    FROM 
        stop
    LEFT OUTER JOIN servicefail ON servicefail.stop_id = stop.id AND servicefail.status != 'V' AND servicefail.company_id = 'TMS2'
    LEFT OUTER JOIN orders ON orders.id = stop.order_id AND orders.company_id = 'TMS2'
    LEFT OUTER JOIN movement ON movement.id = stop.movement_id AND movement.company_id = 'TMS2'
    LEFT OUTER JOIN equipment_item ON equipment_item.equipment_group_id = movement.equipment_group_id AND equipment_item.equipment_type_id = 'D' AND equipment_item.type_sequence = 0 AND equipment_item.company_id = 'TMS2'
    LEFT OUTER JOIN driver ON driver.id = equipment_item.equipment_id AND driver.company_id = 'TMS2'
    WHERE 
        stop.company_id = 'TMS2' 
        AND stop.sched_arrive_early >= {ts '2024-08-11 00:00:00'}
        AND stop.sched_arrive_early <= {ts '2024-08-12 23:59:59'}
        AND stop.stop_type IN ('PU', 'SO')
        AND movement.loaded = 'L'
),
driver_extra_data AS (
    SELECT
        driver_extra_pay.movement_id,
        SUM(driver_extra_pay.amount) AS total_amount,
        MAX(deduct_code.code_type) AS code_type 
    FROM
        driver_extra_pay
    JOIN deduct_code ON deduct_code.id = driver_extra_pay.deduct_code_id
    WHERE
        driver_extra_pay.company_id = 'TMS2'
        AND deduct_code.company_id = 'TMS2'
    GROUP BY
        driver_extra_pay.movement_id
)
SELECT
    s.*,
    COALESCE(d.total_amount, 0) AS additional_amount,
    (s.pay_amt + COALESCE(d.total_amount, 0)) AS truck_hire,
    d.code_type AS deduct_code
FROM
    stop_data s
LEFT JOIN driver_extra_data d ON s.movement_id = d.movement_id
ORDER BY
    s.dispatcher_user_id,
    s.minutes_late;`

	fmt.Println(query)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
		return
	}

	defer rows.Close()

}
