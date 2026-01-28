# Local Development Guide

## Quick Start - Local Development

### Prerequisites
- Docker & Docker Compose installed
- Git

### Setup

1. **Clone the repository**
```bash
cd ~/Documents/quocpt
git clone <repo-url> trading-new-arch
cd trading-new-arch
```

2. **Start services with docker-compose**
```bash
docker compose -f docker-compose.local.yml up -d
```

3. **Verify services are running**
```bash
docker compose -f docker-compose.local.yml ps
```

4. **View logs**
```bash
docker compose -f docker-compose.local.yml logs -f order-core
docker compose -f docker-compose.local.yml logs -f order-state
docker compose -f docker-compose.local.yml logs -f risk-engine
```

### Services Overview

- **Redis (6379)** - Message broker for inter-service communication
- **order-core** - Fast order processing with pre-risk checks
- **order-state** - Order persistence and state management
- **risk-engine** - Post-risk checks and position management
- **tp-sl-engine** - Take-profit and stop-loss automation
- **notify-service** - Alerts and notifications

### Useful Commands

```bash
# Stop all services
docker compose -f docker-compose.local.yml down

# Rebuild a specific service
docker compose -f docker-compose.local.yml build order-core

# Restart a service
docker compose -f docker-compose.local.yml restart order-core

# Execute command in running container
docker compose -f docker-compose.local.yml exec order-core /bin/bash

# View Redis streams
docker compose -f docker-compose.local.yml exec redis redis-cli
> XINFO STREAM orders.fast
> XREAD COUNT 10 STREAMS orders.fast 0
```

### Development Workflow

1. **Make code changes** to any service in `{service}/app/`
2. **Rebuild the service**: `docker compose -f docker-compose.local.yml build {service}`
3. **Restart the service**: `docker compose -f docker-compose.local.yml up -d {service}`
4. **Check logs**: `docker compose -f docker-compose.local.yml logs -f {service}`

### Testing Redis Streams

```bash
# Connect to Redis
docker compose -f docker-compose.local.yml exec redis redis-cli

# View all streams
XINFO STREAM orders.fast
XINFO STREAM orders.filled
XINFO STREAM risk.events
XINFO STREAM tp_sl.events

# Monitor stream in real-time
XREAD BLOCK 0 STREAMS orders.fast $
```

---

## Local Testing Checklist

- [ ] All containers start successfully: `docker compose -f docker-compose.local.yml ps`
- [ ] Redis is accessible: `docker compose -f docker-compose.local.yml exec redis redis-cli ping`
- [ ] Services connect to Redis: Check logs for connection messages
- [ ] Consumer groups are created: Check with `XINFO GROUPS orders.fast`
- [ ] Data flows through streams: Test with manual messages

---

## Troubleshooting

### Container fails to start
```bash
docker compose -f docker-compose.local.yml logs {service}
```

### Redis connection refused
- Check Redis is running: `docker compose -f docker-compose.local.yml ps redis`
- Verify port mapping: `docker port {container_id} 6379`

### Consumer group already exists
- Consumer groups are created automatically on first run
- To reset: `docker compose -f docker-compose.local.yml down -v` (removes volumes)

---

## Next Steps

1. Implement business logic in each service's `app/main.py`
2. Create integration tests using docker-compose
3. Set up GitHub Actions for CI/CD
4. Deploy to VPS using the production blueprint
