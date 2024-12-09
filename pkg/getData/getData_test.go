package getdata

import (
	"fmt"
	"testing"
)

func TestGetData(t *testing.T) {

	data := GetSportsmanFromDB("2024-11-25", "2024-12-05")

	for _, d := range data {
		// print start and end city
		fmt.Printf("Order: %s, City: %s, endCity: %s totalPallets: %d\n", d.Order_id, d.City, d.End_City, d.Total_Pallets)
		fmt.Printf("    freight: %f, fuel: %f, Pickup: %f, Dropoff: %f other: %f\n", d.Freight_Charges, d.Fuel_Surcharge, d.Extra_pickup, d.Extra_drops, d.OtherCharges)
	}
}
