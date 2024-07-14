from datetime import datetime, timedelta
from typing import Optional
from colorama import Fore
from input_utils import get_user_input
from text_utils import print_cyan, print_green, print_stats


def handle_rush(leave_time: datetime) -> datetime:
    """Adds additional time for rush hour."""

    print_cyan(
        "\nThis order is during rush hour. An additional 15 minutes is suggested."
    )
    additional_minutes = get_user_input("extra rush hour", 15)
    return leave_time - timedelta(minutes=additional_minutes)


def handle_suite(leave_time: datetime, suite_info: str) -> datetime:
    """Adds additional time for a suite."""

    print_cyan(f"\nThis order appears to be for a suite:\n{suite_info}")
    print("\nAn additional 10 minutes is suggested to park and navigate the building.")
    additional_minutes = get_user_input("extra building traversal", 10)
    return leave_time - timedelta(minutes=additional_minutes)


def calculate_leave_time(
    event_time: datetime,
    travel_time: int,
    is_boxes: bool,
    suite_info: Optional[str] = None,
) -> datetime:
    """Calculate the leave time based on the event time, travel time, and whether the order is for a suite."""

    if not isinstance(event_time, datetime):
        raise TypeError("event_time must be a datetime object")

    base_setup = 15 if is_boxes else 30
    print_stats(f"Base travel time: {travel_time} minutes")
    print_stats(
        f"Base setup time: {base_setup} minutes{Fore.CYAN + ' (reduced for box lunch)' if is_boxes else ''}"
    )

    leave_time = event_time - timedelta(minutes=base_setup + travel_time)

    if 16 <= event_time.hour <= 18:
        leave_time = handle_rush(leave_time)

    if suite_info:
        leave_time = handle_suite(leave_time, suite_info)

    leave_time = leave_time.replace(
        minute=leave_time.minute - (leave_time.minute % 5), second=0, microsecond=0
    )

    suggested_minutes = int((event_time - leave_time).total_seconds() / 60)
    print_stats(
        f"\nSuggested setup and travel time: {suggested_minutes} minutes (rounded)."
    )

    user_time = get_user_input("travel and setup", suggested_minutes)
    adjustment = suggested_minutes - user_time
    leave_time += timedelta(minutes=adjustment)

    if adjustment != 0:
        print_green(
            f"Adjusted leave time: Leaving {abs(adjustment)} minutes {'later' if adjustment > 0 else 'earlier'}."
        )

    return leave_time
