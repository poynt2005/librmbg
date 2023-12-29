Write-Host "[Builder][Info] ***** DotNet Console Demo builder *****" -ForegroundColor DarkGreen

Write-Host "[Builder][Info] Building demo by .NET sdk" -ForegroundColor DarkGreen
& dotnet publish .\\RemoverDemo.csproj -c Release -r win-x64 --self-contained -p:PublishSingleFile=true  -p:DebugSymbols=false -p:DebugType=None