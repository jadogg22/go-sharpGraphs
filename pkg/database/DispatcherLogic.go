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
        user_name AS dispatcher,
        SUM(move_distance) AS total_miles,
        COUNT(DISTINCT tractor) AS unique_trucks,
        SUM(CASE WHEN loaded = 'E' THEN move_distance ELSE 0 END) AS empty_miles,
        COUNT(DISTINCT order_id) AS total_orders,
        COUNT(DISTINCT CASE WHEN service_fail_count > 0 THEN order_id END) AS orders_with_servicefail
    FROM 
        Transportation_Tractor_Revenue
    WHERE 
        del_date >= (current_date - extract(dow from current_date)::integer)
    GROUP BY 
        user_name
),
order_summary AS (
    SELECT 
        user_name AS dispatcher,
        order_id,
        SUM(stop_count) AS total_stops_per_order,
        SUM(service_fail_count) AS total_servicefail_per_order
    FROM 
        Transportation_Tractor_Revenue
    WHERE 
        del_date >= (current_date - extract(dow from current_date)::integer)
    GROUP BY 
        user_name, order_id
),
dispatcher_summary AS (
    SELECT 
        dispatcher,
        COUNT(DISTINCT order_id) AS unique_orders,
        SUM(total_stops_per_order) AS total_stops,
        SUM(total_servicefail_per_order) AS total_servicefail
    FROM 
        order_summary
    GROUP BY 
        dispatcher
)
SELECT 
    f.dispatcher,
    f.total_miles,
    f.unique_trucks,
    f.empty_miles,
    ((f.total_miles - f.empty_miles * 1.0) / f.total_miles) * 100 AS deadhead_percentage,
    f.orders_with_servicefail * 1.0 / f.total_orders) * 100 AS order_percentage,
    s.total_stops,
    s.total_servicefail,
    s.total_servicefail * 1.0 / s.total_stops * 100 AS stop_percentage
FROM 
    filtered_movements f
    LEFT JOIN dispatcher_summary s ON f.dispatcher = s.dispatcher
ORDER BY 
    f.dispatcher;
`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}

	defer rows.Close()

	var EmptyMiles float64
	for rows.Next() {
		var data models.DailyOpsData
		err = rows.Scan(&data.Manager, &data.Miles, &data.Trucks, &EmptyMiles, &data.Deadhead, &data.Order, &data.Stop)
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
