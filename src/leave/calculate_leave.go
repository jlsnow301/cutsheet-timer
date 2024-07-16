package leave

import (
	"fmt"
	"time"

	"github.com/jlsnow301/cutsheet-timer/input"
	"github.com/jlsnow301/cutsheet-timer/utils"
)

// HandleRush adds additional time for rush hour to the leave time.
func handleRush(leaveTime time.Time) time.Time {
	utils.PrintCyan("\nThis order is during rush hour. An additional 15 minutes is suggested.")
	additionalMinutes := input.GetUserInput("extra rush hour", 15)
	return leaveTime.Add(-time.Minute * time.Duration(additionalMinutes))
}

// HandleSuite adds additional time for a suite to the leave time.
func handleSuite(leaveTime time.Time, suiteInfo string) time.Time {
	utils.PrintCyan("\nThis order appears to be for a suite:")
	fmt.Println(suiteInfo)
	fmt.Println("\nAn additional 10 minutes is suggested to park and navigate the building.")
	additionalMinutes := input.GetUserInput("extra building traversal", 10)
	return leaveTime.Add(-time.Minute * time.Duration(additionalMinutes))
}

// CalculateLeaveTime calculates the leave time based on the event time, travel time, and whether the order is for a suite.
func CalculateLeaveTime(eventTime *time.Time, travelTime int, isBoxes bool, suiteInfo string) time.Time {
	baseSetup := 30
	if isBoxes {
		baseSetup = 15
	}

	utils.PrintStats(fmt.Sprintf("Base travel time: %d minutes", travelTime))
	utils.PrintStats(fmt.Sprintf("Base setup time: %d minutes", baseSetup))

	leaveTime := eventTime.Add(-time.Minute * time.Duration(baseSetup+travelTime))

	if eventTime.Hour() >= 16 && eventTime.Hour() <= 18 {
		leaveTime = handleRush(leaveTime)
	}
	if suiteInfo != "" {
		leaveTime = handleSuite(leaveTime, suiteInfo)
	}

	// Round down to the nearest 5 minutes
	leaveTime = leaveTime.Add(-time.Minute * time.Duration(leaveTime.Minute()%5))
	leaveTime = leaveTime.Add(-time.Second * time.Duration(leaveTime.Second()))
	leaveTime = leaveTime.Add(-time.Nanosecond * time.Duration(leaveTime.Nanosecond()))

	suggestedMinutes := int(eventTime.Sub(leaveTime).Minutes())
	utils.PrintStats(fmt.Sprintf("\nSuggested setup and travel time: %d minutes (rounded).", suggestedMinutes))

	userTime := input.GetUserInput("travel and setup", suggestedMinutes)
	adjustment := suggestedMinutes - userTime
	leaveTime = leaveTime.Add(time.Minute * time.Duration(adjustment))

	if adjustment != 0 {
		if adjustment > 0 {
			utils.PrintGreen(fmt.Sprintf("Adjusted leave time: Leaving %d minutes later.", adjustment))
		} else {
			utils.PrintGreen(fmt.Sprintf("Adjusted leave time: Leaving %d minutes earlier.", -adjustment))
		}
	}

	return leaveTime
}
