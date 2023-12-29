Write-Host "[Builder][Info] ***** EnvPacker builder *****" -ForegroundColor DarkGreen


Write-Host "[Builder][Info] Building envpacker common tool..." -ForegroundColor DarkGreen
pyinstaller -F main.py
Copy-Item -Force -Path .\\dist\\main.exe -Destination .\\packer.exe

Remove-Item -Force -Recurse -Path dist
Remove-Item -Force -Recurse -Path build
Remove-Item -Force -Path .\\main.spec

Write-Host "[Builder][Info] Building python runtime package..." -ForegroundColor DarkGreen
& .\\packer.exe -v ..\\venv -o ..\\venv-package --add-exclude "cv2/config.py,cv2/config-3.py"