FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go mod download && go build -o api ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./api"]
