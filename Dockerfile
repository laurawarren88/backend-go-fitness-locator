FROM golang:1.23-alpine
WORKDIR /app
RUN go install github.com/air-verse/air@latest
RUN apk add --no-cache bash 
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o main main.go
EXPOSE 8081
CMD ["air", "-c", ".air.toml"]