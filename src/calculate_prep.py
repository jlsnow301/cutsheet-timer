import math
from colorama import Fore
from input_utils import confirm_box_lunch, get_user_input
from text_utils import print_green, print_header, print_red, print_stats


def calculate_prep_time(count: int, is_boxes: bool) -> int:
    """Calculate the prep time based on the headcount and whether the sheet is likely to be a box lunch."""

    recommended_prep_time = count / 2.6
    if is_boxes:
        recommended_prep_time /= 2
    total_prep_time = max(15, math.ceil(recommended_prep_time))
    return math.ceil(total_prep_time / 5) * 5  # Round up to the nearest 5


def calculate_and_confirm_prep_time(
    size: str | None, has_boxes: bool
) -> tuple[int, bool]:
    """Get the prep time based on the headcount and whether the sheet is a box lunch."""

    print_header("Prep Time")

    if size is None:
        print_red("Size not provided. Using default size of 15.")
        size_int = 15
    else:
        try:
            size_int = max(15, int(size))  # Ensure minimum size of 15
        except ValueError:
            print_red(f"Invalid size value: {size}. Using default size of 15.")
            size_int = 15

    print_stats(f"Headcount: {size_int}. ")
    print("Base formula is max(15, (count / 2.6))")

    is_boxes = has_boxes and confirm_box_lunch()
    if is_boxes:
        print_green("The suggested prep time will be reduced.")

    prep_time = calculate_prep_time(size_int, is_boxes)

    box_lunch_text = Fore.CYAN + "(reduced for box lunch)" if is_boxes else ""
    print_stats(
        f"\nSuggested prep time: {prep_time} minutes (rounded). {box_lunch_text}"
    )

    final_prep_time = get_user_input("prep", prep_time)

    return final_prep_time, is_boxes
