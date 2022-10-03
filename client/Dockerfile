FROM golang:1.19-alpine3.16

ARG INDEX
ARG HOST
ARG PORT
ARG DELAY

WORKDIR /app
COPY go.mod main.go ./
RUN go build -o main .
CMD ["/app/main", "-host", "${HOST}", "-port", "${PORT}", "-delay", "${DELAY}", "-index", "${INDEX}"]