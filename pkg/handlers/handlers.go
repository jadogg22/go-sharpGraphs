package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/cache"
	"github.com/jadogg22/go-sharpGraphs/pkg/database"
	"github.com/jadogg22/go-sharpGraphs/pkg/getData"
	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"

	"github.com/gin-gonic/gin"
)

// ---------- Transportation Handlers ----------

// This is the handler function for the transportation dashboard data, this function will evenutally return a bunch of data to display on the dash, right now we're working on HOS status.
func Dashboard(c *gin.Context) {
	// get the HOS data for the drivers
	cacheKey := "samsaraHOSData"
	cachedData, typeID, found := cache.MyCache.Get(cacheKey)
	if found {
		if typeID == "[]getdata.DriverInfo" {
			if cachedData, ok := cachedData.([]getdata.DriverInfo); ok {
				c.JSON(200, cachedData)
			}
		}
	}
	data := getdata.GetSamsaraHOSData()
	if data == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
		})
		return
	}

	//Finally update the Response with the json data
	c.JSON(200, data)
}

// This is the handler function for the transportation yearly revenue data, this function will return
// 52 weeks per year of the revenue earned to compair and contrast.
func Trans_year_by_year(c *gin.Context) {
	log.Printf("Trans_year_by_year handler called.")
	cacheKey := "transportationYearByYear"



	log.Printf("Attempting to get weekly revenue data from database.")
	data, err := database.GetWeeklyRevenueData()
	if err != nil {
		log.Printf("Error getting data from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"Error":   err,
		})
		return
	}
	log.Printf("Successfully retrieved %d weekly revenue entries from database.", len(data))

	log.Printf("Attempting to update transportation revenue data.")
	getdata.UpdateTransRevData(data)
	log.Printf("Finished updating transportation revenue data.")

	if len(data) == 0 {
		log.Printf("No data found for the week after update.")
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"Error":   "No data found for the week",
		})
		return
	}

	if len(data) == 53 {
		log.Printf("Data length is 53, trimming to 52.")
		data = data[:52]
	}

	log.Printf("Setting cache for key: %s", cacheKey)
	cache.MyCache.Set(cacheKey, data, "[]models.WeeklyRevenue", time.Hour*2)
	log.Printf("Cache set successfully.")

	c.JSON(200, gin.H{
		"Data": data,
	})
}

// This function returns the weeks/months/quarters
// REVENUE miles. So that we can make sure that we are
// Staying under the 10% empty miles.
func Trans_stacked_miles(c *gin.Context) {
	timePeriod := c.Param("when")

	//sanatize the input
	if timePeriod != "week_to_date" && timePeriod != "month_to_date" && timePeriod != "quarter" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid time period",
		})
		return
	}

	// check the cache
	cacheKey := "stackedMiles" + timePeriod
	cachedData, typeID, found := cache.MyCache.Get(cacheKey)
	if found {
		if typeID == "[]models.MilesData" {
			if cachedData, ok := cachedData.([]models.MilesData); ok {
				c.JSON(200, gin.H{timePeriod: cachedData})
				return
			}
		}
	}

	// Get the data from the database
	data, err := getdata.GetStackedMilesData(timePeriod)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   err,
		})
		return
	}

	milesData := helpers.StackedToMilesData(timePeriod, data)

	// Set the cache
	cache.MyCache.Set(cacheKey, milesData, "[]models.MilesData", time.Hour*2)

	c.JSON(200, gin.H{timePeriod: milesData})
}

type UnimplementedError struct {
	msg string
}

func (e *UnimplementedError) Error() string {
	return e.msg
}

func Trans_coded_revenue(c *gin.Context) {
	when := c.Param("when")
	fmt.Println("Getting coded revenue for ", when)

	//sanatize the input
	if when != "week" && when != "month" && when != "quarter" {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid time period",
		})
		return
	}
	// check the cache
	cacheKey := "codedRevenue" + when
	cachedData, typeID, found := cache.MyCache.Get(cacheKey)
	if found {
		if typeID == "[]models.CodedData" {
			if cachedData, ok := cachedData.([]models.CodedData); ok {
				c.JSON(200, gin.H{"data": cachedData})
				return
			}
		}
	}

	// Get the data from the database
	data, err := getdata.GetCodedRevenueData(when)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": data,
	})
	cache.MyCache.Set(cacheKey, data, "[]models.CodedData", time.Hour*2)
}

