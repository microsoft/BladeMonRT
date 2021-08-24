@echo off
 
:: Usage: Run the BRT service.

go generate ./...
go run .
