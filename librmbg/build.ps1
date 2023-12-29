Write-Host "[Builder][Info] ***** librmbg builder *****" -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Extracting py library" -ForegroundColor DarkGreen
Copy-Item -Force -Path .\\pylib.dat -Destination .\\pylib.zip
Expand-Archive -Force -Path pylib.zip -DestinationPath .\\pylib
Remove-Item -Force -Path .\\pylib.zip


Write-Host "[Builder][Info] Moving venv package zipball..." -ForegroundColor DarkGreen
Copy-Item -Path ..\\venv-package\\py_runtime.zip -Destination .\\pycontext\\py_runtime.zip

Write-Host "[Builder][Info] Building dll core sdk..." -ForegroundColor DarkGreen
go build -buildmode=c-shared -o librmbg.dll

Remove-Item -Recurse -Force .\\pylib
