FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go get github.com/lib/pq

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client wget

RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz -O /tmp/migrate.tar.gz \
    && tar -xzf /tmp/migrate.tar.gz -C /usr/local/bin/ \
    && rm /tmp/migrate.tar.gz

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/.env .
COPY --from=builder /app/entrypoint.sh .

RUN chmod +x /root/main /root/entrypoint.sh /usr/local/bin/migrate

CMD ["sh", "/root/entrypoint.sh"]