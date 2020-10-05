current_dir=$(shell cd ../ && pwd)
check:
	go fmt ./...
	go vet ./...
	$(MAKE) tidy

.PHONY: tidy
tidy:
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on  go build -a -o gateway-server main.go
	
swagger-specs:
	docker run --rm -it -v $(current_dir):/go/src/rep/ quay.io/goswagger/swagger:v0.22.0 generate spec -m -w /go/src/rep/ -o /go/src/rep/swagger_gateway_specs.json