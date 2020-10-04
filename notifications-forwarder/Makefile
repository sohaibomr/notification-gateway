check:
	go fmt ./...
	go vet ./...
	$(MAKE) tidy

.PHONY: tidy
tidy:
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on  go build -a -o notifications-forwarder main.go