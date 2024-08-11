# Set-ExecutionPolicy RemoteSigned
# Get-ExecutionPolicy -List

$GO111MODULE='on'
$GOPROXY='https://goproxy.cn,direct'
$GOOS='windows'
$GOARCH='amd64'
echo "build $GOOS $GOARCH debug..."
go build -mod=mod -ldflags="-s -w" -tags debug -o main.exe main.go
