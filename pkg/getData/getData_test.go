package getdata

import (
	"fmt"
	"testing"
	"time"
)

func TestGetLogisticsMTDData(t *testing.T) {

	t.Log("TestGetLogisticsMTDData")
	today := time.Now()
	// start of the month
	start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

	t.Log("start: ", start)
	t.Log("today: ", today)

	data := GetLogisticsMTDData(start, today)

	for _, d := range data {
		str := fmt.Sprintf(`
--------------------------
Dispacher: %s
TotalOrders: %d
Revenue: %f
TruckHire: %f
NetRevenue: %f
Margin: %f
TotalMiles: %f
RevenuePerMile: %f
StopPercent: %f
OrderPercent: %f
`, d.Dispacher, d.TotalOrders, d.Revenue, d.TruckHire, d.NetRevenue, d.Margins, d.TotalMiles, d.RevPerMile, d.StopPercentage, d.OrderPercentage)
		t.Log(str)
	}
}
