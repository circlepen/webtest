FROM golang:1.19-alpine3.16
WORKDIR /app

COPY . .
ENV GO111MODULE=on
ENV GIN_MODE=release
RUN go install github.com/cosmtrek/air@latest
RUN go mod download

EXPOSE 8000

ENTRYPOINT ["air", "-c", ".air.toml"]