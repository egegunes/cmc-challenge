FROM golang:1.11 as builder

RUN go get github.com/lib/pq
RUN mkdir -p /go/src/github.com/egegunes/cmcgo
WORKDIR /go/src/github.com/egegunes/cmcgo
COPY main.go .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmcgo

FROM alpine
RUN apk add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/egegunes/cmcgo/cmcgo .

CMD ["./cmcgo"]
