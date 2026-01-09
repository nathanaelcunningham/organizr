# Docker Foundation - Build & Run Notes

## Image Sizes

- Backend: ~32MB (target: <50MB) ✅
- Frontend: ~54MB (target: <30MB, acceptable for nginx+assets)

## Building Images

### Backend
```bash
docker build -t organizr-backend ./backend
```

### Frontend
```bash
docker build -t organizr-frontend ./frontend
```

## Running Containers (Basic)

### Backend
```bash
docker run -d -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name organizr-backend \
  organizr-backend
```

Note: Backend requires volume mount for SQLite database persistence.

### Frontend
```bash
docker run -d -p 8081:8080 \
  --name organizr-frontend \
  organizr-frontend
```

Note: Frontend nginx listens on port 8080 (non-privileged port for non-root user).

## Important Notes

- Both containers run as non-root user (uid/gid 1001) for security
- Backend uses CGO_ENABLED=1 for SQLite support
- Frontend uses port 8080 (not 80) due to non-root user restrictions
- Full orchestration with environment variables and networking will be configured in Phase 17 (Docker Compose)
- Production deployment configuration (reverse proxy, SSL, etc.) will be in Phase 18-19

## Next Steps

See Phase 17 for docker-compose configuration that orchestrates both services.
