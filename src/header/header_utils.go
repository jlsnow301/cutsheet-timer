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

// normalizeAddress normalizes an address by adding a space before the last 5 digits if not present.
func normalizeAddress(address string) string {
	re := regexp.MustCompile(`(\D)(\d{5})$`)
	return re.ReplaceAllString(address, `$1 $2`)
}

// parseHeaderInfo parses header information from the given text and returns a HeaderInfo struct.
func ParseHeaderInfo(text string) HeaderInfo {
	info := HeaderInfo{}
	lines := strings.Split(text, "\n")
	collectingAddress := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if info.Origin == "" && (line == "Fremont" || line == "Eastlake") {
			info.Origin = line
		}

		if info.EventTime == "" && strings.HasPrefix(line, "Start Time:") {
			info.EventTime = strings.TrimSpace(strings.SplitN(line, "Start Time:", 2)[1])
		}

		if collectingAddress {
			if line != "" && !strings.HasPrefix(line, "Headcount:") {
				normalizedAddress := normalizeAddress(line)
				if info.Destination != "" {
					info.Destination += ", " + normalizedAddress
				} else {
					info.Destination = normalizedAddress
				}
			} else {
				collectingAddress = false
			}
		}

		if info.Destination == "" {
			match := regexp.MustCompile(`Site Address:\s*(.*)`).FindStringSubmatch(line)
			if match != nil {
				normalizedAddress := normalizeAddress(match[1])
				info.Destination = normalizedAddress
				collectingAddress = true
				if strings.Contains(strings.ToLower(normalizedAddress), "suite") {
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
			}
		}
	}

	return info
}
