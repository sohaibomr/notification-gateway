# Build the manager binary
FROM golang:1.15.0 as builder
RUN mkdir  /notification-gateway
RUN mkdir /notification-gateway/gateway
WORKDIR /notification-gateway/gateway
COPY common /notification-gateway/common/
# Copy the Go Modules manifests
COPY gateway/go.mod go.mod
COPY gateway/go.sum go.sum


# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
COPY gateway .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on  go build -a -o gateway-server main.go
ENTRYPOINT ["./gateway-server"]