func Daily_Ops(c *gin.Context) {
	cacheKey := "dailyOpsData"
	cachedData, _, found := cache.MyCache.Get(cacheKey)
	if found {
		cachedData, ok := cachedData.([]*models.DailyOpsData)
		if ok {
			c.JSON(200, cachedData)
			return
		} else {
			fmt.Println("Error casting the data")
		}
	}
	// cache miss, get the data from the database.
	// get the start and end date for the current week
	fmt.Println("cache miss")
	startDate, endDate := helpers.GetWeekStartAndEndDates()
	data, err := getdata.GetTransportationDailyOps(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   fmt.Sprintf("%s", err),
		})
		return
	}

	if data == nil || len(data) <= 1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   fmt.Sprintf("No data found for the week"),
		})
		return
	}

	// Set the cache
	cache.MyCache.Set(cacheKey, data, "[]*models.DailyOpsData", time.Minute*45)

	fmt.Println("finished getting the data")
	fmt.Println(data)
	//Finally update the Response with the json data
	c.JSON(200, data)
}

func LaneProfit(c *gin.Context) {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startDate and endDate are required"})
		return
	}

	// 1. Get the data as a CSV file
	csvFilePath, err := getdata.GetOrderRevenueReportDataAsCSV(startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data as CSV",
			"error":   err.Error(),
		})
		return
	}
	defer os.Remove(csvFilePath) // Clean up the CSV file

	// 2. Create a temporary file for the PDF output
	pdfFile, err := os.CreateTemp("", "report-*.pdf")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error creating temporary PDF file",
			"error":   err.Error(),
		})
		return
	}
	pdfFilePath := pdfFile.Name()
	pdfFile.Close()
	defer os.Remove(pdfFilePath) // Clean up the PDF file

	// 3. Execute the Python script
	cmd := exec.Command("python3", "Statistics/main.py", csvFilePath, pdfFilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error executing Python script",
			"error":   err.Error(),
			"output":  string(output),
		})
		return
	}

	// 4. Send the PDF file to the client
	c.File(pdfFilePath)
}

// ---------- Logisitics Handlers ----------
func Log_year_by_year(c *gin.Context) {

	cacheKey := "logisticsYearByYear"
	cachedData, typeID, found := cache.MyCache.Get(cacheKey)
	if found {
		if typeID == "[]models.WeeklyRevenue" {
			if cachedData, ok := cachedData.([]models.WeeklyRevenue); ok {
				c.JSON(200, cachedData)
				return
			}
		}
		fmt.Println("Error casting the data")
	}

	//For now we're going to just get all data from the database
	//This data only includes finished weeks.
	fmt.Println("Getting the first data")
	// change to fectch data and use the new struct
	data, err := database.GetCachedData("logistics")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"Error":   fmt.Sprintf("%s", err),
		})
		return
	}

	fmt.Println("finished getting the first data")
	// Because its not very likly that we are at the end of the week
	// and we want to show the most recent data we need to check the
	// Transportation table and get the most recent data

	newData, err := database.GetYearByYearDataRefactored(data, "logistics")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting newest data from the database",
			"Error":   err,
		})
		return
	}

	c.JSON(200, gin.H{
		"Data": newData,
	})
}

func LogisticsMTD(c *gin.Context) {
	// check the cache
	cacheKey := "logisticsMTD"
	cachedData, typeID, found := cache.MyCache.Get(cacheKey)
	if found {
		if typeID == "[]models.LogisticsMTDStats" {
			if cachedData, ok := cachedData.([]models.LogisticsMTDStats); ok {
				c.JSON(200, cachedData)
				return
			}
		}
		fmt.Println("Error casting the data")
	}

	// cache miss, get the data from the database.
	startOfTheMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
	today := time.Now()
	data := getdata.GetLogisticsMTDData(startOfTheMonth, today)
	// need to add error handling here
	if data == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
		})
		return
	}

	// Set the cache
	cache.MyCache.Set(cacheKey, data, "[]models.LogisticsMTDStats", time.Minute*45)
	c.JSON(200, data)
}

// ---------- Dispatch Handlers ----------

func Dispach_week_to_date(c *gin.Context) {

	data, err := database.GetDispacherDataFromDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"Message": "i borke this",
		})
	}

	//Finally update the Response with the json data
	c.JSON(200, data)
}

