package cutsheet

import (
	"strings"
)

// hasBoxes checks if the provided text indicates a bowl or box lunch.
func HasBoxes(text string) bool {
	lines := strings.Split(text, "\n")

	boxCount := 0
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), " box") || strings.Contains(strings.ToLower(line), " bowl") {
			boxCount++
		}
	}

	return boxCount > 2
}
