# librmbg

### 說明

這其實是一個小測試，目的是要將 [rembg](https://github.com/danielgatis/rembg) 這個 python 專案編譯成一個 dll  
並與其他語言進行整合  
_librmbg_ 庫其實很單純，就只有一個主要的函數，_RemoveBackground_ 用於去除圖片中的背景並返回去背後圖片的 buffer

### 建置環境

1. go 1.21.x
2. python 3.10.x
3. .NET sdk 8.x
4. node 20.x (須包含 vs build tools)

注:

1. python 環境中必須包含 pyinstaller, requests, packaging 以及 pyquery，沒有這些套件的話要手動安裝

### 建置

直接右鍵以 powershell 運行 build.ps1 即可，建置完成後檔案會在 _Demo_ 目錄內

### 運行環境

0. **重要: 需要安裝 Visual Studio ReDistribute** 直接裝最新版就好
1. .NET demo: 沒有其他需要裝的
2. node demo: 需安裝 node 20.x

### 運行

1. .NET demo: 運行 Demo/Net/RemoverDemo.exe
2. node demo: 工作目錄切換為 Demo/Node，在這邊運行 npm start

接著，直接提供需要去背的檔案路徑即可，若檔案非圖檔會報錯
