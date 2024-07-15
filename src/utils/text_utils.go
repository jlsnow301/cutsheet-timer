package utils

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ledongthuc/pdf"
)

// extractTextFromPDF extracts text from a PDF file.
func ExtractTextFromPDF(pdfPath string) (string, error) {
	f, r, err := pdf.Open(pdfPath)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

// splitTexts splits the text into header and food service items.
func SplitTexts(text string) (headerText, remainingText string) {
	lines := strings.Split(text, "\n")
	splitIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "Food/Service Item" {
			splitIndex = i
			break
		}
	}
	if splitIndex != -1 {
		headerText = strings.Join(lines[:splitIndex+1], "\n")
		remainingText = strings.Join(lines[splitIndex+1:], "\n")
	} else {
		headerText = text
		remainingText = ""
	}
	return headerText, remainingText
}

// printStars prints a line of stars in yellow.
func PrintStars() {
	PrintYellow("**********************************************")
}

// printRed prints text in red.
func PrintRed(text string) {
	color.Red(text)
}

// printGreen prints text in green.
func PrintGreen(text string) {
	color.Green(text)
}

// printYellow prints text in yellow.
func PrintYellow(text string) {
	color.Yellow(text)
}

// printCyan prints text in cyan.
func PrintCyan(text string) {
	color.Cyan(text)
}

// printStats pretty prints statistics as yellow text until the colon.
func PrintStats(text string) {
	splitText := strings.Split(text, ":")
	PrintYellow(splitText[0] + ":" + splitText[1])
}

// printHeader prints a yellow star and the rest in white.
func PrintHeader(text string) {
	PrintStars()
	PrintYellow("*")
	PrintYellow("* " + text)
	PrintYellow("*")
	fmt.Println()
}
