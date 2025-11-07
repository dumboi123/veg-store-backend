```shell
# Make IDE has full permission on current project
sudo chown -R $USER:$USER path/to/project
sudo chmod +x *.sh

docker rmi benlun1201/veg-store-backend-builder
docker build -f docker/Dockerfile.builder -t learn-go/veg-store-backend-builder .

docker exec -it -uroot veg-store-backend bash
go run cmd/server/main.go
air -c .air.toml

docker exec -it -uroot veg-store-backend bash -c 'cd /app && air -c .air.toml'
docker exec -it -uroot veg-store-backend bash -c 'cd /app && go run cmd/server/main.go'

# Bonus
sudo ./scripts/start.sh
sudo ./scripts/restart.sh
sudo ./scripts/stop.sh
```

```shell
go mod init <module-name>
go run <file-name>.go
go clean -modcache
go mod tidy
go get -u

# Air - live reload for Go apps
go install github.com/air-verse/air@latest

# i18n
go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest
goi18n merge i18n/active.en.toml i18n/translate.vi.toml

# Swagger
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
swag init -g main.go --parseDependency --parseInternal --dir ./cmd/server,./internal/application/dto,./internal/restful -o docs
swag init -g ./cmd/server/main.go -o cmd/docs

go install -v github.com/go-delve/delve/cmd/dlv@latest
dlv debug --headless --listen=:40000 --api-version=2 --accept-multiclient --continue ./cmd/server
dlv exec ./app --headless --listen=:40000 --api-version=2 --accept-multiclient --continue

go clean -modcache
git config --global http.sslVerify false
export GOINSECURE=github.com,go.googlesource.com,golang.org,go.uber.org,google.golang.org,sigs.k8s.io
export GOSUMDB=off
export GOPROXY=direct
GODEBUG=x509ignoreCN=1 go env -w GOINSECURE=*.org
go mod tidy
go get -u=patch
```