func Dispatch_post(c *gin.Context) {
	// receive data from the client
	var data []models.DailyDriverData
	if err := c.ShouldBindJSON(&data); err != nil {
		// if there is an error return a 400 status code
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// driver add the data to the database
	var err error
	for _, driver := range data {
		err = database.Add_DailyDriverData(driver)
		if err != nil {

			fmt.Printf("error %v \n", err)
		}

	}

	c.JSON(200, gin.H{"Message": "Data received"})
}

// function for the vacation endpoint to calculate the vacation days for the drivers and office staff
func Vacation(c *gin.Context) {
	fmt.Println("Getting vacation days")

	typeID := c.Param("type")
	fmt.Println("Getting vacation days for ", typeID)

	switch typeID {
	case "tms", "tms2", "tms3, tms4":
		// get the vacation days for the drivers
		data, err := getdata.GetVacationFromDB(typeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Error getting data from the database",
				"error":   err,
			})
			return
		}
		// sort the data by the greatest amount of vacation days
		data = helpers.SortVacationData(data)
		c.JSON(200, gin.H{
			"Data": gin.H{typeID: data},
		})
	case "all":
		// get the vacation days for all the staff
		companyData := make(map[string][]models.VacationHours)
		companys := []string{"tms", "tms2", "tms3", "tms4", "drivers"}
		for _, company := range companys {
			data, err := getdata.GetVacationFromDB(company)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"Message": "Error getting data from the database",
					"error":   err,
				})
				return
			}
			// sort the data by the greatest amount of vacation days
			sortedData := helpers.SortVacationData(data)
			companyData[company] = sortedData
		}
		c.JSON(200, gin.H{
			"Data": companyData,
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid type",
		})
	}
}

// function for the sportsman endpoint to calculate the sportsman of the week
func Sportsman(c *gin.Context) {
	fmt.Println("Getting sportsman of the week")

	// get the sportsman of the week
	date1 := "2024-11-14"
	date2 := "2024-11-20"
	data, _ := getdata.GetSportsmanFromDB(date1, date2)

	c.JSON(200, gin.H{
		"Data": data,
	})
}

func SportsmanWithDates(c *gin.Context) {
	fmt.Println("Getting sportsman of the week")

	date1 := c.Param("date1")
	date2 := c.Param("date2")

	// check if the dates are valid
	_, err := time.Parse("2006-01-02", date1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid date format",
		})
		return
	}

	_, err = time.Parse("2006-01-02", date2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid date format",
		})
		return
	}

	// get the sportsman of the week
	data, _ := getdata.GetSportsmanFromDB(date1, date2)

	c.JSON(200, gin.H{
		"Data": data,
	})

}

type SportsmansRequest struct {
	OrderNumber string `json:"orderNumber"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

func Generate_Sportsmans(c *gin.Context) {
	var request SportsmansRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date1, err := time.Parse("2006-01-02T15:04:05Z", request.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid Startdate format",
		})
		return
	}

	date2, err := time.Parse("2006-01-02T15:04:05Z", request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Invalid Enddate format",
		})
		return
	}

	// format the dates to be strings in the correct format "2024-11-14"
	dateStr1 := date1.Format("2006-01-02")
	dateStr2 := date2.Format("2006-01-02")

	// get the sportsman of the week
	myData, err := getdata.GetSportsmanFromDB(dateStr1, dateStr2)
	if err != nil {
		if len(myData) < 1 {
			errMessage := fmt.Sprintf("No data found for the selected dates (Are you sure %s is the right bill date?)", dateStr1)
			c.JSON(http.StatusBadRequest, gin.H{
				"Message": errMessage,
				"error":   err,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Message": "Error getting data from the database",
				"error":   err,
			})
		}
		return
	}

	// Generate the excel file
	fmt.Println("Generating the excel file")
	spreadSheetData, err := helpers.GenerateSportsmansSpreadsheet(myData, request.OrderNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error generating the excel file",
			"error":   err,
		})
	}

	// Send back the excel file
	c.Header("Content-Description", "attachment; filename=Sportsman.xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", spreadSheetData)
}

func DriverManager(c *gin.Context) {
	// get the data from the database
	data, err := getdata.GetDriverManagerData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"Data": data,
	})
}
