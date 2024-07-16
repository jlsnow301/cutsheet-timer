package travel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jlsnow301/cutsheet-timer/input"
	"github.com/jlsnow301/cutsheet-timer/utils"
	"googlemaps.github.io/maps"
)

// getEventTime gets the event time as a time.Time object.
func GetEventTime(eventTime string) (*time.Time, error) {
	if eventTime == "" {
		color.Red("No event time provided.")
		eventTime = input.PromptForEventTime()
	}

	// Convert eventTime to uppercase to handle lowercase am/pm
	eventTime = strings.ToUpper(eventTime)

	parsedTime, err := time.Parse("03:04 PM", eventTime)
	if err != nil {
		color.Red(fmt.Sprintf("Invalid event time: %s. Please re-enter.", eventTime))
		eventTime = input.PromptForEventTime()
		// Ensure the re-entered time is also converted to uppercase
		eventTime = strings.ToUpper(eventTime)
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
func parseDurationAndDistance(directionsResult *maps.Route) (int, string) {
	if directionsResult == nil {
		color.Red("No directions found.")
		return 0, ""
	}

	leg := directionsResult.Legs[0]
	duration := int(leg.DurationInTraffic.Minutes())
	distanceText := leg.Distance.HumanReadable

	return duration, distanceText
}

// GetBaseTravelTime gets the base travel time based on the origin and destination.
func GetBaseTravelTime(origin, destination string) (int, error) {
	directionsResult := getDirections(origin, destination)
	durationMins, distanceText := parseDurationAndDistance(directionsResult)
	if durationMins == 0 || distanceText == "" {
		return 0, errors.New("no directions found")
	}

	utils.PrintStats(fmt.Sprintf("Total roundtrip mileage: %s\n", distanceText))

	return durationMins, nil
}
