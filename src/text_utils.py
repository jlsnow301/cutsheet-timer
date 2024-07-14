import PyPDF2
from colorama import Fore


def extract_text_from_pdf(pdf_path):
    """Extract text from PDF."""

    with open(pdf_path, "rb") as file:
        reader = PyPDF2.PdfReader(file)
        text = ""
        for page in reader.pages:
            text += page.extract_text()

    return text


def split_texts(text: str) -> tuple:
    """Split the text into header and food service items."""

    lines = text.split("\n")
    split_index = None
    for i, line in enumerate(lines):
        if line.strip() == "Food/Service Item":
            split_index = i
            break
    if split_index is not None:
        header_text = "\n".join(lines[: split_index + 1])
        remaining_text = "\n".join(lines[split_index + 1 :])
    else:
        header_text = text
        remaining_text = ""

    return (header_text, remaining_text)


def print_stars():
    print_yellow("**********************************************")


def print_red(text):
    print(Fore.RED + text)


def print_green(text):
    print(Fore.GREEN + text)


def print_yellow(text):
    print(Fore.YELLOW + text)


def print_cyan(text):
    print(Fore.CYAN + text)


def print_stats(text):
    """Pretty prints statistics as yellow text until the colon."""
    split_text = text.split(":")
    print_yellow(split_text[0] + ":" + Fore.RESET + split_text[1])


def print_header(text):
    """Prints a yellow star and the rest in white."""
    print_stars()
    print_yellow("*")
    print_yellow("* " + Fore.RESET + text)
    print_yellow("*\n")
