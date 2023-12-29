Write-Host "[Builder][Info] ***** Nodejs demo builder *****" -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Installing dependencies..." -ForegroundColor DarkGreen
& npm i


Write-Host "[Builder][Info] Copying files..." -ForegroundColor DarkGreen
Copy-Item -Force -Path ..\\BackgroundRemoverNode\\dist\\* -Destination .\\
Copy-Item -Force -Path ..\\venv-package\\python3.dll -Destination .\\python3.dll
Copy-Item -Force -Path ..\\venv-package\\python310.dll -Destination .\\python310.dll
Copy-Item -Force -Path ..\\librmbg\\librmbg.dll -Destination .\\librmbg.dll