package models

import (
	"database/sql"
	"strings"
	"time"
)

type OrderDetails struct {
	OrderID          string    `json:"order_id"`
	OperationsUser   string    `json:"operations_user"`
	RevenueCodeID    string    `json:"revenue_code_id"`
	FreightCharge    float64   `json:"freight_charge"`
	BillMiles        float64   `json:"bill_miles"`
	BillDate         time.Time `json:"bill_date"`
	OrderTrailerType string    `json:"order_trailer_type"`
	OriginValue      string    `json:"origin_value"`
	DestinationValue string    `json:"destination_value"`
	CustomerID       string    `json:"customer_id"`
	CustomerName     string    `json:"customer_name"`
	CustomerCategory string    `json:"customer_category"`
	CategoryDescr    string    `json:"category_descr"`
	MovementID       string    `json:"movement_id"`
	Loaded           string    `json:"loaded"`
	MoveDistance     float64   `json:"move_distance"`
	Brokerage        string    `json:"brokerage"`
	OriginCity       string    `json:"origin_city"`
	OriginState      string    `json:"origin_state"`
	DestCity         string    `json:"dest_city"`
	DestState        string    `json:"dest_state"`
	ReportDate       time.Time `json:"report_date"`
	ActualDate       time.Time `json:"actual_date"`
	EmptyMiles       float64   `json:"empty_miles"`
	LoadedMiles      float64   `json:"loaded_miles"`
	TotalMiles       float64   `json:"total_miles"`
	TotalRevenue     float64   `json:"total_revenue"`
	WeekValue        string    `json:"week_value"`
	MonthValue       string    `json:"month_value"`
	QuarterValue     string    `json:"quarter_value"`
	DetailID         string    `json:"detail_id"`
}

type TransportationOrder struct {
	TotalRevenue     sql.NullFloat64 `json:"total_revenue"`
	BillMiles        sql.NullInt64   `json:"bill_miles"`
	LoadedMiles      sql.NullInt64   `json:"loaded_miles"`
	Billed           sql.NullBool    `json:"billed"`
	EmptyMiles       sql.NullInt64   `json:"empty_miles"`
	TotalMiles       sql.NullInt64   `json:"total_miles"`
	EmptyPercentage  sql.NullFloat64 `json:"empty_pct"`
	RevLoadedMile    sql.NullFloat64 `json:"rev_per_loaded_mile"`
	RevTotalMile     sql.NullFloat64 `json:"rev_per_total_mile"`
	Freight          sql.NullFloat64 `json:"freight"`
	FuelSurcharge    sql.NullFloat64 `json:"fuel_surcharge"`
	RemainingCharges sql.NullFloat64 `json:"remaining_charges"`
	Brokered         sql.NullBool    `json:"brokered"`
	DestinationState sql.NullString  `json:"destination_state"`
	Week             sql.NullString  `json:"week"`
	Month            sql.NullString  `json:"month"`
	Quarter          sql.NullString  `json:"quarter"`
	Origin           sql.NullString  `json:"origin"`
	OrderNumber      string          `json:"order"`
	OrderType        sql.NullString  `json:"order_type"`
	DeliveryDate     sql.NullString  `json:"delivery_date"`
	RevenueCode      sql.NullString  `json:"revenue_code"`
	Destination      sql.NullString  `json:"destination"`
	Customer         sql.NullString  `json:"customer"`
	CustomerCategory sql.NullString  `json:"customer_category"`
	OperationsUser   sql.NullString  `json:"operations_user"`
	ControllingParty sql.NullString  `json:"controlling_party"`
	Commodity        sql.NullString  `json:"commodity"`
	TrailerType      sql.NullString  `json:"trailer_type"`
	OriginState      sql.NullString  `json:"origin_state"`
}

type TractorRevenue struct {
	MoveID           int
	MoveDistance     float64
	Loaded           string
	OrderID          int
	Charges          float64
	BillDistance     float64
	FreightCharge    float64
	OriginCity       string
	OriginState      sql.NullString
	EquipID          sql.NullString
	ActualArrival    time.Time
	DelDate          time.Time
	Tractor          sql.NullString
	EquipmentTypeID  sql.NullString
	Dispatcher       sql.NullString
	FleetID          sql.NullString
	FleetDescription sql.NullString
	UserName         sql.NullString
	ServiceFailCount int
	HasServiceFail   bool
	StopCount        int
}

