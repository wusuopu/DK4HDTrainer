#!/bin/sh

GO111MODULE=on
GOPROXY=https://goproxy.cn,direct
GOOS=windows
GOARCH=amd64
echo build $GOOS $GOARCH deubg...
go build -mod=mod -ldflags="-s -w" -o main-cli.exe main-cli.go

