# 1) Сборка
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bot ./cmd/bot

# 2) Рантайм (минимальный образ без лишнего)
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/bot /app/bot
USER 65532:65532
ENTRYPOINT ["/app/bot"]
