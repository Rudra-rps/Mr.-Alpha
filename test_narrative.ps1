# Quick test of narrative detection
Write-Host "Testing Mr.Alpha - Narrative Detection" -ForegroundColor Cyan
Write-Host ""

cd "narrative_radar'"
python detect_narrative.py

Write-Host ""
Write-Host "Test complete!" -ForegroundColor Green
