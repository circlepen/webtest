# Build stage
FROM golang:1.19-alpine3.16
WORKDIR /app

COPY . .
ENV GO111MODULE=on
ENV GIN_MODE=release
RUN go install github.com/cosmtrek/air@latest
RUN go mod download

COPY app.env .

EXPOSE 8000

# ENTRYPOINT ["/app/main"]
# ENTRYPOINT [ "go", "run", "main.go" ]
ENTRYPOINT ["air", "-c", ".air.toml"]