#!/bin/sh

GO111MODULE=on
GOPROXY=https://goproxy.cn,direct
GOOS=windows
GOARCH=amd64
echo build $GOOS $GOARCH release...
go build -ldflags="-s -w -H=windowsgui"  -o src/main.exe src/main.go
