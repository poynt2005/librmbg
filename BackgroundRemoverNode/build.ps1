Write-Host "[Builder][Info] ***** Nodejs binding builder *****" -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Installing dependencies..." -ForegroundColor DarkGreen
& npm i

Write-Host "[Builder][Info] Building by npm script..." -ForegroundColor DarkGreen
& npm run build