FROM golang:1.18

WORKDIR /app
COPY src go.mod main.go config.json ./
RUN go build -o main .
CMD ["/app/main", "-config", "config.json"]