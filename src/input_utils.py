from datetime import datetime
from colorama import Fore

from text_utils import print_cyan, print_green, print_red


def get_user_input(reason: str, default_minutes: int) -> int:
    """Prompts the user to input a valid number of minutes."""

    print(
        f"Enter a different number or press {Fore.BLUE}ENTER{Fore.RESET} to add suggested {default_minutes} minutes.\n"
    )
    user_input = input(f"Set {reason} time in minutes: ").strip()
    if user_input == "":
        print_green(f"{default_minutes} minutes added.")
        return default_minutes
    try:
        minutes = int(user_input)
        print_green(f"\n{minutes} minutes added.")
        return minutes
    except ValueError:
        print_red("Invalid input. No additional time added.")
        return 0


def confirm_box_lunch() -> bool:
    """Prompts the user to confirm if the sheet is a box lunch."""

    print_cyan("\nThis appears to be a box/bowl lunch.")
    print("If this is correct, the suggested prep time will be reduced.")
    user_input = (
        input(
            f"Is this correct? Type "
            + Fore.RED
            + "n"
            + Fore.RESET
            + " to revert, or "
            + Fore.BLUE
            + "ENTER"
            + Fore.RESET
            + " to continue: "
        )
        .strip()
        .lower()
    )

    return user_input != "n"


def prompt_for_event_time() -> str:
    """Prompt the user to input a valid event time."""

    while True:
        event_time = input("Please enter the event time (HH:MM AM/PM): ")
        try:
            datetime.strptime(event_time, "%I:%M %p")
            return event_time
        except ValueError:
            print_red("Invalid format. Please use HH:MM AM/PM.")
