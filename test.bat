@echo off
 
:: Usage: Generate all mocks and run tests.
 
go generate ./...
go test ./... -short