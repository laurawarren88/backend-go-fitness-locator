FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/air /usr/local/bin/air
COPY --from=builder /app/main /app/main
COPY --from=builder /app/.air.toml /app/.air.toml

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
