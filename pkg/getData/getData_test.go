package getdata

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// TestMain runs before all tests in the package
func TestMain(m *testing.M) {
	// Load the .env file from the project root
	err := godotenv.Load("../../.env")
	if err != nil {
		// If this fails, the tests can't run correctly.
		// This is better than letting each test fail individually.
		panic("Failed to load .env file for testing: " + err.Error())
	}
	// Run the tests
	os.Exit(m.Run())
}

func TestGetLaneData(t *testing.T) {
	// Define a date range
	startDate, _ := time.Parse("2006-01-02", "2024-01-01")
	endDate, _ := time.Parse("2006-01-02", "2024-01-31")

	// Call the function
	data, err := GetLaneData(startDate, endDate)

	// Assert that there is no error and data is returned.
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// Print the number of rows returned
	fmt.Printf("Returned %d rows\n", len(data))

	// Print the first 5 rows for verification
	for i := 0; i < 5 && i < len(data); i++ {
		fmt.Printf("Row %d: %+v\n", i+1, data[i])
	}
}
