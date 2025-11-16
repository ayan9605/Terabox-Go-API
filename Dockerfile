FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Install tools
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    go install golang.org/x/tools/cmd/goimports@latest

COPY . .

# Auto-fix imports and format
RUN goimports -w . && go fmt ./...

# Tidy and download
RUN go mod tidy && go mod download

# Generate docs
RUN swag init

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
EXPOSE 8080
CMD ["./main"]
