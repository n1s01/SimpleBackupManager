@echo off
setlocal enabledelayedexpansion

echo.
echo ================================================================
echo            Backup Tool Universal Installer v1.0                
echo ================================================================
echo.

REM Check if Go is installed
echo [INFO] Checking Go installation...
go version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH!
    echo Please install Go from: https://golang.org/dl/
    pause
    exit /b 1
)

for /f "tokens=*" %%i in ('go version 2^>nul') do set GO_VERSION=%%i
echo [SUCCESS] Go found: !GO_VERSION!

REM Check if go.mod exists
if not exist "go.mod" (
    echo [ERROR] go.mod not found in current directory!
    echo Please run this installer from the backup-tool source directory
    pause
    exit /b 1
)

echo.
echo [INFO] Building backup tool...

REM Build the project
go build -ldflags="-s -w" -o backup.exe
if errorlevel 1 (
    echo [ERROR] Build failed!
    pause
    exit /b 1
)

if not exist "backup.exe" (
    echo [ERROR] backup.exe was not created!
    pause
    exit /b 1
)

echo [SUCCESS] Build completed successfully!

REM Create installation directory
set INSTALL_DIR=%LOCALAPPDATA%\BackupTool
echo.
echo [INFO] Creating installation directory: !INSTALL_DIR!

if exist "!INSTALL_DIR!" (
    echo [INFO] Directory already exists, cleaning...
    del /q "!INSTALL_DIR!\*" >nul 2>&1
) else (
    mkdir "!INSTALL_DIR!"
)

REM Copy executable
echo [INFO] Installing backup.exe to !INSTALL_DIR!
copy "backup.exe" "!INSTALL_DIR!\backup.exe" >nul

if not exist "!INSTALL_DIR!\backup.exe" (
    echo [ERROR] Failed to copy backup.exe to installation directory!
    pause
    exit /b 1
)

echo [SUCCESS] Executable installed successfully!

REM Add to PATH using PowerShell
echo.
echo [INFO] Adding to PATH...

powershell -Command "$currentPath = [Environment]::GetEnvironmentVariable('Path', 'User'); if ($currentPath -notlike '*%LOCALAPPDATA%\BackupTool*') { $newPath = $currentPath + ';%LOCALAPPDATA%\BackupTool'; [Environment]::SetEnvironmentVariable('Path', $newPath, 'User'); Write-Host '[SUCCESS] Added to PATH successfully!' -ForegroundColor Green } else { Write-Host '[INFO] Already in PATH' -ForegroundColor Yellow }"

echo [INFO] You may need to restart your terminal for PATH changes to take effect

REM Test installation
echo.
echo [INFO] Testing installation...
"!INSTALL_DIR!\backup.exe" --help >nul 2>&1
if not errorlevel 1 (
    echo [SUCCESS] Installation test passed!
) else (
    echo [WARNING] Installation test failed, but files are installed
)

REM Clean up
if exist "backup.exe" (
    del "backup.exe"
    echo [INFO] Cleaned up build artifacts from source directory
)

echo.
echo ================================================================
echo                     Installation Complete!                     
echo ================================================================
echo.
echo Backup tool has been installed to: !INSTALL_DIR!
echo Added to PATH for current user
echo.
echo Usage:
echo   backup init          - Initialize project for backups
echo   backup create        - Create new backup
echo   backup list          - List all backups (interactive)
echo   backup load          - Load backup
echo.
echo You can now use 'backup' command from any directory!
echo.

pause