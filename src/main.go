package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/jlsnow301/cutsheet-timer/cutsheet"
	"github.com/jlsnow301/cutsheet-timer/header"
	"github.com/jlsnow301/cutsheet-timer/leave"
	"github.com/jlsnow301/cutsheet-timer/prep"
	timeutils "github.com/jlsnow301/cutsheet-timer/time"
	"github.com/jlsnow301/cutsheet-timer/travel"
	"github.com/jlsnow301/cutsheet-timer/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.PrintRed("Please drag a PDF file onto the script.")
		os.Exit(1)
	}

	envPath := filepath.Join(filepath.Dir(os.Args[0]), ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		utils.PrintRed("Please create a .env file in the project root.")
		os.Exit(1)
	}

	err := godotenv.Load(envPath)
	if err != nil {
		utils.PrintRed("Error loading .env file")
		os.Exit(1)
	}

	pdfPath := os.Args[1]
	pdfText, err := utils.ExtractTextFromPDF(pdfPath)
	if err != nil {
		utils.PrintRed(fmt.Sprintf("Error extracting text from PDF: %v", err))
		os.Exit(1)
	}

	headerText, remainingText := utils.SplitTexts(pdfText)
	headerInfo := header.ParseHeaderInfo(headerText)

	if headerInfo.Destination == "" {
		utils.PrintRed("Unable to determine destination address.")
		os.Exit(1)
	}

	hasBox := cutsheet.HasBoxes(remainingText)
	prepTime, isBoxes := prep.CalculateAndConfirmPrepTime(headerInfo.Size, hasBox)

	fmt.Println()
	utils.PrintHeader("Travel Time")
	utils.PrintStats(fmt.Sprintf("Site Address: %s", headerInfo.Destination))

	origin := headerInfo.Origin
	if origin == "" {
		utils.PrintRed("No origin specified.")
		os.Exit(1)
	}

	utils.PrintStats(fmt.Sprintf("Origin: %s\n", origin))

	originAddress := os.Getenv(strings.ToUpper(origin) + "_ADDRESS")
	if originAddress == "" {
		utils.PrintRed(fmt.Sprintf("Unknown origin: %s", headerInfo.Origin))
		os.Exit(1)
	}

	eventTime, err := timeutils.GetEventTime(headerInfo.EventDate, headerInfo.EventTime)
	if err != nil {
		utils.PrintRed("Invalid event time. Please use HH:MM AM/PM.")
		os.Exit(1)
	}

	travelTime, err := travel.GetBaseTravelTime(originAddress, headerInfo.Destination, eventTime)
	if err != nil {
		utils.PrintRed(fmt.Sprintf("Unable to calculate travel time: %v", err))
		os.Exit(1)
	}

	leaveTime := leave.CalculateLeaveTime(eventTime, travelTime, isBoxes, headerInfo.SuiteInfo)
	readyByTime := leaveTime.Add(-time.Minute * time.Duration(prepTime))

	fmt.Println()
	utils.PrintStars()
	fmt.Println()
	fmt.Printf("Ready by: %s\n", readyByTime.Format("03:04 PM"))
	fmt.Printf("Leave by: %s\n", leaveTime.Format("03:04 PM"))
	fmt.Printf("Event time: %s\n", eventTime.Format("03:04 PM"))
	fmt.Println()
	fmt.Println()
}
