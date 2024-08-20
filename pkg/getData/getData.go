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

	// format the dates to be used in the query "2024-08-11 00:00:00"
	startDateStr := startDate.Format("2006-01-02 15:04:05")
	endDateStr := endDate.Format("2006-01-02 15:04:05")

	fmt.Println(startDateStr, endDateStr)

	orderQuery := `SELECT 
    users.name AS dispatcher,
    movement.override_pay_amt AS truck_hire,
    orders.total_charge AS charges,
	movement.move_distance AS miles
FROM 
    movement
    LEFT OUTER JOIN payee 
        ON payee.id = movement.override_payee_id 
        AND payee.company_id = 'TMS2'
    LEFT OUTER JOIN users 
        ON users.id = movement.dispatcher_user_id 
        AND users.company_id = 'TMS2'
    INNER JOIN movement_order 
        ON movement.id = movement_order.movement_id 
        AND movement_order.company_id = 'TMS2'
    INNER JOIN orders 
        ON orders.id = movement_order.order_id 
        AND orders.company_id = 'TMS2'
    INNER JOIN stop 
        ON stop.id = orders.shipper_stop_id 
        AND stop.company_id = 'TMS2'
WHERE 
    movement.company_id = 'TMS2'
    AND movement.status <> 'V'
    AND movement.loaded = 'L'
    AND stop.actual_arrival BETWEEN {ts '2024-08-18 00:00:00'} AND {ts '2024-08-18 23:59:59'}
ORDER BY 
    dispatcher, 
    orders.revenue_code_id, 
    orders.id;`

	rows, err := conn.Query(orderQuery)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
		return nil
	}

	defer rows.Close()

	var OrdersData []models.LogisticsOrdersData
	for rows.Next() {
		var d models.LogisticsOrdersData
		err := rows.Scan(
			&d.Dispacher,
			&d.Truck_hire,
			&d.Charges,
			&d.Miles,
		)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			log.Println("Error scanning row: " + err.Error())
			return nil
		}

		OrdersData = append(OrdersData, d)
	}

	stopsQuery := `
WITH FaultyStops AS (
    SELECT DISTINCT
        stop.order_id,
        movement.dispatcher_user_id
    FROM
        stop
        LEFT JOIN servicefail 
            ON servicefail.stop_id = stop.id 
            AND servicefail.status != 'V' 
            AND servicefail.company_id = 'TMS2'
        LEFT JOIN movement 
            ON movement.id = stop.movement_id
            AND movement.company_id = 'TMS2'
    WHERE
        stop.company_id = 'TMS2'
        AND stop.sched_arrive_early BETWEEN {ts '2024-08-18 00:00:00'} AND {ts '2024-08-18 23:59:59'}
        AND stop.stop_type IN ('PU', 'SO')
        AND servicefail.minutes_late IS NOT NULL
),
OrderFaults AS (
    SELECT
        dispatcher_user_id,
        COUNT(DISTINCT order_id) AS order_faults
    FROM
        FaultyStops
    GROUP BY
        dispatcher_user_id
),
TotalOrders AS (
    SELECT
        movement.dispatcher_user_id,
        COUNT(DISTINCT stop.order_id) AS total_orders
    FROM
        stop
        LEFT JOIN movement 
            ON movement.id = stop.movement_id
            AND movement.company_id = 'TMS2'
    WHERE
        stop.company_id = 'TMS2'
        AND stop.sched_arrive_early BETWEEN {ts '2024-08-18 00:00:00'} AND {ts '2024-08-18 23:59:59'}
        AND stop.stop_type IN ('PU', 'SO')
    GROUP BY
        movement.dispatcher_user_id
)
SELECT
    users.name AS dispatcher,
    COUNT(DISTINCT stop.id) AS total_stops,
    COUNT(DISTINCT CASE WHEN servicefail.minutes_late IS NOT NULL THEN stop.id END) AS stop_faults,
    COALESCE(OrderFaults.order_faults, 0) AS order_faults,
    COALESCE(TotalOrders.total_orders, 0) AS total_orders
FROM
    stop
    LEFT JOIN servicefail 
        ON servicefail.stop_id = stop.id 
        AND servicefail.status != 'V'
        AND servicefail.company_id = 'TMS2'
    LEFT JOIN movement 
        ON movement.id = stop.movement_id 
        AND movement.company_id = 'TMS2'
    LEFT JOIN OrderFaults 
        ON OrderFaults.dispatcher_user_id = movement.dispatcher_user_id
    LEFT JOIN TotalOrders 
        ON TotalOrders.dispatcher_user_id = movement.dispatcher_user_id
    LEFT JOIN users 
        ON users.id = movement.dispatcher_user_id 
        AND users.company_id = 'TMS2'
WHERE
    stop.company_id = 'TMS2'
    AND stop.sched_arrive_early BETWEEN {ts '2024-08-18 00:00:00'} AND {ts '2024-08-18 23:59:59'}
    AND stop.stop_type IN ('PU', 'SO')
    AND movement.loaded = 'L'
GROUP BY
    users.name,
    OrderFaults.order_faults,
    TotalOrders.total_orders
ORDER BY
    dispatcher;`

	stopsRows, err := conn.Query(stopsQuery)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
		return nil
	}

	defer stopsRows.Close()

	var StopsData []models.LogisticsStopOrdersData
	for stopsRows.Next() {
		var d models.LogisticsStopOrdersData
		err := stopsRows.Scan(
			&d.Dispacher,
			&d.Total_stops,
			&d.Stop_faults,
			&d.Order_faults,
			&d.Total_orders,
		)
		if err != nil {
			fmt.Println("Error scanning row: " + err.Error())
			log.Println("Error scanning row: " + err.Error())
			return nil
		}

		StopsData = append(StopsData, d)
	}

	//agrigate data
	agregedData := helpers.AgregateLogisticMTDStats(OrdersData, StopsData)
	conn.Close()

	return agregedData
}

