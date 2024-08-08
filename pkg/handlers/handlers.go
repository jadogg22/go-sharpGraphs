package handlers

import (
	"fmt"
	"net/http"

	"github.com/jadogg22/go-sharpGraphs/pkg/database"
	"github.com/jadogg22/go-sharpGraphs/pkg/getData"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"

	"github.com/gin-gonic/gin"
)

// Handler function for	root route for debuging
func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "This came from the test handler",
	})
}

func Test_db(c *gin.Context) {
	// connect to the database
	db, err := database.PG_Make_connection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error connecting to the database",
		})
		return
	}

	rows, err := db.Query("SELECT count(*) FROM transportation")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error querying the database",
		})
		return
	}

	defer rows.Close()

	var count int

	for rows.Next() {
		rows.Scan(&count)
	}

	c.JSON(200, gin.H{
		"message": "Successfully connected to the database, there are " + fmt.Sprint(count) + " rows in the transportation table"})

}

// ---------- Transportation Handlers ----------
func Trans_year_by_year(c *gin.Context) {

	//For now we're going to just get all data from the database
	//This data only includes finished weeks.
	fmt.Println("Getting the first data")
	// change to fectch data and use the new struct
	data, err := database.GetCachedData("transportation")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"Error":   err,
		})
		return
	}

	fmt.Println("finished getting the first data")
	// Because its not very likly that we are at the end of the week
	// and we want to show the most recent data we need to check the
	// Transportation table and get the most recent data

	newData, err := database.GetYearByYearDataRefactored(data, "transportation")
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

// This function returns the weeks/months/quarters
// REVENUE miles. So that we can make sure that we are
// Staying under the 10% empty miles.
func Trans_stacked_miles(c *gin.Context) {
	timePeriod := c.Param("when")

	data, err2 := database.GetMilesData(timePeriod, "transportation")
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   err2,
		})
		return
	}

	c.JSON(200, data) // We're gonna come up with something better here
}

func Trans_coded_revenue(c *gin.Context) {
	when := c.Param("when")
	fmt.Println("Getting coded revenue for ", when)

	// parts := strings.Split(when, "-")

	// if len(parts) == 1 {
	// 	fmt.Println("TODO, write funtion for coded revnue for one peram")
	// }

	// if len(parts) == 2 {
	// 	fmt.Println("TODO, write function for coded revenue from one date to another.")
	// }

	// if len(parts) < 1 && len(parts) > 2 {
	// 	fmt.Println("Sorry but, WTF")
	// }

	data, revenue, count, err2 := database.GetCodedRevenueData(when)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"error":   err2,
		})
		return
	}

	c.JSON(200, gin.H{
		"CodedRevenue": data,
		"TotalRevenue": revenue,
		"TotalCount":   count,
	})
}

// This function returns the daily operations data
// for the transportation department

// Daily ops page includes two tables of data.
// I'm not really super worried about bandwidth so for now I'm going to just send two arays of maps

// first table is | Manager | # trucks | Miles | Deadhead | order | stop |
// secound table is | manager | Average MPTPD | Average RPTPD | DH% | ORDER OTP |STOP OTP | AVG MPTPD Needed to Make Goal

// I think for the second table we're just going to include the color of the data like item.AverageMPTPDCOlOR: "Green"

// then on the front end we're going to be able to use the correct colors with a simple funtion.

func Daily_Ops(c *gin.Context) {

	c.JSON(200, gin.H{"Message": "Working on it"})
}

func Daily_Ops_TEST(c *gin.Context) {
	// purely for testing if the database connection works from the server
	data, err := getdata.TestConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
		})
		return
	}

	c.JSON(200, data)
}

func Transportation_post(c *gin.Context) {
	var data []models.LoadData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Adding data to the database")

	var err error
	err = database.AddOrderToDB(&data, "transportation")
	if err != nil {
		fmt.Printf("error %v \n", err)
	}
}

// ---------- Logisitics Handlers ----------
func Log_year_by_year(c *gin.Context) {

	//For now we're going to just get all data from the database
	//This data only includes finished weeks.
	fmt.Println("Getting the first data")
	// change to fectch data and use the new struct
	data, err := database.GetCachedData("logistics")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
			"Error":   err,
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

func Log_stacked_miles(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

func Logistics_post(c *gin.Context) {
	var data []models.LoadData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Adding data to the database")

	var err error
	err = database.AddOrderToDB(&data, "logistics")
	if err != nil {
		fmt.Printf("error %v \n", err)
	}
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

func Dispatch_post_WTDOT(c *gin.Context) {
	// receive data from the client
	var data []models.OTWTDStats
	if err := c.ShouldBindJSON(&data); err != nil {
		// if there is an error return a 400 status code
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// driver add the data to the database
	conn, err := database.Make_connection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// make sure the table exists
	query := `CREATE TABLE IF NOT EXISTS WTDOTStats (
		dispatcher TEXT,
		date DATE,
		startDate DATE,
		endDate DATE,
		TotalOrders INT,
		TotalStops INT,
		ServiceIncidents INT,
		OrderOnTime float,
		StopOnTime float,
		PRIMARY KEY (dispatcher, date)
);`

	_, err = conn.Exec(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, driver := range data {
		err = database.Add_OTWTDStats(conn, driver)
		if err != nil {
			fmt.Printf("error %v \n", err)
		}
	}

	c.JSON(200, gin.H{"Message": "Data received"})
}
