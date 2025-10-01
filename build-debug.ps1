$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -tags debug -o .\target\redis-go_windows_amd64-debug.exe

$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -tags debug -o .\target\redis-go_linux_amd64-debug

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -tags debug -o .\target\redis-go_darwin_amd64-debug