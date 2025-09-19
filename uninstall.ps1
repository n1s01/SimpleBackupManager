# Backup Tool Uninstaller

Write-Host ""
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "              Backup Tool Uninstaller v1.0                     " -ForegroundColor Cyan  
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host ""

$installDir = "$env:LOCALAPPDATA\BackupTool"

Write-Host "[INFO] Checking installation..." -ForegroundColor Blue

if (-not (Test-Path $installDir)) {
    Write-Host "[INFO] Backup tool is not installed (directory not found)" -ForegroundColor Yellow
    pause
    exit 0
}

Write-Host "[INFO] Found installation at: $installDir" -ForegroundColor Blue

# Ask for confirmation
Write-Host ""
$confirmation = Read-Host "Are you sure you want to uninstall Backup Tool? (y/N)"

if ($confirmation -ne "y" -and $confirmation -ne "Y") {
    Write-Host "[INFO] Uninstallation cancelled" -ForegroundColor Yellow
    pause
    exit 0
}

# Remove from PATH
Write-Host ""
Write-Host "[INFO] Removing from PATH..." -ForegroundColor Blue

$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$installDir*") {
    $newPath = $currentPath -replace [regex]::Escape(";$installDir"), ""
    $newPath = $newPath -replace [regex]::Escape("$installDir;"), ""
    $newPath = $newPath -replace [regex]::Escape("$installDir"), ""
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "[SUCCESS] Removed from PATH" -ForegroundColor Green
} else {
    Write-Host "[INFO] Not found in PATH" -ForegroundColor Yellow
}

# Remove installation directory
Write-Host "[INFO] Removing installation files..." -ForegroundColor Blue

try {
    Remove-Item $installDir -Recurse -Force
    Write-Host "[SUCCESS] Installation files removed" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Failed to remove some files: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "                 Uninstallation Complete!                      " -ForegroundColor Green
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Backup tool has been removed from your system" -ForegroundColor White
Write-Host "You may need to restart your terminal for PATH changes to take effect" -ForegroundColor Yellow
Write-Host ""

pause