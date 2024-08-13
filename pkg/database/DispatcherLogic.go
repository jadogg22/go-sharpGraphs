package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"
	"strings"
	"time"
)

type DispatcherData struct {
	Dispatcher   string
	Trucks       int
	Miles        int
	WeekDeadhead float64
	MPTPD        float64
	RPTPD        float64
	MPTPDColor   string
	RPTPDColor   string
	DHColor      string
	Location     string
}

const (
	MPTD_OTR_Goal   = 564.0
	MPTD_Local_Goal = 183.0
	MPTD_Texas_Goal = 271.0

	RPTPD_OTR_Goal   = 1249.0
	RPTPD_Local_Goal = 499.0
	RPTPD_Texas_Goal = 605.0

	DH_OTR_Goal   = 9.0
	DH_Local_Goal = 44.0
	DH_Texas_Goal = 29.0
)

type Goals map[string]map[string]float64

var goals = Goals{
	"OTR": {
		"MPTD":  MPTD_OTR_Goal,
		"RPTPD": RPTPD_OTR_Goal,
		"DH":    DH_OTR_Goal,
	},
	"Local": {
		"MPTD":  MPTD_Local_Goal,
		"RPTPD": RPTPD_Local_Goal,
		"DH":    DH_Local_Goal,
	},
	"Texas": {
		"MPTD":  MPTD_Texas_Goal,
		"RPTPD": RPTPD_Texas_Goal,
		"DH":    DH_Texas_Goal,
	},
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
func GetDailyOpsData(company string) ([]models.DailyOpsData, error) {
	// find the date for the last sunday
	var myData []models.DailyOpsData

	// get the data from the database
	query := `
WITH filtered_movements AS (
    SELECT 
        dispatcher,
        COUNT(DISTINCT tractor) AS unique_trucks,
        SUM(move_distance) AS total_miles,
        SUM(CASE WHEN loaded = 'E' THEN move_distance ELSE 0 END) AS empty_miles,
        COUNT(DISTINCT order_id) AS total_orders
    FROM 
        Transportation_Tractor_Revenue
    WHERE 
        del_date >= (current_date - extract(dow FROM current_date)::integer)
    GROUP BY 
        dispatcher
),
order_summary AS (
    SELECT 
        dispatcher,
        COUNT(DISTINCT order_id) AS unique_orders,
        SUM(stop_count) AS total_stops_per_dispatcher
    FROM 
        Transportation_Tractor_Revenue
    WHERE 
        del_date >= (current_date - extract(dow FROM current_date)::integer)
    GROUP BY 
        dispatcher
),
bad_stops AS (
    SELECT
        fleet_manager AS dispatcher,
        COUNT(DISTINCT order_id) AS BAD_orders,
        COUNT(DISTINCT stop_id) AS BAD_stops
    FROM 
        bad_stops
    WHERE 
        sched_arrive_early >= (current_date - extract(dow FROM current_date)::integer)
    GROUP BY
        fleet_manager
),
date_range AS (
    SELECT
        CURRENT_DATE - EXTRACT(DOW FROM CURRENT_DATE)::integer AS start_date,  -- Most recent Sunday
        CURRENT_DATE AS end_date
),
date_diff AS (
    SELECT
        (end_date - start_date) AS days_in_range
    FROM
        date_range
)
SELECT 
    f.dispatcher,
    f.unique_trucks,
    --(f.total_miles / NULLIF(f.unique_trucks, 0)) AS miles_per_truck,  -- Average miles per truck
    (f.total_miles / NULLIF(f.unique_trucks, 0) / NULLIF(d.days_in_range, 0)) AS miles_per_truck_per_day,  -- Miles per truck per day
    (f.empty_miles / NULLIF(f.total_miles, 0)) * 100 AS deadhead_percentage,  -- Deadhead percentage
    ((f.total_orders - COALESCE(b.BAD_orders, 0)) * 1.0 / NULLIF(f.total_orders, 0)) * 100 AS order_percentage,
    ((s.total_stops_per_dispatcher - COALESCE(b.BAD_stops, 0)) * 1.0 / NULLIF(s.total_stops_per_dispatcher, 0)) * 100 AS stop_percentage  
    --COALESCE(b.BAD_orders, 0) AS BAD_orders,
    --COALESCE(b.BAD_stops, 0) AS BAD_stops
FROM 
    filtered_movements f
    LEFT JOIN order_summary s ON f.dispatcher = s.dispatcher
    LEFT JOIN bad_stops b ON f.dispatcher = b.dispatcher
    CROSS JOIN date_diff d
ORDER BY 
    f.dispatcher;
`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var data models.DailyOpsData
		err = rows.Scan(&data.Manager, &data.Trucks, &data.Miles, &data.Deadhead, &data.Order, &data.Stop)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		myData = append(myData, data)
	}

	return myData, nil
}

