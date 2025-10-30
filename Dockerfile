FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build tools and dependencies, including PostgreSQL client
RUN apk add --no-cache git gcc musl-dev postgresql-client

# Copy source code
COPY . .

# Download dependencies
RUN go mod tidy && go mod verify

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Expose port
EXPOSE 8080

# Run the application with wait-for-db.sh
CMD ["./wait-for-db.sh", "library-db", "./main"]
