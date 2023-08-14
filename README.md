# chatgogogoproxy

## build script
Ubuntu

env GOOS=linux GOARCH=amd64 go build

Optimizing

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o echo-server server.go