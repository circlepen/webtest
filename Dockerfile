# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app

COPY . .
ENV GO111MODULE=on

RUN go mod download
RUN go mod tidy
RUN go build -o /app/main
RUN apk add curl

# Run Stage
FROM alpine:3.16
WORKDIR /app
COPY  --from=builder /app/main .

COPY app.env .

EXPOSE 8000

CMD ["/app/main"]
