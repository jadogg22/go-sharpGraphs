package database

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
)

func TestCountWorkingDays(t *testing.T) {
	// Test for a single day
	fmt.Println("testing working days for the dispacher functions")

	workingDays, _ := helpers.CountWorkingDays()

	// in june there are 20 working days
	// get the month we are currently in
	month := time.Now().Month()
	if month == 6 {
		if workingDays != 20 {
			t.Errorf("Expected 20 working days, got %d", workingDays)
		}
	}
}

func TestGetStartDayOfWeek(t *testing.T) {
	// Test for a single day
	fmt.Println("testing start day of the week for the dispacher functions")

	startDate := helpers.GetStartDayOfWeek()

	// get the month we are currently in
	month := time.Now().Month()
	if month == 6 {
		if startDate.Weekday() != time.Sunday {
			t.Errorf("Expected Monday, got %s", startDate.Weekday())
		}
	}
}

func TestDatabaseConnection(t *testing.T) {
	// get the db path
	dbPath := filepath.Join("..", "..", "Data", "Production.db")

	db, err := test_make_connection(dbPath)
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	query := `SELECT * FROM daily_driver_data
			ORDER BY date DESC
			LIMIT 1;`

	rows, err := db.Query(query)
	if err != nil {
		t.Fatalf("error querying database: %v", err)
	}
	defer rows.Close()

	// should be one row
	if !rows.Next() {
		t.Fatalf("no rows returned")
	}

	// get the columns
	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("error getting columns: %v", err)
	}

	// get the values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// scan the values
	if err := rows.Scan(valuePtrs...); err != nil {
		t.Fatalf("error scanning values: %v", err)
	}

	// print the values
	for i, col := range columns {
		fmt.Printf("%s: %s\n", col, values[i])
	}

	defer db.Close()
}

func TestWeekToDateDispacherStats(t *testing.T) {
	// get the db path
	dbPath := filepath.Join("..", "..", "Data", "Production.db")

	db, err := test_make_connection(dbPath)
	if err != nil {
		t.Fatalf("error opening database: %v", err)
	}

	data, err := WeekToDateDispatcherStats(db)
	if err != nil {
		switch err.Error() {
		case "no rows returned, dataStaleError":
			fmt.Println("Data is stale")
		default:
			t.Fatalf("error getting week to date dispatcher stats: %v", err)
		}
	}

	if len(data) == 0 {
		t.Errorf("Got no data")
	}

	// print the data nicely formatted

	for item := range data {
		fmt.Printf("Dispatcher: %s\n", data[item].Dispatcher)
		fmt.Printf("\tTrucks: %d\n", data[item].Trucks)
		fmt.Printf("\tMiles: %d\n", data[item].Miles)
		fmt.Printf("\tWeek Deadhead: %f\n", data[item].WeekDeadhead)
		fmt.Printf("\tMPTPD: %f\n", data[item].MPTPD)
		fmt.Printf("\tRPTPD: %f\n", data[item].RPTPD)
		fmt.Printf("\tLocation %s\n", data[item].Location)
	}
}

func TestGetColor(t *testing.T) {
	// this is pretty odvious for the most part but im done and I need to actually have the same for Deadhead

	DHGoal1 := 10.0
	//DHGoal2 := 45.0
	//DHGoal3 := 30.0

	test1 := getColor(-5.0+100.0, DHGoal1)
	if test1 != "green" {
		t.Errorf("Expected green, got %s", test1)
	}
	test2 := getColor(-15.0+100.0, DHGoal1)
	if test2 != "yellow" {
		t.Errorf("Expected yellow, got %s", test2)
	}
	test3 := getColor(-35.0+100.0, DHGoal1)
	if test3 != "red" {
		t.Errorf("Expected red, got %s", test3)
	}
}
