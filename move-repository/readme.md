
### build in linux
# Mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build filename.go

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build filename.go