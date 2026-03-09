FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS TARGETARCH

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /app/backuper cmd/api/main.go

FROM --platform=$TARGETPLATFORM alpine:latest

RUN apk add --no-cache \
    postgresql-client \
    ca-certificates \
    tzdata \
    file  # Добавляем file для диагностики

WORKDIR /app

COPY --from=builder /app/backuper .
RUN chmod +x /app/backuper

RUN mkdir -p /var/log && \
    touch /var/log/app.log && \
    chmod 666 /var/log/app.log

EXPOSE 8080

CMD ["./backuper"]