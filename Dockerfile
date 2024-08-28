FROM golang:1.23.0-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o telegram-bot cmd/link-saver-bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/telegram-bot .
COPY --from=builder /app/.env .

CMD ["./telegram-bot"]
