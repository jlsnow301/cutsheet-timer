#!/bin/bash

# Get the full path to the file passed as the first argument
file_path="$1"

# Change to the directory where the script is located
cd "$(dirname "$0")"

# Activate the virtual environment
source venv/bin/activate

# Execute the Python script within the src directory with the full path to the dragged file
python src/main.py "$file_path"

read -p "Press [Enter] key to continue..."