FROM golang:alpine as builder
RUN apk add --no-cache git && \
    mkdir /build
ADD . /build/
WORKDIR /build 

RUN export GO111MODULE=on && \
    go mod init bitqubic.com/kubestress

RUN go get k8s.io/client-go@master
RUN go build -o kubestress .

FROM alpine
RUN apk add --no-cache stress-ng
COPY --from=builder /build/kubestress /app/
WORKDIR /app
CMD ["./kubestress"]
