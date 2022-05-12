FROM golang:1.17 as builder

MAINTAINER KAIFENG_ch

ENV GO111MODULE = on \
    goproxy = https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS = Linux \
    GOARCH = amd64  \

WORKDIR /go/src
COPY . .

RUN go build -t imchat:1.0 .

WORKDIR /dist

EXPOSE 8000

CMD ['/dist/app']