func WeekToDateDispatcherStats(db *sql.DB) ([]DispatcherData, error) {
	today := time.Now()
	startDate := helpers.GetStartDayOfWeek()

	workingDays, err := helpers.CountWorkingDaysInTimeframe(startDate, today)
	if err != nil {
		return nil, fmt.Errorf("error counting working days: %v", err)
	}
	if workingDays == 0 {
		workingDays = 1
	}

	query := `SELECT dispatcher, SUM(trucks) as trucks, SUM(total_miles) as miles, SUM(deadhead_percent) as deadhead, SUM(revenue) as revenue
	          FROM daily_driver_data
	          WHERE date >= $1 AND date <= $2
	          GROUP BY dispatcher`

	rows, err := db.Query(query, startDate, today)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	var data []DispatcherData

	// check if there are any rows
	if !rows.Next() {
		return nil, errors.New("no rows returned, dataStaleError")
	}

	for rows.Next() {
		var dispatcher string
		var trucks int
		var miles int
		var deadhead float64
		var revenue float64

		err = rows.Scan(&dispatcher, &trucks, &miles, &deadhead, &revenue)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		weekDeadhead := deadhead / float64(workingDays)
		MPTPD := float64((miles / trucks) / workingDays)
		RPTPD := float64(revenue / float64(trucks) / float64(workingDays))
		location := getLocation(dispatcher)

		dispatcherData := DispatcherData{
			Dispatcher:   dispatcher,
			Trucks:       trucks,
			Miles:        miles,
			WeekDeadhead: weekDeadhead,
			MPTPD:        MPTPD,
			RPTPD:        RPTPD,
			MPTPDColor:   getColor(MPTPD, goals[location]["MPTD"]),
			RPTPDColor:   getColor(RPTPD, goals[location]["RPTPD"]),
			DHColor:      getColorDH(weekDeadhead, goals[location]["DH"]),
			Location:     location,
		}

		fmt.Println(dispatcherData.Dispatcher)

		data = append(data, dispatcherData)
	}

	return data, nil
}

func getLocation(name string) string {
	// make name lowercase
	name = strings.ToLower(name)
	if name == "rochelle genera" {
		return "Texas"
	} else if name == "stephanie bingham" {
		return "Local"
	} else {
		return "OTR"
	}
}

func getColor(actual, goal float64) string {
	if actual >= goal {
		return "Green"
	} else if actual < goal {
		return "Red"
	}
	return "Black"
}

func getColorDH(actual, goal float64) string {
	if actual <= goal {
		return "Green"
	} else if actual > goal {
		return "Red"
	}
	return "Black"
}

func Add_OTWTDStats(db *sql.DB, data models.OTWTDStats) error {
	query := `INSERT INTO WTDOTStats(dispatcher, date, start_date, end_date, mptpd, rptpd, deadhead)
	          VALUES($1, $2, $3, $4, $5, $6, $7)`

	dispatcher := data.Dispatcher
	date := data.Date
	startDate := data.StartDate
	endDate := data.EndDate
	totalOrders := data.TotalOrders
	totalStops := data.TotalStops
	serviceIncidents := data.ServiceIncidents
	orderOnTime := data.OrderOnTime
	stopOnTime := data.StopOnTime

	_, err := db.Exec(query, dispatcher, date, startDate, endDate, totalOrders, totalStops, serviceIncidents, orderOnTime, stopOnTime)
	if err != nil {
		return fmt.Errorf("error inserting data into ot_wtd_stats: %v", err)
	}

	return nil
}

func AddBadStops(badStops []*models.BadStop) error {
	query := `
        INSERT INTO bad_stops (
            order_id, movement_id, actual_arrival, sched_arrive_early, id, sched_arrive_late,
            movement_sequence, equipment_group_id, dispatcher_user_id, equipment_id, equipment_type_id,
            fleet_manager, driver_id, stop_id, minutes_late, appt_required, stop_type, entered_user_id,
            entered_date, edi_standard_code, dsp_comment, sf_fault_of_carrier_or_driver, customer_id,
            operations_user, order_status
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25
        )
ON CONFLICT (id) DO UPDATE
        SET
            order_id = EXCLUDED.order_id,
            movement_id = EXCLUDED.movement_id,
            actual_arrival = EXCLUDED.actual_arrival,
            sched_arrive_early = EXCLUDED.sched_arrive_early,
            sched_arrive_late = EXCLUDED.sched_arrive_late,
            movement_sequence = EXCLUDED.movement_sequence,
            equipment_group_id = EXCLUDED.equipment_group_id,
            dispatcher_user_id = EXCLUDED.dispatcher_user_id,
            equipment_id = EXCLUDED.equipment_id,
            equipment_type_id = EXCLUDED.equipment_type_id,
            fleet_manager = EXCLUDED.fleet_manager,
            driver_id = EXCLUDED.driver_id,
            stop_id = EXCLUDED.stop_id,
            minutes_late = EXCLUDED.minutes_late,
            appt_required = EXCLUDED.appt_required,
            stop_type = EXCLUDED.stop_type,
            entered_user_id = EXCLUDED.entered_user_id,
            entered_date = EXCLUDED.entered_date,
            edi_standard_code = EXCLUDED.edi_standard_code,
            dsp_comment = EXCLUDED.dsp_comment,
            sf_fault_of_carrier_or_driver = EXCLUDED.sf_fault_of_carrier_or_driver,
            customer_id = EXCLUDED.customer_id,
            operations_user = EXCLUDED.operations_user,
            order_status = EXCLUDED.order_status
    `

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, stop := range badStops {
		_, err := stmt.Exec(
			stop.OrderID, stop.MovementID, stop.ActualArrival, stop.SchedArriveEarly, stop.ID,
			stop.SchedArriveLate, stop.MovementSequence, stop.EquipmentGroupID, stop.DispatcherUserID,
			stop.EquipmentID, stop.EquipmentTypeID, stop.FleetManager, stop.DriverID, stop.StopID,
			stop.MinutesLate, stop.ApptRequired, stop.StopType, stop.EnteredUserID, stop.EnteredDate,
			stop.EDIStandardCode, stop.DSPComment, stop.SFFaultOfCarrierOrDriver, stop.CustomerID,
			stop.OperationsUser, stop.OrderStatus,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
