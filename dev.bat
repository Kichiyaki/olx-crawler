@echo off
SET DEFAULT_HANDLER=true
SET DISABLE_MENU=true
IF NOT EXIST "config.json" COPY "default_config.json" "config.json"
go run main.go