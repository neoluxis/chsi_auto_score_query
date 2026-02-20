# Build stage
FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o chsi-query ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/chsi-query .
COPY --from=builder /app/.env.example .

EXPOSE 8080

CMD ["./chsi-query"]
