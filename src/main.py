from datetime import timedelta
import sys
import os
from colorama import init


from dotenv import load_dotenv

from calculate_leave import calculate_leave_time
from calculate_prep import calculate_and_confirm_prep_time
from cutsheet_utils import has_boxes
from header_utils import parse_header_info
from time_utils import get_base_travel_time, get_event_time
from text_utils import (
    extract_text_from_pdf,
    print_header,
    print_red,
    print_stars,
    print_stats,
    split_texts,
)

# Initialize colorama
init(autoreset=True)


def main():
    if len(sys.argv) < 2:
        print_red("Please drag a PDF file onto the script.")
        return 1

    # Construct the path to the .env file in the src directory
    env_path = os.path.join(os.path.dirname(__file__), ".env")

    if not os.path.exists(env_path):
        print_red("Please create a .env file in the project root.")
        return 1

    load_dotenv(env_path)

    pdf_path = sys.argv[1]
    pdf_text = extract_text_from_pdf(pdf_path)
    header_text, remaining_text = split_texts(pdf_text)
    header = parse_header_info(header_text)

    if header["destination"] is None:
        print_red("Unable to determine destination address.")
        return 1

    has_box = has_boxes(remaining_text)
    prep_time, is_boxes = calculate_and_confirm_prep_time(header["size"], has_box)

    print()
    print_header("Travel Time")
    print_stats(f"Site Address: {header['destination']}")

    origin = header["origin"]
    if origin is None:
        print_red("No origin specified.")
        return 1

    print_stats(f"Origin: {origin}\n")

    origin_address = os.getenv(f"{origin.upper()}_ADDRESS")
    if origin_address is None:
        print_red(f"Unknown origin: {header['origin']}")
        return 1

    travel_time = get_base_travel_time(origin_address, header["destination"])
    if travel_time is None:
        print_red("Unable to calculate travel time.")
        return 1

    event_time = get_event_time(header["event time"])
    if event_time is None:
        print_red("Invalid event time. Please use HH:MM AM/PM.")
        return 1

    leave_time = calculate_leave_time(
        event_time, travel_time, is_boxes, header["suite info"]
    )
    ready_by_time = leave_time - timedelta(minutes=prep_time)

    print()
    print_stars()
    print()
    print(f"Ready by: {ready_by_time.strftime('%I:%M %p')}")
    print(f"Leave by: {leave_time.strftime('%I:%M %p')}")
    print(f"Event time: {event_time.strftime('%I:%M %p')}\n")

    return 0


if __name__ == "__main__":
    try:
        sys.exit(main())
    except Exception as e:
        print_red(f"An error occurred: {e}")
        sys.exit(1)
