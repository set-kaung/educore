# build app for linux azure machine

GOOS=linux GOARCH=amd64 go build -o app_linux ./cmd/server
echo "Go server built successfully"

scp -i ~/.ssh/SK-CW_key.pem -r app_linux web azureuser@20.41.112.19:~/educore

# restart educore app and secret fetching service
sudo systemctl reload educore.service
sudo systemctl reload educore-secrets.service
