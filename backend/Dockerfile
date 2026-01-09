# Stage 1 - Build
FROM golang:1.25-alpine AS build

# Install build dependencies required for CGO/SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go.mod and go.sum for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary with CGO enabled for SQLite support
# Use ldflags to strip debug symbols for smaller binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o organizr ./cmd/api

# Stage 2 - Runtime
FROM alpine:3.21

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite-libs

# Create non-root user
RUN addgroup -g 1001 app && adduser -S app -u 1001 -G app

WORKDIR /app

# Copy binary and assets from build stage
COPY --from=build --chown=app:app /app/organizr /app/organizr
COPY --from=build --chown=app:app /app/assets /app/assets

# Switch to non-root user
USER app

# Expose backend port
EXPOSE 8080

# Run the application
CMD ["/app/organizr"]
