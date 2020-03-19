@echo off
rm -rf build\windows
go build -ldflags "-H windowsgui" -o build\windows\crawler.exe
copy %cd%\default_config.json %cd%\build\windows\config.json
copy %cd%\crawler.exe.manifest %cd%\build\windows\crawler.exe.manifest
copy %cd%\icon.ico %cd%\build\windows\icon.ico
Xcopy /E /I %cd%\client\build %cd%\build\windows\public