type BadStop struct {
	OrderID                  string
	MovementID               string
	ActualArrival            time.Time
	SchedArriveEarly         *time.Time
	ID                       string
	SchedArriveLate          time.Time
	MovementSequence         int
	EquipmentGroupID         string
	DispatcherUserID         string
	EquipmentID              string
	EquipmentTypeID          string
	FleetManager             string
	DriverID                 string
	StopID                   string
	MinutesLate              int
	ApptRequired             string
	StopType                 string
	EnteredUserID            string
	EnteredDate              time.Time
	EDIStandardCode          string
	DSPComment               string
	SFFaultOfCarrierOrDriver string
	CustomerID               string
	OperationsUser           string
	OrderStatus              string
}

// This is every stop in the mcloud database and this should be enough ingo
// to get the data that we need for the front end projects
type LogisticsStops struct {
	ID                       string     `json:"id"`
	OrderID                  string     `json:"order_id"`
	MovementID               string     `json:"movement_id"`
	ActualArrival            *time.Time `json:"actual_arrival"`
	SchedArriveEarly         *time.Time `json:"sched_arrive_early"`
	SchedArriveLate          *time.Time `json:"sched_arrive_late"`
	MovementSequence         int        `json:"movement_sequence"`
	DispatcherUserID         *string    `json:"dispatcher_user_id"`
	FleetManager             *string    `json:"fleet_manager"`
	DriverID                 *string    `json:"driver_id"`
	ServiceFailStopID        *string    `json:"servicefail_stop_id"`
	MinutesLate              *float64   `json:"minutes_late"`
	ApptRequired             *string    `json:"appt_required"`
	StopType                 *string    `json:"stop_type"`
	EnteredUserID            *string    `json:"entered_user_id"`
	EnteredDate              *time.Time `json:"entered_date"`
	EDIStandardCode          *string    `json:"edi_standard_code"`
	DSPComment               *string    `json:"dsp_comment"`
	SFFaultOfCarrierOrDriver *string    `json:"sf_fault_of_carrier_or_driver"`
	OverridePayeeID          *string    `json:"override_payee_id"`
	CustomerID               *string    `json:"customer_id"`
	PayAmt                   *float64   `json:"pay_amt"`
	OrdersID                 *string    `json:"orders_id"`
	OperationsUser           *string    `json:"operations_user"`
	TotalCharge              *float64   `json:"total_charge"`
	OrderStatus              *string    `json:"order_status"`
	BillDistance             *float64   `json:"bill_distance"`
	AdditionalAmount         *float64   `json:"additional_amount"`
	TruckHire                *float64   `json:"truck_hire"`
	DeductCode               *string    `json:"deduct_code"`
}

type LogisticsOrdersData struct {
	Dispacher  string
	Truck_hire float64
	Charges    float64
	Miles      float64
}

type LogisticsStopOrdersData struct {
	Dispacher    sql.NullString
	Total_stops  int
	Total_orders int
	Order_faults int
	Stop_faults  int
}

type LogisticsMTDStats struct {
	Dispacher       string  `json:"dispacher"`
	TotalOrders     int     `json:"total_orders"`
	Revenue         float64 `json:"revenue"`
	TruckHire       float64 `json:"truck_hire"`
	NetRevenue      float64 `json:"net_revenue"`
	Margins         float64 `json:"margins"`
	TotalMiles      float64 `json:"total_miles"`
	RevPerMile      float64 `json:"rev_per_mile"`
	StopPercentage  float64 `json:"stop_percentage"`
	OrderPercentage float64 `json:"order_percentage"`
}

