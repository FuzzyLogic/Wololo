FROM debian:buster-slim
RUN apt update && apt install -y net-tools iputils-ping wget
ENTRYPOINT [ "/bin/bash", "-c", "while [ true ]; do echo hello; sleep 1; done"]

FROM golang:1.16 as builder
WORKDIR /go/src/github.com/FuzzyLogic/Wololo/test/device
COPY ./test/device/main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /root/
COPY --from=builder /go/src/github.com/FuzzyLogic/Wololo/test/device/main ./device
RUN chmod a+x ./device
#ENTRYPOINT [ "/bin/sh", "-c", "echo sleeping for 10s... ; sleep 10; ./device" ]  
ENTRYPOINT [ "./device" ]  
