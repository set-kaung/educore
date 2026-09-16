GOOS=linux GOARCH=amd64 go build -o app_linux ./cmd/server
echo "Go server built successfully"
file app_linux
scp -r app_linux web bad-azure:~/educore 