# Build the manager binary
FROM golang:1.15.0 as builder
RUN mkdir  /notification-gateway
RUN mkdir /notification-gateway/rate-limiter
WORKDIR /notification-gateway/rate-limiter
COPY common /notification-gateway/common/
# Copy the Go Modules manifests
COPY rate-limiter/go.mod go.mod
COPY rate-limiter/go.sum go.sum


# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
COPY rate-limiter .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on  go build -a -o rate-limiter main.go
ENTRYPOINT ["./rate-limiter"]