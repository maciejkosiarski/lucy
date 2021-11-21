# syntax = docker/dockerfile:1-experimental
FROM golang:1.15.8-alpine AS container
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
RUN go mod init github.com/maciejkosiarski/lucy
RUN go mod tidy
RUN go mod download
COPY . .
RUN go build -o main .
WORKDIR /dist
RUN cp /build/main .
EXPOSE 3000
CMD ["/dist/main"]
