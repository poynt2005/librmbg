Write-Host "[Builder][Info] ***** Rembg library builder *****" -ForegroundColor DarkGreen
Write-Host "[Builder][Info] Make sure you have dotnet sdk 8.x, go 1.21.x, nodejs(with npm) 20.x, Microsoft C++ Build Tools and python 3.10 installed on this machine to complete the building"  -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Initialize python venv" -ForegroundColor DarkGreen
& .\\init-venv.ps1

Write-Host "[Builder][Info] Build env packer" -ForegroundColor DarkGreen
Push-Location .\\EnvPacker
& .\\build.ps1
Pop-Location

Write-Host "[Builder][Info] Build core lib" -ForegroundColor DarkGreen
Push-Location .\\librmbg
& .\\build.ps1
Pop-Location

Write-Host "[Builder][Info] Build .NET demo" -ForegroundColor DarkGreen
Push-Location .\\RemoverDemo
& .\\build.ps1
Pop-Location

Write-Host "[Builder][Info] Build node native binding" -ForegroundColor DarkGreen
Push-Location .\\BackgroundRemoverNode
& .\\build.ps1
Pop-Location

Write-Host "[Builder][Info] Build node demo" -ForegroundColor DarkGreen
Push-Location .\\RemoverDemoNode
& .\\build.ps1
Pop-Location

Write-Host "[Builder][Info] Constructing demo..." -ForegroundColor DarkGreen
New-Item -Path .\\Demo\\Node -ItemType Directory -Force
New-Item -Path .\\Demo\\Net -ItemType Directory -Force


Write-Host "[Builder][Info] Constructing demo..." -ForegroundColor DarkGreen
Copy-Item -Path .\\RemoverDemoNode\\* -Destination .\\Demo\\Node -Recurse -Force
Remove-Item -Force -Path .\\Demo\\Node\\build.ps1
Remove-Item -Force -Path .\\Demo\\Node\\package-lock.json
Copy-Item -Path .\\RemoverDemo\\bin\\Release\\net8.0\\win-x64\\publish\\* -Destination .\\Demo\\Net -Recurse -Force

Write-Host "[Builder][Info] Done" -ForegroundColor DarkGreen
