@echo off
set "file_path=%~1"

:: Change to the directory where the batch script is located
cd /d "%~dp0"

:: Activate the virtual environment
call venv\Scripts\activate.bat

:: Execute the Python script within the src directory with the full path to the dragged file
python src\main.py "%file_path%"

pause