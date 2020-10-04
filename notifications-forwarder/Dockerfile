# Build the manager binary
FROM golang:1.15.0 as builder
RUN mkdir  /notification-gateway
RUN mkdir /notification-gateway/notifications-forwarder
WORKDIR /notification-gateway/notifications-forwarder
COPY common /notification-gateway/common/
# Copy the Go Modules manifests
COPY notifications-forwarder/go.mod go.mod
COPY notifications-forwarder/go.sum go.sum


# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
COPY notifications-forwarder .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on  go build -a -o notifications-forwarder main.go
ENTRYPOINT ["./notifications-forwarder"]