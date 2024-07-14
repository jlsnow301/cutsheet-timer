import os
import re
import googlemaps
from datetime import datetime

from input_utils import prompt_for_event_time
from text_utils import print_red, print_stats


from datetime import datetime


def get_event_time(event_time: str | None) -> datetime | None:
    """Get the event time as a datetime object."""

    if event_time is None:
        print_red("No event time provided.")
        event_time = prompt_for_event_time()

    try:
        return datetime.strptime(event_time, "%I:%M %p")
    except ValueError:
        print_red(f"Invalid event time: {event_time}. Please re-enter.")
        event_time = prompt_for_event_time()
        try:
            return datetime.strptime(event_time, "%I:%M %p")
        except ValueError:
            return None


def get_directions(origin: str, destination: str):
    """Fetch directions using Google Maps API."""

    gmaps = googlemaps.Client(key=os.getenv("GOOGLE_MAPS_API_KEY"))
    try:
        return gmaps.directions(origin, destination, departure_time=datetime.now())  # type: ignore
    except Exception as e:
        print_red(f"Error fetching directions: {e}")
        return None


def parse_duration_and_distance(directions_result):
    """Extract duration and distance from the directions result."""

    if not directions_result:
        print_red("No directions found.")
        return None, None

    leg = directions_result[0]["legs"][0]
    duration_text = leg["duration"]["text"]
    distance_text = leg["distance"]["text"]

    return duration_text, distance_text


def extract_time_from_text(duration_text: str) -> int:
    """Extract total time in minutes from duration text."""

    hours = minutes = 0
    if hours_match := re.search(r"(\d+)\s*hour", duration_text):
        hours = int(hours_match.group(1))
    if minutes_match := re.search(r"(\d+)\s*min", duration_text):
        minutes = int(minutes_match.group(1))
    return hours * 60 + minutes


def extract_distance_and_unit(distance_text: str):
    """Extract distance and unit from distance text."""

    if distance_match := re.search(r"(\d+(?:\.\d+)?)\s*(km|mi)", distance_text):
        distance = float(distance_match.group(1))
        unit = distance_match.group(2)
        return distance, unit
    return None, None


def get_base_travel_time(origin: str, destination: str) -> int | None:
    """Get the base travel time based on the origin and destination."""

    directions_result = get_directions(origin, destination)
    duration_text, distance_text = parse_duration_and_distance(directions_result)
    if duration_text is None or distance_text is None:
        return None

    total_minutes = extract_time_from_text(duration_text)
    distance, unit = extract_distance_and_unit(distance_text)
    if distance and unit:
        roundtrip_distance = distance * 2
        print_stats(f"Total roundtrip mileage: {roundtrip_distance} {unit}\n")

    return total_minutes
