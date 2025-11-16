# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install swag CLI and build tools
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Copy go.mod first
COPY go.mod ./

# Copy go.sum if it exists, otherwise create empty file
COPY go.su[m] ./
RUN test -f go.sum || touch go.sum

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
