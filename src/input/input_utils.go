package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jlsnow301/cutsheet-timer/utils"
)

// get_user_input prompts the user to input a valid number of minutes.
func GetUserInput(reason string, default_minutes int) int {
	fmt.Printf("Enter a different number or press ENTER to add suggested %d minutes.\n", default_minutes)
	fmt.Printf("Set %s time in minutes: ", reason)

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		user_input := scanner.Text()
		if user_input == "" {
			utils.PrintGreen(fmt.Sprintf("%d minutes added.", default_minutes))
			return default_minutes
		}
		minutes, err := strconv.Atoi(user_input)
		if err != nil {
			utils.PrintRed("Invalid input. No additional time added.")
			return 0
		}
		utils.PrintGreen(fmt.Sprintf("\n%d minutes added.", minutes))
		return minutes
	}
	return 0
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

// prompt_for_event_time prompts the user to input a valid event time.
func PromptForEventTime() string {
	for {
		fmt.Print("Please enter the event time (HH:MM AM/PM): ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			event_time := scanner.Text()
			_, err := time.Parse("03:04 PM", event_time)
			if err == nil {
				return event_time
			}
			utils.PrintRed("Invalid format. Please use HH:MM AM/PM.")
		}
	}
}
