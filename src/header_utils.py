import re
from typing import Dict, Optional


def normalize_address(address: str) -> str:
    """Normalize an address by adding a space before the last 5 digits if not present."""

    return re.sub(r"(\D)(\d{5})$", r"\1 \2", address)


def parse_header_info(text: str) -> Dict[str, Optional[str]]:
    """Parse header information from the given text."""

    info: Dict[str, Optional[str]] = {
        "origin": None,
        "destination": None,
        "size": None,
        "event time": None,
        "suite info": None,
    }

    lines = text.split("\n")
    collecting_address = False

    for line in lines:
        line = line.strip()

        if not info["origin"] and line in ["Fremont", "Eastlake"]:
            info["origin"] = line

        if not info["event time"] and line.startswith("Start Time:"):
            info["event time"] = line.split("Start Time:", 1)[1].strip()

        if collecting_address:
            if line and not line.startswith("Headcount:"):
                info["destination"] = (
                    f"{info['destination']}, {normalize_address(line)}"
                    if info["destination"]
                    else normalize_address(line)
                )
            else:
                collecting_address = False

        if not info["destination"]:
            match = re.search(r"Site Address:\s*(.*)", line)
            if match:
                info["destination"] = normalize_address(match.group(1))
                collecting_address = True
                if "suite" in info["destination"].lower():
                    info["suite info"] = line.split("Site Name:", 1)[1].strip()

        if line.startswith("Site Name:") and "suite" in line.lower():
            info["suite info"] = line.split("Site Name:", 1)[1].strip()

        if not info["size"]:
            match = re.search(r"Headcount:\s*(\d+)", line)
            if match:
                info["size"] = match.group(1)
                collecting_address = False

    return info
