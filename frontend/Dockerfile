# Stage 1 - Build
FROM node:24-alpine AS build

WORKDIR /app

# Accept build argument
ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

# Copy package files for layer caching
COPY package*.json ./

# Install production dependencies only
RUN npm ci --no-audit --no-fund

# Copy source code
COPY . .

# Build application
RUN npm run build

# Stage 2 - Runtime (nginx)
FROM nginx:alpine

# Create non-root user
RUN addgroup -g 1001 app && adduser -S app -u 1001 -G app

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built assets from build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Create temp directories for nginx cache with correct permissions
RUN mkdir -p /tmp/client_temp /tmp/proxy_temp_path /tmp/fastcgi_temp /tmp/uwsgi_temp /tmp/scgi_temp && \
    chown -R app:app /tmp/client_temp /tmp/proxy_temp_path /tmp/fastcgi_temp /tmp/uwsgi_temp /tmp/scgi_temp /usr/share/nginx/html

# Expose port (changed to 8080 for non-root)
EXPOSE 8080

# Switch to non-root user
USER app

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
