package getdata

import (
	"testing"
	"time"
)

func TestTransportationDailyOps(t *testing.T) {
	t.Log("TestTransportationDailyOps")

	today := time.Now()
	//first day of the month
	currentYear, currentMonth, _ := today.Date()
	startDate := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, today.Location())

	myData, err := GetTransportationDailyOps(startDate, today)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(myData) < 1 {
		t.Fail()
	}

}

/*
func TestGetTransportationOrdersData(t *testing.T) {
	t.Log("TestGetTransportationOrdersData")

	// create a db connection

	conn, err := sql.Open("mssql", URL)
	if err != nil {
		fmt.Println("Error creating connection pool: " + err.Error())
		return
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		fmt.Println("Error pinging database: " + err.Error())
		return
	}

	RevenueReport := make(map[string]float64)

	today := time.Now()
	for i := 2020; i <= today.Year(); i++ {
		for month := time.January; month <= time.December; month++ {
			// get the start and end date of the month
			firstDay := time.Date(i, month, 1, 0, 0, 0, 0, today.Location())

			// break loop if first day is in the future
			if firstDay.After(today) {
				break
			}

			lastDay := firstDay.AddDate(0, 1, -1)

			data, err := GetTransportationOrdersData(conn, firstDay, lastDay)
			if err != nil {
				t.Error(err)
			}

			for _, d := range data {
				RevenueReport[d.WeekValue] += d.TotalRevenue
			}
		}

	}

	for k, v := range RevenueReport {
		t.Logf("Week: %s, Revenue: %f", k, v)
		if len(k) <= 5 {
			t.Error("Week value is not correct")
			continue
		}

		year, week := k[:4], k[len(k)-2:]

		// cast year and week to int
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			t.Error(err)
		}
		weekInt, err := strconv.Atoi(week)
		if err != nil {
			t.Error(err)
		}

		// add data to my database
		err = database.AddCacheData(v, weekInt, yearInt)
		if err != nil {
			t.Error(err)
		}
	}

	if len(RevenueReport) != 0 {
		t.Log("Test Failed")
	} else {
		t.Error("Test Passed")
	}

	conn.Close()
}
*/
