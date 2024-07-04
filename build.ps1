$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -o .\target\redis-go_windows_amd64.exe

$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o .\target\redis-go_linux_amd64

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -o .\target\redis-go_darwin_amd64