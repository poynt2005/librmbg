Write-Host "[Builder][Info] ***** python virtual env initiazlier *****" -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Initializing python venv..." -ForegroundColor DarkGreen
& python -m venv venv 

Write-Host "[Builder][Info] Entering python evnv..." -ForegroundColor DarkGreen
.\\venv\\Scripts\\Activate.ps1

Write-Host "[Builder][Info] Install requirements..." -ForegroundColor DarkGreen
& pip install -r req.txt

Write-Host "[Builder][Info] Exiting venv..." -ForegroundColor DarkGreen
deactivate