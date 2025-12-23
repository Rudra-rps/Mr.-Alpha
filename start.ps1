# Mr.Alpha Quick Start Script for PowerShell

Write-Host "================================" -ForegroundColor Cyan
Write-Host "Mr.Alpha - Quick Start" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Starting services..." -ForegroundColor Yellow
Write-Host ""

# Start Alpha Radar (Python)
Write-Host "[1/3] Starting Alpha Radar (Python)..." -ForegroundColor Green
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd 'narrative_radar'''; python api.py"
Start-Sleep -Seconds 3

# Start Smart Money Tracker (Go)
Write-Host "[2/3] Starting Smart Money Tracker (Go)..." -ForegroundColor Green
Start-Process pwsh -ArgumentList "-NoExit", "-Command", "cd wallet_watcher; go run main.go"
Start-Sleep -Seconds 3

# Open Frontend
Write-Host "[3/3] Opening Frontend..." -ForegroundColor Green
Start-Process "frontend\index.html"

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "All services started!" -ForegroundColor Green
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Alpha Radar:        http://localhost:5000" -ForegroundColor White
Write-Host "Smart Money:        http://localhost:8080" -ForegroundColor White
Write-Host "Frontend:           Opening in browser..." -ForegroundColor White
Write-Host ""
Write-Host "Press Ctrl+C in each terminal window to stop services" -ForegroundColor Yellow
