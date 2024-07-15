package travel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/jlsnow301/cutsheet-timer/input"
	"googlemaps.github.io/maps"
)

// getEventTime gets the event time as a time.Time object.
func GetEventTime(eventTime string) (*time.Time, error) {
	if eventTime == "" {
		color.Red("No event time provided.")
		eventTime = input.PromptForEventTime()
	}

	parsedTime, err := time.Parse("03:04 PM", eventTime)
	if err != nil {
		color.Red(fmt.Sprintf("Invalid event time: %s. Please re-enter.", eventTime))
		eventTime = input.PromptForEventTime()
		parsedTime, err = time.Parse("03:04 PM", eventTime)
		if err != nil {
			return nil, err
		}
	}
	return &parsedTime, nil
}

// getDirections fetches directions using Google Maps API.
func getDirections(origin, destination string) *maps.Route {
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		color.Red(fmt.Sprintf("Error creating Google Maps client: %v", err))
		return nil
	}

	request := &maps.DirectionsRequest{
		Origin:        origin,
		Destination:   destination,
		DepartureTime: "now",
	}

	routes, _, err := client.Directions(context.Background(), request)
	if err != nil {
		color.Red(fmt.Sprintf("Error fetching directions: %v", err))
		return nil
	}

	if len(routes) > 0 {
		return &routes[0]
	}
	return nil
}

// parseDurationAndDistance extracts duration and distance from the directions result.
func parseDurationAndDistance(directionsResult *maps.Route) (string, string) {
	if directionsResult == nil {
		color.Red("No directions found.")
		return "", ""
	}

	leg := directionsResult.Legs[0]
	durationText := leg.Duration.String()
	distanceText := leg.Distance.HumanReadable

	return durationText, distanceText
}

// extractTimeFromText extracts total time in minutes from duration text.
func extractTimeFromText(durationText string) int {
	hours, minutes := 0, 0
	hoursMatch := regexp.MustCompile(`(\d+)\s*hour`).FindStringSubmatch(durationText)
	if len(hoursMatch) > 1 {
		hours, _ = strconv.Atoi(hoursMatch[1])
	}
	minutesMatch := regexp.MustCompile(`(\d+)\s*min`).FindStringSubmatch(durationText)
	if len(minutesMatch) > 1 {
		minutes, _ = strconv.Atoi(minutesMatch[1])
	}
	return hours*60 + minutes
}

// extractDistanceAndUnit extracts distance and unit from distance text.
func extractDistanceAndUnit(distanceText string) (float64, string) {
	distanceMatch := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(km|mi)`).FindStringSubmatch(distanceText)
	if len(distanceMatch) > 2 {
		distance, _ := strconv.ParseFloat(distanceMatch[1], 64)
		unit := distanceMatch[2]
		return distance, unit
	}
	return 0, ""
}

// GetBaseTravelTime gets the base travel time based on the origin and destination.
func GetBaseTravelTime(origin, destination string) (*int, error) {
	directionsResult := getDirections(origin, destination)
	durationText, distanceText := parseDurationAndDistance(directionsResult)
	if durationText == "" || distanceText == "" {
		return nil, errors.New("no directions found")
	}

	totalMinutes := extractTimeFromText(durationText)
	distance, unit := extractDistanceAndUnit(distanceText)
	if distance > 0 && unit != "" {
		roundtripDistance := distance * 2
		color.Green(fmt.Sprintf("Total roundtrip mileage: %.2f %s\n", roundtripDistance, unit))
	}

	return &totalMinutes, nil
}
