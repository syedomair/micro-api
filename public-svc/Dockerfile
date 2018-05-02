FROM golang:1.9.0 as builder

WORKDIR /go/src/github.com/syedomair/micro-api/public-svc

COPY src .
COPY proto proto
#RUN go get -u github.com/golang/dep/cmd/dep
#RUN dep init && dep ensure
COPY vendor vendor

RUN CGO_ENABLED=0 GOOS=linux go build -o public-svc -a -installsuffix cgo main.go repository.go handler.go nats.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/syedomair/micro-api/public-svc/public-svc /public-svc 

EXPOSE 8180

CMD ["./public-svc"]
