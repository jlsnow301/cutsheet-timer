@echo off

REM Must be run from the directory containing main.py
if not exist src\main.py (
    echo This script must be run from the importer directory.
    exit /b
)

REM Function to check if Python 3.9 or higher is installed
python -c "import sys; exit(0 if sys.version_info >= (3,9) else 1)" >nul 2>&1
echo. 
if %errorlevel% equ 0 (
    echo Python 3.9 or newer is already installed.
) else (
    echo Python 3.9 or newer not found.
    if not exist python-installer.exe (
        echo Downloading Python 3.12.4 installer...
        powershell -Command "(New-Object System.Net.WebClient).DownloadFile('https://www.python.org/ftp/python/3.12.4/python-3.12.4-amd64.exe', 'python-installer.exe')"
    )
    echo Installing Python 3.12.4...
    start /wait python-installer.exe /quiet InstallAllUsers=1 PrependPath=1
    echo Deleting installer...
    del python-installer.exe
)

REM Create a virtual environment if it doesn't exist
echo.
if not exist venv (
    echo Creating virtual environment...
    python -m venv venv
) else (
    echo Virtual environment already exists.
)

REM Activate the virtual environment
call venv\Scripts\activate.bat

REM Install required packages
pip install -r src\requirements.txt

REM Copy the run script from /scripts to the current directory
copy scripts\run_script.bat .

REM Move setup_other.sh into the /src directory
move setup_other.sh src\

echo.
echo Setup complete.
echo.
pause

REM Move this script into the /src directory
move %0 src\
