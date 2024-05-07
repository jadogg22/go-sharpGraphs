package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler function for	root route for debuging
func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "This came from the test handler",
	})
}

// ---------- Transportation Handlers ----------
func Trans_year_by_year(c *gin.Context) {
	// get date from system
	// conncet to database
	// pull all year by year data and and compair data
	db, err := Make_connection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error connecting to the database",
		})
		return
	}

	//For now we're going to just get all data from the database
	//This data only includes finished weeks.
	data, err := GetYearByYearData(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error getting data from the database",
		})
		return
	}

	fmt.Println("finished getting the first data")
	// Because its not very likly that we are at the end of the week
	// and we want to show the most recent data we need to check the
	// Transportation table and get the most recent data

	newData, err := GetNewestYearByYearData(db, data)
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

	// query db, and get the stacked miles

	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

func Trans_coded_revenue(c *gin.Context) {
	when := c.Param("when")
	fmt.Println("Getting coded revenue for ", when)

	conn, err := Make_connection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": "Error connecting to the database",
		})
		return
	}

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

	data, revenue, count, err2 := GetCodedRevenueData(conn, when)
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

// ---------- Logisitics Handlers ----------
func Log_year_by_year(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

func Log_stacked_miles(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Woring on it",
	})
}

// ---------- Dispatch Handlers ----------

func Dispach_week_to_date(c *gin.Context) {

	conn, _ := Make_connection()
	data, err := GetDispacherDataFromDB(conn)
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
	var data []DailyDriverData
	if err := c.ShouldBindJSON(&data); err != nil {
		// if there is an error return a 400 status code
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// driver add the data to the database
	conn, err := Make_connection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, driver := range data {
		err = Add_DailyDriverData(conn, driver)
		if err != nil {

			fmt.Printf("error %v \n", err)
		}

	}

	c.JSON(200, gin.H{"Message": "Data received"})
}
