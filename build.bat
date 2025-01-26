  SET CGO_ENABLED=0
  SET GOOS=linux
  SET GOARCH=amd64
  go build -ldflags="-s -w" -o go_process_manager cmd/go_process_manager/main.go