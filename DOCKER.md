# 🐳 Docker Quick Start Guide

## One-Command Deployment

```bash
docker-compose up --build
```

This will:
- Build Python narrative API
- Build Go wallet tracker
- Start nginx frontend proxy
- Wire everything together

**Access:** http://localhost:3000

## Services

| Service | Port | URL |
|---------|------|-----|
| Frontend | 3000 | http://localhost:3000 |
| Python API | 5000 | http://localhost:5000 |
| Go Backend | 8080 | http://localhost:8080 |

## Commands

### Start Services
```bash
docker-compose up
```

### Start in Background
```bash
docker-compose up -d
```

### View Logs
```bash
docker-compose logs -f
```

### Stop Services
```bash
docker-compose down
```

### Rebuild After Changes
```bash
docker-compose up --build
```

### Check Status
```bash
docker-compose ps
```

## Configuration

Environment variables in `docker-compose.yml`:
- `DEMO_MODE=true` - Use simulator for stable demo
- `TWITTER_BEARER_TOKEN` - Optional (uses demo data)

## Troubleshooting

**Port already in use:**
```bash
docker-compose down
# Kill local services first
taskkill /F /IM python.exe
taskkill /F /IM go.exe
```

**Rebuild from scratch:**
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up
```

**View container logs:**
```bash
docker logs mr-alpha-narrative
docker logs mr-alpha-wallet
docker logs mr-alpha-frontend
```

## Production Deployment

For production, update `docker-compose.yml`:
1. Set `DEMO_MODE=false`
2. Add actual `TWITTER_BEARER_TOKEN`
3. Configure Alchemy webhook to container URL
4. Use proper domain instead of localhost

## Architecture

```
                    Docker Network
┌─────────────────────────────────────────────┐
│                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │ Nginx    │  │ Python   │  │   Go     │ │
│  │ :3000    │──│  :5000   │  │  :8080   │ │
│  │ Frontend │  │ API      │  │ Backend  │ │
│  └──────────┘  └──────────┘  └──────────┘ │
│                                             │
└─────────────────────────────────────────────┘
        │
        ▼
   User Browser
   localhost:3000
```

Nginx proxies API requests to backend containers automatically.
