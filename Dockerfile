FROM golang:1.17 as builder

MAINTAINER KAIFENG_ch 3184218074@qq.com

ENV GO111MODULE = on \
    GOPROXY = https://goproxy.cn,direct \
    CGO_ENABLE = 0 \
    GOOS = linux \
    GOARCH = amd64

WORKDIR /app
COPY . .
RUN go build -o IMChat .

ENV GIN_MODE = release
EXPOSE 8000

ENTRYPOINT ./main