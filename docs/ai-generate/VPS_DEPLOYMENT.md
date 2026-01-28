# VPS Deployment Guide

This guide covers the **production setup** on VPS/EC2. Source code is NOT included on VPS.

## VPS Setup (One-time)

### 1. Install Docker

```bash
curl -fsSL https://get.docker.com | sh
```

### 2. Enable Docker daemon

```bash
sudo systemctl enable docker
sudo systemctl start docker
```

### 3. Allow non-root user to run Docker

```bash
sudo usermod -aG docker $USER
exit  # logout and login to apply group changes
```

### 4. Create trading directory

```bash
sudo mkdir -p /opt/trading
sudo chown -R $USER:$USER /opt/trading
cd /opt/trading
```

### 5. Verify Docker is working

```bash
docker ps
docker compose version
```

---

## File Structure on VPS

```
/opt/trading/
├── docker-compose.yml  (copied by CI/CD)
├── .env                 (copied by CI/CD)
└── logs/                      (created by services)
    ├── order-core.log
    ├── order-state.log
    ├── risk-engine.log
    ├── tp-sl-engine.log
    └── notify-service.log
```

**Important**: Source code is NOT on VPS. Only compose and env files.

---

## Initial .env

After copying files from CI/CD, `.env` contains:

```env
ORDER_CORE_DIGEST=sha256:xxxxx
ORDER_STATE_DIGEST=sha256:xxxxx
RISK_ENGINE_DIGEST=sha256:xxxxx
TPSL_ENGINE_DIGEST=sha256:xxxxx
NOTIFY_DIGEST=sha256:xxxxx

REDIS_HOST=managed-redis.example.com
REDIS_PORT=6379
```

---

## systemd Auto-start Service

### 1. Create service file

```bash
sudo cp trading.service /etc/systemd/system/
```

Or create manually:

```bash
sudo vim /etc/systemd/system/trading.service
```

With content from `trading.service` in repo root.

### 2. Enable and start

```bash
sudo systemctl daemon-reload
sudo systemctl enable trading
sudo systemctl start trading
```

### 3. Verify

```bash
sudo systemctl status trading
sudo systemctl is-enabled trading
```

### 4. View logs

```bash
sudo journalctl -u trading -f
```

---

## Deployment via CI/CD

When you push a **git tag** (e.g., `v1.2.0`):

1. **GitHub Actions** builds Docker images
2. **Pushes to Docker Hub** with digest
3. **Updates `.env`** with new digests
4. **Pulls latest images** on VPS
5. **Restarts services** with `docker compose up -d`

Manual VPS intervention **NOT needed**.

---

## Manual Operations on VPS

### View running services

```bash
cd /opt/trading
docker compose -f docker-compose.yml ps
```

### Check service logs

```bash
docker compose -f docker-compose.yml logs -f order-core
```

### Stop all services

```bash
docker compose -f docker-compose.yml down
```

### Manually start (if using systemd, not needed)

```bash
cd /opt/trading
docker compose -f docker-compose.yml up -d
```

### Rollback to previous version

1. Edit `.env` with old digest
2. Run: `docker compose -f docker-compose.yml up -d`

---

## Monitoring & Alerting

Check service status regularly:

```bash
# Health check
docker compose -f docker-compose.yml ps

# Check disk usage
docker system df

# Cleanup old images (optional)
docker image prune -a --filter "until=168h"
```

---

## Troubleshooting

### Services won't start

```bash
docker compose -f docker-compose.yml logs
```

### Image pull fails (authentication)

```bash
docker login  # with Docker Hub credentials
docker compose -f docker-compose.yml pull
```

### Out of disk space

```bash
docker system prune -a
```

---

## Important: When to SSH VPS

| Scenario | Action |
|----------|--------|
| Deploy new version | ❌ Don't SSH - Let CI/CD handle |
| Rollback version | ❌ Update `.env` from local, push |
| Change compose file | ✅ SSH + copy new file manually |
| Troubleshoot issues | ✅ SSH to check logs & status |
| Update infra | ✅ SSH for infrastructure changes |

---

## References

- [Docker Installation](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/)
- [systemd Service Unit](https://systemd.io/)
