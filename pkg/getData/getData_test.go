package getdata

import (
	"fmt"
	"testing"
)

func TestGetData(t *testing.T) {

	data := GetSportsmanFromDB()

	for _, d := range data {
		// print start and end city
		fmt.Printf("Order: %s, City: %s, endCity: %s\n", d.Order_id, d.City, d.End_City)
	}
}
