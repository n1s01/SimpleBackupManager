$installDir = "$env:LOCALAPPDATA\BackupTool"

Write-Host "Backup Tool Uninstaller" -ForegroundColor White
Write-Host "Removing installation..." -ForegroundColor Gray

if (-not (Test-Path $installDir)) {
    Write-Host "Not installed" -ForegroundColor Yellow
    pause
    exit 0
}

$confirmation = Read-Host "Remove Backup Tool? (y/N)"

if ($confirmation -ne "y" -and $confirmation -ne "Y") {
    Write-Host "Cancelled" -ForegroundColor Yellow
    pause
    exit 0
}

Write-Host "Removing from PATH..." -ForegroundColor Gray

$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$installDir*") {
    $newPath = $currentPath -replace [regex]::Escape(";$installDir"), ""
    $newPath = $newPath -replace [regex]::Escape("$installDir;"), ""
    $newPath = $newPath -replace [regex]::Escape("$installDir"), ""
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Removed from PATH" -ForegroundColor Green
} else {
    Write-Host "Not in PATH" -ForegroundColor Yellow
}

Write-Host "Removing files..." -ForegroundColor Gray

try {
    Remove-Item $installDir -Recurse -Force
    Write-Host "Files removed" -ForegroundColor Green
} catch {
    Write-Host "Remove failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "Uninstall complete" -ForegroundColor Green
Write-Host "Restart terminal for PATH changes" -ForegroundColor Yellow

pause