package getdata

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

//func TestGetData(t *testing.T) {
//	RunUpdater()
//}

func TestGetLogisticsMTDData(t *testing.T) {

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

	data := GetLogisticsMTDData(time.Now(), time.Now())
	if len(data) == 0 {
		t.Errorf("expected data to have length > 0")
	}

	if data == nil {
		t.Errorf("expected data to not be nil")
	}

	for _, v := range data {
		fmt.Println(v.Dispacher)
		fmt.Println(v.Revenue)
		fmt.Println(v.TotalMiles)

	}

}
