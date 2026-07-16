@echo off
echo Building COTW Save Backup Maker...

:: Ensure we are in the right directory
cd /d "%~dp0"

:: Package the Fyne app
echo Packaging app with Fyne CLI...
cd cmd\backup-maker
fyne package -os windows -icon ..\..\assets\icon.png
cd ..\..

:: Check if ISCC (Inno Setup Compiler) is in the PATH
where iscc >nul 2>nul
if %ERRORLEVEL% equ 0 (
    echo Compiling installer with Inno Setup...
    iscc installer.iss
) else (
    echo [WARNING] Inno Setup ISCC not found in PATH.
    echo Please install Inno Setup and compile 'installer.iss' manually to generate the installer.
)

echo Done!