func NewLogisticsMTDStats(dispacher string, truck_hire, charges, miles float64, total_stops, total_orders, Order_faults, stop_faults int) *LogisticsMTDStats {
	if total_orders == 0 {
		total_orders = 1
	}

	if total_stops == 0 {
		total_stops = 1
	}

	if charges == 0 {
		charges = 1
	}

	if miles == 0 {
		miles = 1
	}

	var margins float64
	if charges-truck_hire == 0 {
		margins = 1
	} else {
		margins = (charges - truck_hire) / charges
	}

	var rev_per_mile float64
	if miles == 0 {
		rev_per_mile = 0
	} else {
		rev_per_mile = (charges - truck_hire) / miles
	}

	var stop_percentage float64
	if total_stops == 0 {
		stop_percentage = 100.0
	} else if stop_faults == 0 {
		stop_percentage = 100.0
	} else {
		stop_percentage = ((float64(total_stops) - float64(stop_faults)) / float64(total_stops) * 100)
	}

	var order_percentage float64
	if total_orders == 0 {
		order_percentage = 100.0
	} else if Order_faults == 0 {
		order_percentage = 100.0
	} else {
		order_percentage = ((float64(total_orders) - float64(Order_faults)) / float64(total_orders) * 100)
	}

	return &LogisticsMTDStats{
		Dispacher:       dispacher,
		TotalOrders:     total_orders,
		Revenue:         charges,
		TruckHire:       truck_hire,
		NetRevenue:      charges - truck_hire,
		Margins:         margins,
		TotalMiles:      miles,
		RevPerMile:      rev_per_mile,
		StopPercentage:  stop_percentage,
		OrderPercentage: order_percentage,
	}

}

type DailyOpsData struct {
	Dispatcher     string  `json:"driverManager"`
	NumberOfTrucks int     `json:"numberOfTrucks"`
	MilesPerTruck  float64 `json:"milesPerTruck"`
	Deadhead       float64 `json:"deadhead"`
	OrderPercent   float64 `json:"order"`
	StopPercent    float64 `json:"stop"`
}

func NewDailyOpsDataFromDB(dispatcher string, total_loaded_distance, total_empty_distance sql.NullFloat64, total_stops, total_servicefail_count, orders_with_service_fail, total_orders, total_unique_trucks int) *DailyOpsData {
	// avoid divide by zero
	if total_unique_trucks == 0 {
		total_unique_trucks = 1
	}
	if total_orders == 0 {
		total_orders = 1
	}
	if total_stops == 0 {
		total_stops = 1
	}
	if total_loaded_distance.Float64 == 0 || !total_loaded_distance.Valid {
		total_loaded_distance.Float64 = 1
	}
	if total_empty_distance.Float64 == 0 || !total_empty_distance.Valid {
		total_empty_distance.Float64 = 1
	}

	total_distance := total_loaded_distance.Float64 + total_empty_distance.Float64

	// calculate miles per truck
	var miles_per_truck float64
	// avoid divide by zero error
	if total_unique_trucks == 0 {
		miles_per_truck = 0
	} else {
		miles_per_truck = total_distance / float64(total_unique_trucks)
	}
	// calculate the deadhead percentage
	var deadhead float64
	if total_distance == 0 {
		deadhead = 0
	} else {
		deadhead = ((total_distance - total_empty_distance.Float64) / total_distance)
		deadhead = deadhead * 100
	}

	// calculate the order percentage
	var order_percent float64
	if total_orders == 0 {
		order_percent = 1.0
	} else if orders_with_service_fail == 0 {
		order_percent = 1.0
	} else {
		order_percent = ((float64(total_orders) - float64(orders_with_service_fail)) / float64(total_orders))
	}

	// calculate the stop percentage
	var stop_percent float64
	if total_stops == 0 {
		stop_percent = 1.0
	} else if total_servicefail_count == 0 {
		stop_percent = 1.0
	} else {
		stop_percent = ((float64(total_stops) - float64(total_servicefail_count)) / float64(total_stops))
	}

	return &DailyOpsData{
		Dispatcher:     dispatcher,
		NumberOfTrucks: total_unique_trucks,
		MilesPerTruck:  miles_per_truck,
		Deadhead:       deadhead,
		OrderPercent:   order_percent,
		StopPercent:    stop_percent,
	}
}