func getTransportationOrders(conn *sql.DB, startDate, endDate time.Time) {
	// format start and endDates to be "2024-08-11"
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	query := fmt.Sprintf(`

select 
	orders.id order_id, orders.operations_user operations_user, orders.revenue_code_id revenue_code_id,
	orders.freight_charge freight_charge, orders.bill_distance bill_miles, orders.bill_date bill_date,
	orders.ctrl_party_id controlling_party, orders.commodity_id commodity, orders.order_type_id order_type,
	orders.equipment_type_id order_trailer_type, origin.state origin_value, dest.state destination_value, customer.id customer_id,
	customer.name customer_name, customer.category customer_category, category.descr category_descr, movement.id movement_id,
	loaded, move_distance, movement.brokerage, trailer.trailer_type trailer_type, origin.city_name origin_city, origin.state origin_state,
	dest.city_name dest_city, dest.state dest_state, other_charge.amount oc_amount, charge_code.is_fuel_surcharge is_fuel_surcharge,
	dest.sched_arrive_early report_date, dest.actual_arrival actual_date, prorated_orderdist.empty_distance empty_miles,
	prorated_orderdist.loaded_distance loaded_miles, (prorated_orderdist.empty_distance+prorated_orderdist.loaded_distance) total_miles,
	orders.id record_count, orders.id fuel_surcharge, orders.id remaining_charges, orders.id total_revenue, orders.id empty_pct,
	orders.id rev_loaded_mile, orders.id rev_total_mile, orders.id billed, orders.id week_value, orders.id month_value,
	orders.id quarter_value, revenue_code_id detail_id 
from 
	orders left outer join customer on customer.id = orders.customer_id and customer.company_id = 'TMS'  
	left outer join category on category.id = customer.category and category.company_id = 'TMS'  
	left outer join movement_order on movement_order.order_id = orders.id and movement_order.company_id = 'TMS'  
	left outer join movement on movement.id = movement_order.movement_id and movement.company_id = 'TMS'  
	left outer join continuity trailercont on (movement.id = trailercont.movement_id)and(trailercont.equipment_type_id='L') and  trailercont.company_id = 'TMS'  
	left outer join trailer on trailer.id = trailercont.equipment_id and trailer.company_id = 'TMS'  left outer join other_charge on other_charge.order_id = orders.id and other_charge.company_id = 'TMS'  
	left outer join charge_code on charge_code.id = other_charge.charge_id  
	left outer join prorated_orderdist on prorated_orderdist.order_id = orders.id and prorated_orderdist.company_id = 'TMS'  ,stop origin ,stop dest 
where 
	orders.company_id = 'TMS' and orders.status <> 'Q' and orders.status <> 'V' and (orders.subject_order_status is null or orders.subject_order_status <> 'S') and loaded = 'L' 
	and ((dest.actual_arrival is not null and dest.actual_arrival >= {ts '2024-08-11 00:00:00'}) or dest.actual_arrival is null and dest.sched_arrive_early >= {ts '%s 00:00:00'}) and ((dest.actual_arrival is not null and dest.actual_arrival <= {ts '%s 23:59:59'}) 
	or dest.actual_arrival is null and dest.sched_arrive_early <= {ts '2024-08-15 23:59:59'}) and origin.id = orders.shipper_stop_id  and  origin.company_id = 'TMS' and dest.id = orders.consignee_stop_id  and  dest.company_id = 'TMS' order by revenue_code_id, order_id, movement_id`, startDateStr, endDateStr)

	rows, err := conn.Query(query)
	if err != nil {
		fmt.Println("Error querying database: " + err.Error())
		log.Println("Error querying database: " + err.Error())
		return
	}

	defer rows.Close()

}
