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

# Install runtime dependencies (shadow for usermod/groupmod, su-exec for dropping privileges)
RUN apk add --no-cache ca-certificates sqlite-libs shadow su-exec

# Create app user/group with default Unraid IDs (99:100)
RUN addgroup -g 1000 app && adduser -S app -u 1000 -G app && \
    mkdir -p /data && chown -R app:app /data

WORKDIR /app

# Copy binary and assets from build stage
COPY --from=build /app/organizr /app/organizr
COPY --from=build /app/assets /app/assets

# Copy entrypoint script
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Environment variables for user/group ID (Unraid defaults)
ENV PUID=99
ENV PGID=100

# Expose backend port
EXPOSE 8080

# Run entrypoint
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/app/organizr"]