type LoadData struct {
	RevenueCode      string  `json:"revenue_code"`
	Order            string  `json:"order"`
	OrderType        string  `json:"order_type"`
	Freight          int     `json:"freight"`
	FuelSurcharge    float64 `json:"fuel_surcharge"`
	RemainingCharges float64 `json:"remaining_charges"`
	TotalRevenue     float64 `json:"total_revenue"`
	BillMiles        int     `json:"bill_miles"`
	LoadedMiles      int     `json:"loaded_miles"`
	EmptyMiles       int     `json:"empty_miles"`
	TotalMiles       int     `json:"total_miles"`
	EmptyPct         float64 `json:"empty_pct"`
	RevPerLoadedMile float64 `json:"rev_per_loaded_mile"`
	RevPerTotalMile  float64 `json:"rev_per_total_mile"`
	DeliveryDate     string  `json:"delivery_date"`
	Origin           string  `json:"origin"`
	Destination      string  `json:"destination"`
	Customer         string  `json:"customer"`
	CustomerCategory string  `json:"customer_category"`
	OperationsUser   string  `json:"operations_user"`
	Billed           string  `json:"billed"`
	ControllingParty string  `json:"controlling_party"`
	Commodity        string  `json:"commodity"`
	TrailerType      string  `json:"trailer_type"`
	OriginState      string  `json:"origin_state"`
	DestinationState string  `json:"destination_state"`
	Week             string  `json:"week"`
	Month            string  `json:"month"`
	Quarter          string  `json:"quarter"`
	Brokered         string  `json:"brokered"`
}

type DriverData struct {
	Dispatcher                string      `json:"dispatcher"`
	Deadhead_percent          []float64   `json:"Deadhead"`
	Freight                   []float64   `json:"Freight"`
	Fuel_Surcharge            []float64   `json:"Fuel_Surcharge"`
	Remain_Chgs               []float64   `json:"Remain_Chgs"`
	Revenue                   []float64   `json:"Revenue"`
	Total_Rev_per_rev_miles   []float64   `json:"Total_Rev_per_rev_miles"`
	Total_Rev_per_Total_Miles []float64   `json:"Total_Rev_per_Total_Miles"`
	Average_weekly_rev        []float64   `json:"Average_weekly_rev"`
	Average_weekly_Rev_Miles  []float64   `json:"Average_weekly_Rev_Miles"`
	Average_rev_miles         []float64   `json:"Average_rev_miles"`
	Revenue_Miles             []float64   `json:"Revenue_Miles"`
	Total_Miles               []float64   `json:"Total_Miles"`
	Trucks                    []int64     `json:"Trucks"`
	Date                      []time.Time `json:"Date"`
}

type DailyDriverData struct {
	Dispatcher                string
	Deadhead_percent          float64
	Freight                   float64
	Fuel_Surcharge            float64
	Remain_Chgs               float64
	Revenue                   float64
	Total_Rev_per_rev_miles   float64
	Total_Rev_per_Total_Miles float64
	Average_weekly_rev        float64
	Average_weekly_Rev_Miles  float64
	Average_rev_miles         float64
	Revenue_Miles             float64
	Total_Miles               float64
	Trucks                    int64
	Date                      time.Time
}

type CodedRevenueData struct {
	Code    []string    `json:"Code"`
	Revenue []float64   `json:"Revenue"`
	Date    []time.Time `json:"Date"`
}

type CodedData struct {
	Name    string
	Revenue float64
	Count   int
}

type MilesData struct {
	DeliveryDate     string  `json:"DeliveryDate"`
	Name             string  `json:"Name"`
	NameStr          string  `json:"NameStr"`
	TotalLoadedMiles float64 `json:"Total_Loaded_Miles"`
	TotalEmptyMiles  float64 `json:"Total_Empty_Miles"`
	TotalMiles       float64 `json:"Total_Actual_Miles"`
	PercentEmpty     float64 `json:"Percent_empty"`
}

type OTWTDStats struct {
	Dispatcher       string
	Date             time.Time
	StartDate        time.Time
	EndDate          time.Time
	TotalOrders      int
	TotalStops       int
	ServiceIncidents int
	OrderOnTime      float32
	StopOnTime       float32
}

type WeeklyRevenue struct {
	Name        int      `json:"Name"`
	Revenue2021 *float64 `json:"2021 Revenue,omitempty"`
	Revenue2022 *float64 `json:"2022 Revenue,omitempty"`
	Revenue2023 *float64 `json:"2023 Revenue,omitempty"`
	Revenue2024 *float64 `json:"2024 Revenue,omitempty"`
	Revenue2025 *float64 `json:"2025 Revenue,omitempty"`
}

