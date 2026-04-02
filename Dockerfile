FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy depending
COPY go.mod go.sum ./
RUN go mod download

# Copy the code and build the application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o todo-app ./cmd/main.go

# The final image
FROM alpine:3.20

WORKDIR /app

# Installing CA certificates and the PostgreSQL client
RUN apk add --no-cache ca-certificates && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup

# Copy the binary from builder
COPY --from=builder /app/todo-app .

EXPOSE 8000

CMD ["./todo-app"]