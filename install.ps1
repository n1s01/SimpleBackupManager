# Backup Tool Installer
# Universal installer that builds and installs backup tool to PATH

Write-Host ""
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "           Backup Tool Universal Installer v1.0                " -ForegroundColor Cyan  
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
Write-Host "[INFO] Checking Go installation..." -ForegroundColor Blue
$goVersion = go version 2>$null
if (-not $goVersion) {
    Write-Host "[ERROR] Go is not installed or not in PATH!" -ForegroundColor Red
    Write-Host "Please install Go from: https://golang.org/dl/" -ForegroundColor Yellow
    pause
    exit 1
}
Write-Host "[SUCCESS] Go found: $goVersion" -ForegroundColor Green

# Check if current directory contains go.mod
if (-not (Test-Path "go.mod")) {
    Write-Host "[ERROR] go.mod not found in current directory!" -ForegroundColor Red
    Write-Host "Please run this installer from the backup-tool source directory" -ForegroundColor Yellow
    pause
    exit 1
}

Write-Host ""
Write-Host "[INFO] Building backup tool..." -ForegroundColor Blue

# Build the project
$buildResult = go build -ldflags="-s -w" -o backup.exe 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "[ERROR] Build failed!" -ForegroundColor Red
    Write-Host $buildResult -ForegroundColor Red
    pause
    exit 1
}

if (-not (Test-Path "backup.exe")) {
    Write-Host "[ERROR] backup.exe was not created!" -ForegroundColor Red
    pause
    exit 1
}

Write-Host "[SUCCESS] Build completed successfully!" -ForegroundColor Green

# Create installation directory
$installDir = "$env:LOCALAPPDATA\BackupTool"
Write-Host ""
Write-Host "[INFO] Creating installation directory: $installDir" -ForegroundColor Blue

if (Test-Path $installDir) {
    Write-Host "[INFO] Directory already exists, cleaning..." -ForegroundColor Yellow
    Remove-Item "$installDir\*" -Force -ErrorAction SilentlyContinue
} else {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

# Copy executable to installation directory
Write-Host "[INFO] Installing backup.exe to $installDir" -ForegroundColor Blue
Copy-Item "backup.exe" "$installDir\backup.exe" -Force

if (-not (Test-Path "$installDir\backup.exe")) {
    Write-Host "[ERROR] Failed to copy backup.exe to installation directory!" -ForegroundColor Red
    pause
    exit 1
}

Write-Host "[SUCCESS] Executable installed successfully!" -ForegroundColor Green

# Add to PATH
Write-Host ""
Write-Host "[INFO] Adding to PATH..." -ForegroundColor Blue

# Get current user PATH
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")

# Check if already in PATH
if ($currentPath -like "*$installDir*") {
    Write-Host "[INFO] Installation directory already in PATH" -ForegroundColor Yellow
} else {
    # Add to PATH
    $newPath = $currentPath + ";" + $installDir
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "[SUCCESS] Added to PATH successfully!" -ForegroundColor Green
    Write-Host "[INFO] You may need to restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
}

# Test installation
Write-Host ""
Write-Host "[INFO] Testing installation..." -ForegroundColor Blue

# Refresh PATH for current session
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User") + ";" + [Environment]::GetEnvironmentVariable("Path", "Machine")

$testResult = & "$installDir\backup.exe" --help 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Installation test passed!" -ForegroundColor Green
} else {
    Write-Host "[WARNING] Installation test failed, but files are installed" -ForegroundColor Yellow
}

# Clean up build file from source directory
if (Test-Path "backup.exe") {
    Remove-Item "backup.exe" -Force
    Write-Host "[INFO] Cleaned up build artifacts from source directory" -ForegroundColor Blue
}

Write-Host ""
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "                    Installation Complete!                     " -ForegroundColor Green
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Backup tool has been installed to: $installDir" -ForegroundColor White
Write-Host "Added to PATH for current user" -ForegroundColor White
Write-Host ""
Write-Host "Usage:" -ForegroundColor Yellow
Write-Host "  backup init          - Initialize project for backups" -ForegroundColor White
Write-Host "  backup create        - Create new backup" -ForegroundColor White
Write-Host "  backup list          - List all backups (interactive)" -ForegroundColor White
Write-Host "  backup load          - Load backup" -ForegroundColor White
Write-Host ""
Write-Host "You can now use 'backup' command from any directory!" -ForegroundColor Green
Write-Host ""

pause