func (wr *WeeklyRevenue) GetRevenue(year int) *float64 {
	switch year {
	case 2021:
		return wr.Revenue2021
	case 2022:
		return wr.Revenue2022
	case 2023:
		return wr.Revenue2023
	case 2024:
		return wr.Revenue2024
	default:
		return nil
	}
}

type VacationHours struct {
	EmployeeID        string
	EmployeeName      string
	VacationHoursDue  string
	VacationHoursRate string // sql null float
	AmountDue         string
}

// stacked miles data
type StackedMilesData struct {
	ID          string
	Date        string
	EmptyMiles  float64
	LoadedMiles float64
}

type SportsmanData struct {
	Order_id              string
	Order_date            string
	Delivery_Date         string
	Bill_date             string
	City                  string
	State                 string
	Zip                   string
	End_City              string
	End_State             string
	End_Zip               string
	Consignee             string
	Miles                 string
	BOL_Number            string
	Commodity             string
	Weight                string
	Movement              int64
	Total_Pallets         int64
	Pallets_Droped        int64
	Pallets_Picked        int64
	Freight_Charges       float64
	Fuel_Surcharge        float64
	Detention_and_layover float64
	OtherCharges          float64
	Total_Charges         float64
}

// Helper function to handle NullString to default "N/A"
func nullStringToStr(ns sql.NullString, defaultVal string) string {
	if ns.Valid {

		return strings.TrimRight(ns.String, " ")
	}
	return defaultVal
}

// Helper function to handle NullInt64 to default value
func nullInt64ToInt(n sql.NullInt64, defaultVal int64) int64 {
	if n.Valid {
		return n.Int64
	}
	return defaultVal
}

// Helper function to handle NullFloat64 to default value
func nullFloat64ToFloat(n sql.NullFloat64, defaultVal float64) float64 {
	if n.Valid {
		return n.Float64
	}
	return defaultVal
}

func NewSportsmanData(order_id, ordered_date, delivery_date, bill_date, city, state, zip, end_city, end_state, end_zip, consignee, miles, bol_number, commodity, weight sql.NullString, movement, pallets_droped, pallets_picked, TotalPallets sql.NullInt64, freight_charges, other_charges, total_charges, fuel_surcharge, Detention_and_layover sql.NullFloat64) *SportsmanData {
	return &SportsmanData{
		Order_id:              nullStringToStr(order_id, "N/A"),
		Order_date:            nullStringToStr(ordered_date, "N/A"),
		Delivery_Date:         nullStringToStr(delivery_date, "N/A"),
		Bill_date:             nullStringToStr(bill_date, "N/A"),
		City:                  nullStringToStr(city, "N/A"),
		State:                 nullStringToStr(state, "N/A"),
		Zip:                   nullStringToStr(zip, "N/A"),
		End_City:              nullStringToStr(end_city, "N/A"),
		End_State:             nullStringToStr(end_state, "N/A"),
		End_Zip:               nullStringToStr(end_zip, "N/A"),
		Consignee:             nullStringToStr(consignee, "N/A"),
		Miles:                 nullStringToStr(miles, "N/A"),
		BOL_Number:            nullStringToStr(bol_number, "N/A"),
		Commodity:             nullStringToStr(commodity, "N/A"),
		Weight:                nullStringToStr(weight, "N/A"),
		Movement:              nullInt64ToInt(movement, 0),
		Total_Pallets:         nullInt64ToInt(TotalPallets, 0),
		Pallets_Droped:        nullInt64ToInt(pallets_droped, 0),
		Pallets_Picked:        nullInt64ToInt(pallets_picked, 0),
		Freight_Charges:       nullFloat64ToFloat(freight_charges, 0),
		Fuel_Surcharge:        nullFloat64ToFloat(fuel_surcharge, 0),
		Detention_and_layover: nullFloat64ToFloat(Detention_and_layover, 0),
		OtherCharges:          nullFloat64ToFloat(other_charges, 0),
		Total_Charges:         nullFloat64ToFloat(total_charges, 0),
	}
}
