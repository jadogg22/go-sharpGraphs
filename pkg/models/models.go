package models

import (
	"database/sql"
	"time"
)

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

type DailyOpsData struct {
	Manager  string  `json:"driverManager"`
	Trucks   int     `json:"numberOfTrucks"`
	Miles    float64 `json:"milesPerTruck"`
	Deadhead float64 `json:"deadhead"`
	Order    float64 `json:"order"`
	Stop     float64 `json:"stop"`
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
}
