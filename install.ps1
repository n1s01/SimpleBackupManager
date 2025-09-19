# Backup Tool Installer

$downloadUrl = "https://github.com/n1s01/SimpleBackupManager/releases/download/v1.0.0/backup.exe"
$installDir = "$env:LOCALAPPDATA\BackupTool"
$executablePath = "$installDir\backup.exe"

Write-Host "Backup Tool Installer" -ForegroundColor White
Write-Host "Creating installation directory..." -ForegroundColor Gray

if (Test-Path $installDir) {
    Remove-Item "$installDir\*" -Force -ErrorAction SilentlyContinue
} else {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

Write-Host "Downloading backup.exe..." -ForegroundColor Gray

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $executablePath -ErrorAction Stop
    Write-Host "Download completed" -ForegroundColor Green
} catch {
    Write-Host "Download failed: $($_.Exception.Message)" -ForegroundColor Red
    pause
    exit 1
}

if (-not (Test-Path $executablePath)) {
    Write-Host "Installation failed" -ForegroundColor Red
    pause
    exit 1
}

$fileInfo = Get-Item $executablePath
$size = [math]::Round($fileInfo.Length / 1MB, 2)
Write-Host "File size: $size MB" -ForegroundColor Gray

Write-Host "Testing executable..." -ForegroundColor Gray

try {
    & $executablePath --help | Out-Null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Executable test passed" -ForegroundColor Green
    } else {
        Write-Host "Executable test warning" -ForegroundColor Yellow
    }
} catch {
    Write-Host "Executable test warning" -ForegroundColor Yellow
}

Write-Host "Adding to PATH..." -ForegroundColor Gray

$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$installDir*") {
    $newPath = $currentPath + ";" + $installDir
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "Added to PATH" -ForegroundColor Green
} else {
    Write-Host "Already added in PATH" -ForegroundColor Yellow
}
Write-Host "Installation complete" -ForegroundColor Green
Write-Host "Installed to: $installDir" -ForegroundColor Gray
Write-Host "Commands:" -ForegroundColor White
Write-Host "  backup init" -ForegroundColor Gray
Write-Host "  backup create" -ForegroundColor Gray
Write-Host "  backup list" -ForegroundColor Gray
Write-Host "  backup load" -ForegroundColor Gray
Write-Host "Restart ur PC to use backup command" -ForegroundColor Red

pause