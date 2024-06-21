package _test

import (
	"fmt"
	"github.com/jadogg22/go-sharpGraphs/pkg/helpers"
	"testing"
)

func TestCountWorkingDays(t *testing.T) {
	// Test the CountWorkingDays function

	workingDays, currentDay := helpers.CountWorkingDays()
	fmt.Println("Working days: ", workingDays)
	fmt.Println("Current day: ", currentDay)

	// make sure the working days are less then 30
	if workingDays > 30 {
		t.Errorf("Working days are more than 30")
	}
}
