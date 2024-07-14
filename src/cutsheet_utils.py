from colorama import Fore


def has_boxes(text: str) -> bool:
    """If the sheet info is for a bowl or box lunch, return True."""

    lines = text.split("\n")

    box_count = 0
    for line in lines:
        # Line contains the word "Box" or "Bowl"
        if " box" in line.lower() or " bowl" in line.lower():
            box_count += 1

    return box_count > 2
