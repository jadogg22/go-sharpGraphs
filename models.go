package main

type DriverData struct {
	Driver               string    `json:"driver"`
	Day                  []string  `json:"day"`
	TotalTrucks          []int     `json:"# Total Trucks"`
	TotalMiles           []int     `json:"Total Miles"`
	LoadedMiles          []int     `json:"Loaded Miles"`
	RevenueWithoutFuel   []float64 `json:"Revenue w/o fuel"`
	MilesPerTruck        []int     `json:"Miles Per Truck"`
	RevenuePerTruck      []float64 `json:"Revenue Per Truck"`
	Deadhead             []float64 `json:"Deadhead"`
	RevenuePerLoadedMile []float64 `json:"Revenue Per Loaded Mile"`
	OnTimeService        []float64 `json:"On-Time Service %"`
}

type YearByYear struct {
	Name        string  `json:"Name"`
	Revenue2021 float64 `json:"2021 Revenue"`
	Revenue2022 float64 `json:"2022 Revenue"`
	Revenue2023 float64 `json:"2023 Revenue"`
	Revenue2024 float64 `json:"2024 Revenue"`
}
