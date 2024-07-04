# shellcheck disable=SC2098
# shellcheck disable=SC2086
GOOS=linux GOARCH=amd64 go build -o redis-go_"${GOOS}"_"${GOARCH}";
GOOS=windows GOARCH=amd64 go build -o redis-go_"${GOOS}"_"${GOARCH}";
GOOS=darwin GOARCH=amd64 go build -o redis-go_${GOOS}_${GOARCH};