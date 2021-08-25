@echo off
 
:: Usage: Build the BRT service to generate a .exe file.

go generate ./...
go build
