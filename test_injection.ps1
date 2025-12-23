# Test trade injection
Write-Host "Testing Mr.Alpha - Trade Injection" -ForegroundColor Cyan
Write-Host ""

$payload = @{
    wallet_name = "TestWhale"
    token = "ETH"
    value_usd = 25000
    conviction = "High"
    narrative = "Restaking"
} | ConvertTo-Json

Write-Host "Injecting trade..." -ForegroundColor Yellow
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/inject" -Method Post -Body $payload -ContentType "application/json"

Write-Host ""
Write-Host "✅ Trade injected!" -ForegroundColor Green
Write-Host ($response | ConvertTo-Json -Depth 3)
Write-Host ""
Write-Host "Check http://localhost:8080/api/alerts to see it" -ForegroundColor Cyan
