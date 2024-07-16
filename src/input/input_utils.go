package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jlsnow301/cutsheet-timer/utils"
)

// get_user_input prompts the user to input a valid number of minutes.
func GetUserInput(defaultMinutes int, reason string) int {
	blue := color.New(color.FgBlue).SprintFunc()

	fmt.Printf("Enter a different number or press %s to add suggested %d minutes.\n", blue("ENTER"), defaultMinutes)
	fmt.Println()

	var lowerReason string
	if reason != "" {
		lowerReason = strings.ToLower(reason)
	} else {
		lowerReason = "any additional"
	}

	fmt.Printf("Set %s time in minutes: ", lowerReason)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return 0
	}

	var outMsg string

	user_input := scanner.Text()
	if user_input == "" {
		if reason != "" {
			outMsg = fmt.Sprintf("%s set to %d minutes.", reason, defaultMinutes)
		} else {
			outMsg = fmt.Sprintf("%d minutes added.", defaultMinutes)
		}

		utils.PrintGreen(outMsg)
		return defaultMinutes
	}
	minutes, err := strconv.Atoi(user_input)
	if err != nil {
		utils.PrintRed("Invalid input. No additional time added.")
		return 0
	}

	if reason != "" {
		outMsg = fmt.Sprintf("%s set to %d minutes.", reason, minutes)
	} else {
		outMsg = fmt.Sprintf("%d minutes added.", minutes)
	}

	utils.PrintGreen(outMsg)
	return minutes

}

// confirm_box_lunch prompts the user to confirm if the sheet is a box lunch.
func ConfirmBoxLunch() bool {
	utils.PrintCyan("\nThis appears to be a box/bowl lunch.")
	fmt.Println("If this is correct, the suggested prep time will be reduced.")
	fmt.Print("Is this correct? Type n to revert, or ENTER to continue: ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		user_input := scanner.Text()
		return strings.ToLower(user_input) != "n"
	}
	return false
}

func PromptForEventTime() string {
	for {
		fmt.Print("Please enter the event time (HH:MM AM/PM): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			eventTime := scanner.Text()

			// Convert AM/PM to uppercase for parsing
			upperEventTime := strings.ToUpper(eventTime)

			// Try parsing with both "03:04 PM" and "3:04 PM" formats
			_, err1 := time.Parse("03:04 PM", upperEventTime)
			_, err2 := time.Parse("3:04 PM", upperEventTime)

			if err1 == nil || err2 == nil {
				return eventTime // Return the original input
			}
			utils.PrintRed("Invalid format. Please use HH:MM AM/PM.")
		}
	}
}
