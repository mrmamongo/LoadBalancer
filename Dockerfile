FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/loadBalancer .


FROM scratch

WORKDIR /app
ADD config.json .
COPY --from=builder /app/loadBalancer /app/loadBalancer

CMD ["./loadBalancer", "-config", "config.json"]
