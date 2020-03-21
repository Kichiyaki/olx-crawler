@echo off
rm -rf build\windows
windres -O coff -o crawler.syso crawler.rc
go build -ldflags "-H windowsgui -X olx-crawler/config.Version=0.1.7" -o build\windows\crawler.exe
copy %cd%\default_config.json %cd%\build\windows\config.json
copy %cd%\icon.ico %cd%\build\windows\icon.ico
copy %cd%\crawler.exe.manifest %cd%\build\windows\crawler.exe.manifest
Xcopy /E /I %cd%\client\build %cd%\build\windows\public
Xcopy /E /I %cd%\i18n\locales %cd%\build\windows\i18n\locales