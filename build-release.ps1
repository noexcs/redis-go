$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-s -w -extldflags '-static'" -o .\target\redis-go_windows_amd64-release.exe

$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-s -w -extldflags '-static'" -o .\target\redis-go_linux_amd64-release

$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags "-s -w -extldflags '-static'" -o .\target\redis-go_darwin_amd64-release