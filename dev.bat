@echo off
SET MODE=development
IF NOT EXIST "config.json" COPY "default_config.json" "config.json"
go run main.go