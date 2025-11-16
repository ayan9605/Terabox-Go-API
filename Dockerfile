# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install all build tools
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Copy entire project first
COPY . .

# Install swag and generate docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init

# Ensure all dependencies are available
RUN go mod tidy && \
    go mod download && \
    go mod verify

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
