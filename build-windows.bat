@echo off
rm -rf build\windows
go build -ldflags "-H windowsgui" -o build\windows\crawler.exe
copy %cd%\default_config.json %cd%\build\windows\config.json
Xcopy /E /I %cd%\client\build %cd%\build\windows\public