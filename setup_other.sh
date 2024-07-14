#!/bin/bash

# Check if main.py exists in the current directory
if [ ! -f "./src/main.py" ]; then
    echo "Error: This script must be run from the primary directory where main.py exists."
    exit 1
fi

# Function to check if Python version is 3.9 or higher
check_python() {
    # Check for python command and get its version
    python_version=$(python3 --version 2>/dev/null | grep -oP 'Python \K\d+\.\d+')
    if [[ $? -eq 0 ]]; then
        # Compare the version to 3.9
        if [[ $(echo -e "3.9\n${python_version}" | sort -V | head -n1) == "3.9" ]]; then
            return 0
        fi
    fi
    return 1
}

# Install Python 3.12.4 if version is lower than 3.9
if ! check_python; then
    echo "Python version is lower than 3.9."
    if [ -f "python-installer.exe" ]; then
        echo "Using existing Python installer."
        chmod +x python-installer.exe
        ./python-installer.exe /quiet InstallAllUsers=1 PrependPath=1
        rm python-installer.exe
    else
        echo "Downloading and installing Python 3.12.4..."
        wget https://www.python.org/ftp/python/3.12.4/Python-3.12.4.tgz
        tar -xf Python-3.12.4.tgz
        cd Python-3.12.4
        ./configure --enable-optimizations
        make -j 8
        sudo make altinstall
        cd ..
        rm -rf Python-3.12.4 Python-3.12.4.tgz
    fi
else
    echo "Python version is 3.9 or higher."
fi

# Create a virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "Creating virtual environment..."
    python3.12 -m venv venv
else
    echo "Virtual environment already exists."
fi

# Activate the virtual environment
source venv/bin/activate

# Install required packages
pip install -r requirements.txt

# Copy the run_script.sh from /scripts to the current directory
cp ./scripts/run_script.sh .

# Move setup_windows.bat into the /src directory
mv setup_windows.bat ./src/

echo "Setup complete."
read -n 1 -s -r -p "Press any key to continue..."

# Move this script into the /src directory
mv -- "$0" ./src/