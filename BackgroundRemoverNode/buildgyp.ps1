Write-Host "[Builder][Info] ***** NodeJs binding builder *****" -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Building Node-GYP" -ForegroundColor DarkGreen

if (Test-Path -Path .\\build) {
    Remove-Item -Force -Recurse .\\build
}

node-gyp configure
node-gyp build

if (!(Test-Path -Path .\\dist)) {
    New-Item -ItemType Directory -Path .\\dist
}

Move-Item -Force -Path .\\build\Release\\BackgroundRemover.node -Destination .\\dist\\BackgroundRemover.node
Remove-Item -Force -Recurse .\\build
