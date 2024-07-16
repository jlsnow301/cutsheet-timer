package header

import (
	"regexp"
	"strings"
)

type HeaderInfo struct {
	Origin      string
	Destination string
	Size        string
	EventTime   string
	SuiteInfo   string
}

func normalizeAddress(address string) string {
	// Remove any occurrence of "Headcount" and everything after it
	headcountIndex := strings.Index(strings.ToLower(address), "headcount")
	if headcountIndex != -1 {
		address = address[:headcountIndex]
	}

	// Trim any trailing commas and spaces
	address = strings.TrimRight(address, ", ")

	// Add a space before any capitalized letter unless one already exists
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	address = re.ReplaceAllString(address, "$1 $2")

	// Ensure there's a space before the ZIP code if it exists
	zipRe := regexp.MustCompile(`(\D)(\d{5})$`)
	address = zipRe.ReplaceAllString(address, "$1 $2")

	// Check if the address contains a ZIP code
	hasZip := regexp.MustCompile(`\d{5}$`).MatchString(address)

	// Check if the address contains "Seattle" (case-insensitive)
	hasSeattle := strings.Contains(strings.ToLower(address), "seattle")

	// If there's no ZIP code and no "Seattle", add "Seattle"
	if !hasZip && !hasSeattle {
		address += " Seattle"
	}

	return address
}

func ParseHeaderInfo(content []string) HeaderInfo {
	info := HeaderInfo{}

	collectingAddress := false
	var addressParts []string

	for _, line := range content {
		line = strings.TrimSpace(line)

		if info.Origin == "" && (line == "Fremont" || line == "Eastlake") {
			info.Origin = line
		}

		if info.EventTime == "" && strings.HasPrefix(line, "Start Time:") {
			info.EventTime = strings.TrimSpace(strings.SplitN(line, "Start Time:", 2)[1])
		}

		if collectingAddress {
			if strings.Contains(line, "Headcount:") {
				collectingAddress = false
				info.Destination = normalizeAddress(strings.Join(addressParts, ", "))
			} else if line != "" {
				addressParts = append(addressParts, line)
			}
		}

		if info.Destination == "" {
			match := regexp.MustCompile(`Site Address:\s*(.*)`).FindStringSubmatch(line)
			if match != nil {
				addressParts = []string{match[1]}
				collectingAddress = true
				if strings.Contains(strings.ToLower(match[1]), "suite") {
					siteName := strings.SplitN(line, "Site Name:", 2)[1]
					info.SuiteInfo = strings.TrimSpace(siteName)
				}
			}
		}

		if strings.HasPrefix(line, "Site Name:") && strings.Contains(strings.ToLower(line), "suite") {
			siteName := strings.SplitN(line, "Site Name:", 2)[1]
			info.SuiteInfo = strings.TrimSpace(siteName)
		}

		if info.Size == "" {
			match := regexp.MustCompile(`Headcount:\s*(\d+)`).FindStringSubmatch(line)
			if match != nil {
				info.Size = match[1]
				collectingAddress = false
				if len(addressParts) > 0 {
					info.Destination = normalizeAddress(strings.Join(addressParts, ", "))
				}
			}
		}
	}

	// In case the address collection wasn't terminated by a Headcount line
	if collectingAddress && len(addressParts) > 0 {
		info.Destination = normalizeAddress(strings.Join(addressParts, ", "))
	}

	return info
}
