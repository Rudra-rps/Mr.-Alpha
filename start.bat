@echo off
echo ================================
echo Mr.Alpha - Quick Start
echo ================================
echo.

echo Starting services...
echo.

echo [1/3] Starting Alpha Radar (Python)...
cd "narrative_radar'"
start "Alpha Radar API" cmd /k "python api.py"
timeout /t 3 /nobreak >nul
cd ..

echo [2/3] Starting Smart Money Tracker (Go)...
cd wallet_watcher
start "Smart Money Tracker" cmd /k "go run main.go"
timeout /t 3 /nobreak >nul
cd ..

echo [3/3] Opening Frontend...
cd frontend
start "" "index.html"
cd ..

echo.
echo ================================
echo All services started!
echo ================================
echo.
echo Alpha Radar:        http://localhost:5000
echo Smart Money:        http://localhost:8080
echo Frontend:           Opening in browser...
echo.
echo Press any key to stop all services...
pause >nul

taskkill /FI "WINDOWTITLE eq Alpha Radar API*" /F /T >nul 2>&1
taskkill /FI "WINDOWTITLE eq Smart Money Tracker*" /F /T >nul 2>&1

echo Services stopped.
