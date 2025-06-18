FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o events-service .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/events-service .

EXPOSE 8082
CMD ["./events-service"]