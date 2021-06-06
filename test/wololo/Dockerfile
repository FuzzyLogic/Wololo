FROM golang:1.16 as builder
WORKDIR /go/src/github.com/FuzzyLogic/Wololo
COPY ./ ./
RUN go get ./cmd/wololo
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/wololo

FROM alpine:latest  
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /root/
COPY --from=builder /go/src/github.com/FuzzyLogic/Wololo/wololo .
RUN chmod a+x ./wololo
RUN mkdir -p /etc/wololo
COPY test/wololo/config.json /etc/wololo/config.json
CMD ["./wololo"]  