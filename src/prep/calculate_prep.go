package prep

import (
	"fmt"
	"math"
	"strconv"

	"github.com/fatih/color"
	"github.com/jlsnow301/cutsheet-timer/input"
	"github.com/jlsnow301/cutsheet-timer/utils"
)

// calculatePrepTime calculates the prep time based on the headcount and whether the order is for a box lunch.
func calculatePrepTime(count int, isBoxes bool) int {
	recommendedPrepTime := float64(count) / 2.6
	if isBoxes {
		recommendedPrepTime /= 2
	}
	totalPrepTime := math.Max(15, math.Ceil(recommendedPrepTime))
	return int(math.Ceil(totalPrepTime/5) * 5) // Round up to the nearest 5
}

// calculateAndConfirmPrepTime gets the prep time based on the headcount and whether the order is for a box lunch.
func CalculateAndConfirmPrepTime(size string, hasBoxes bool) (int, bool) {
	utils.PrintHeader("Prep Time")

	sizeInt := 15
	if size == "" {
		utils.PrintRed("Size not provided. Using default size of 15.")
	} else {
		var err error
		sizeInt, err = strconv.Atoi(size)
		if err != nil {
			utils.PrintRed(fmt.Sprintf("Invalid size value: %s. Using default size of 15.", size))
			sizeInt = 15
		}
	}

	utils.PrintStats(fmt.Sprintf("Headcount: %d. ", sizeInt))
	fmt.Println("Base formula is max(15, (count / 2.6))")

	isBoxes := hasBoxes && input.ConfirmBoxLunch()
	if isBoxes {
		utils.PrintGreen("The suggested prep time will be reduced.")
	}

	prepTime := calculatePrepTime(sizeInt, isBoxes)

	boxLunchText := ""
	if isBoxes {
		cyan := color.New(color.FgCyan).SprintFunc()
		boxLunchText = cyan(" (reduced for box lunch)")
	}
	utils.PrintStats(fmt.Sprintf("\nSuggested prep time: %d minutes (rounded)%s.", prepTime, boxLunchText))

	finalPrepTime := input.GetUserInput(prepTime, "Prep")

	return finalPrepTime, isBoxes
}
