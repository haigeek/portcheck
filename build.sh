#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o portcheck-linux-amd64 main.go
GOOS=linux GOARCH=arm64 go build -o portcheck-linux-arm64 main.go
