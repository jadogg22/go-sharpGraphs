package getdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// get enviroment variables
	"github.com/joho/godotenv"
	"os"
)

// Tag represents a tag associated with a driver.
type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Driver represents a driver with ID, name, and associated tags.
type Driver struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Tags []Tag  `json:"tags"`
}

// Response represents the overall response structure.
type Response struct {
	Data []Driver `json:"data"`
}

type DriverInfo struct {
	ID        string
	Name      string
	DM        string
	HosStatus string
}

func GetSamsaraHOSData() []DriverInfo {
	// get a list of all drivers and their DMs
	driverMap := getDriverInfo()

	// get the HOS data for the drivers and transform into a slice
	mySlice := get_HOS_Data(driverMap)

	return mySlice
}

func getToken() string {
	// get the token from the environment variable
	// need to go up the directory to get the token
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	token := os.Getenv("SAMSARA_TOKEN")
	if token == "" {
		fmt.Println("Please set the SAMSARA_TOKEN environment variable")
		// panic if the token is not found
		panic("No token found")
	}
	return token
}

func getDriverInfo() *map[string]*DriverInfo {
	url := "https://api.samsara.com/fleet/drivers?driverActivationStatus=active&parentTagIds=&attributeValueIds="
	token := getToken()

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response Response
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println(err)
	}

	driverMap := make(map[string]*DriverInfo)

	for _, driver := range response.Data {
		driverName := driver.Name
		driverID := driver.ID
		var DMName string
		for _, tag := range driver.Tags {
			DMName = tag.Name
		}

		myDriverInfo := &DriverInfo{
			ID:   driverID,
			Name: driverName,
			DM:   DMName,
		}

		driverMap[driverID] = myDriverInfo
	}

	return &driverMap

}

func get_HOS_Data(driverMap *map[string]*DriverInfo) []DriverInfo {

	// Driver represents the driver information.
	type Driver struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// CurrentDutyStatus represents the current duty status of the driver.
	type CurrentDutyStatus struct {
		HosStatusType string `json:"hosStatusType"`
	}

	// Clocks represents the clocks information (not needed for your query but included for completeness).
	type Clocks struct {
		Break struct {
			TimeUntilBreakDurationMs int64 `json:"timeUntilBreakDurationMs"`
		} `json:"break"`
		Drive struct {
			DriveRemainingDurationMs int64 `json:"driveRemainingDurationMs"`
		} `json:"drive"`
		Shift struct {
			ShiftRemainingDurationMs int64 `json:"shiftRemainingDurationMs"`
		} `json:"shift"`
		Cycle struct {
			CycleStartedAtTime       string `json:"cycleStartedAtTime"`
			CycleRemainingDurationMs int64  `json:"cycleRemainingDurationMs"`
			CycleTomorrowDurationMs  int64  `json:"cycleTomorrowDurationMs"`
		} `json:"cycle"`
	}

	// DriverData represents the overall driver data structure.
	type DriverData struct {
		Driver            Driver            `json:"driver"`
		CurrentDutyStatus CurrentDutyStatus `json:"currentDutyStatus"`
		Violations        struct{}          `json:"violations"` // You can define this struct as needed
		Clocks            Clocks            `json:"clocks"`
	}

	type Response struct {
		Data []DriverData `json:"data"`
	}

	url := "https://api.samsara.com/fleet/hos/clocks?parentTagIds=&driverIds="
	token := getToken()

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response Response
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println(err)
	}

	// loop through the response and update the HOS status for each driver
	for _, driverData := range response.Data {
		driverID := driverData.Driver.ID
		if driver, exists := (*driverMap)[driverID]; exists {
			// if the driver exists in the map, then we can update the HOS status
			(*driver).HosStatus = driverData.CurrentDutyStatus.HosStatusType
		} else {
			// these guys are not active and so we can skip them
		}
	}

	var driverInfo_HOS []DriverInfo
	for key := range *driverMap {
		driverInfo := (*driverMap)[key]
		if driverInfo != nil {
			driverInfo_HOS = append(driverInfo_HOS, *driverInfo)
		}
	}

	/*
		for _, driver := range driverInfo_HOS {
			fmt.Println("--------------------")
			fmt.Println("Driver ID: ", driver.ID)
			fmt.Println("Driver Name: ", driver.Name)
			fmt.Println("Driver DM: ", driver.DM)
			fmt.Println("Driver HOS Status: ", driver.HosStatus)
		}
	*/

	return driverInfo_HOS
}
