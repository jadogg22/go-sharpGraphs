package database

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"database/sql"
	"fmt"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Load .env file for tests
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}

	// Search for the .env file in the current directory and up to 3 parent directories.
	for i := 0; i < 4; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			err := godotenv.Load(filepath.Join(dir, ".env"))
			if err != nil {
				log.Fatalf("Error loading .env file for tests: %v", err)
			}
			break
		}
		dir = filepath.Dir(dir)
	}

	// Initialize database connection for tests
	PG_HOST = os.Getenv("PG_HOST")
	PG_PORT = os.Getenv("PG_PORT")
	PG_USER = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	PG_DATABASE = os.Getenv("PG_DATABASE")

	DB, err = PG_Make_connection()
	if err != nil {
		log.Fatalf("Error connecting to the database for tests: %v", err)
	}

	// Run tests
	exitCode := m.Run()

	// Close database connection after tests
	if DB != nil {
		DB.Close()
	}

	os.Exit(exitCode)
}

func TestCheckSpecificWeeks(t *testing.T) {
	

	// Test specific weeks and years
	tests := []struct {
		year int
		week int
		expectError bool
	}{
		{year: 2020, week: 1, expectError: false}, // Assuming data exists for week 1, 2020
		{year: 2020, week: 34, expectError: false}, // Assuming data exists for week 34, 2020
		{year: 2023, week: 1, expectError: false},  // Assuming data exists for week 1, 2023
		{year: 2023, week: 34, expectError: false}, // Assuming data exists for week 34, 2023
		{year: 9999, week: 1, expectError: true},   // Assuming no data for future year
	}

	for _, tt := range tests {
		log.Printf("Running health check for Year=%d, Week=%d", tt.year, tt.week)
		err := CheckTransWeeklyRevHealthForSpecificWeek(tt.year, tt.week)
		if (err != nil) != tt.expectError {
			t.Errorf("CheckTransWeeklyRevHealthForSpecificWeek(Year=%d, Week=%d) got error = %v, want error = %v", tt.year, tt.week, err, tt.expectError)
		}
	}
}

// CheckTransWeeklyRevHealthForSpecificWeek checks if data exists for a specific week and year.
func CheckTransWeeklyRevHealthForSpecificWeek(year, week int) error {
	var totalRevenue sql.NullFloat64
	query := "SELECT totalrevenue FROM transportation_weekly_revenue WHERE year = $1 AND week = $2"
	
	err := DB.QueryRow(query, year, week).Scan(&totalRevenue)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no data found for Year=%d, Week=%d", year, week)
		} else {
			return fmt.Errorf("query error for Year=%d, Week=%d: %w", year, week, err)
		}
	}

	if !totalRevenue.Valid || totalRevenue.Float64 <= 0 {
		return fmt.Errorf("invalid or zero revenue for Year=%d, Week=%d (Value: %v)", year, week, totalRevenue)
	}
	return nil
}