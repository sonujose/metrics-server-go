FROM golang:1.17.3-alpine3.14 AS builder

WORKDIR /go/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server-metrics-exporter .

FROM alpine:3.14 
RUN apk --no-cache add ca-certificates

WORKDIR /usr/app/

COPY --from=builder /go/src/app .

EXPOSE 8080

ENTRYPOINT ["./server-metrics-exporter"]