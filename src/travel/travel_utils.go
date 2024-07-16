package travel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jlsnow301/cutsheet-timer/utils"
	"googlemaps.github.io/maps"
)

// getDirections fetches directions using Google Maps API.
func getDirections(origin, destination string, event *time.Time) *maps.Route {
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		color.Red(fmt.Sprintf("Error creating Google Maps client: %v", err))
		return nil
	}

	// If the time is in the past, just say "now"
	var eventTime string
	if event.Before(time.Now()) {
		eventTime = "now"
	} else {
		eventTime = strconv.FormatInt(event.Unix(), 10)
	}

	request := &maps.DirectionsRequest{
		Origin:        origin,
		Destination:   destination,
		DepartureTime: eventTime,
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

func getDoubleDistance(distanceText string) (string, error) {
	// Split the miles off the end of the string
	splitText := strings.SplitN(distanceText, " ", 2)
	if len(splitText) < 2 {
		return "", errors.New("no distance found")
	}

	distance, err := strconv.ParseFloat(splitText[0], 64)
	if err != nil {
		return "", err
	}

	double := 2.0
	// Multiply by 2 for round trips
	distance = distance * double

	return fmt.Sprintf("%.2f %s", distance, splitText[1]), nil
}

// GetBaseTravelTime gets the base travel time based on the origin and destination.
func GetBaseTravelTime(origin, destination string, event *time.Time) (int, error) {
	directionsResult := getDirections(origin, destination, event)
	durationMins, distanceText := parseDurationAndDistance(directionsResult)
	if durationMins == 0 || distanceText == "" {
		return 0, errors.New("no directions found")
	}

	roundTripMiles, err := getDoubleDistance(distanceText)
	if err != nil {
		return 0, err
	}

	utils.PrintStats(fmt.Sprintf("Total roundtrip mileage: %s\n", roundTripMiles))

	return durationMins, nil
}
