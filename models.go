package main

import (
	"time"
)

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
