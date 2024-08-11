$GO111MODULE='on'
$GOPROXY='https://goproxy.cn,direct'
$GOOS='windows'
$GOARCH='amd64'
echo "build $GOOS $GOARCH cli..."
go build -mod=mod -ldflags="-s -w" -o DK4HDTrainer-cli.exe main-cli.go

