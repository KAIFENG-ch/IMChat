FROM golang:1.16 as build

MAINTAINER KAIFENG_ch 3184218074@qq.com

ENV GO111MODULE = on \
    CGO_ENABLE = 0 \
    GOOS = linux \
    GOARCH = amd64

WORKDIR /app
COPY . .
RUN go build -o IMChat .

WORKDIR /dist

EXPOSE 8000