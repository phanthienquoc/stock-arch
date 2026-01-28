# Trading System - Docker + Redis Stream Architecture

A modern, scalable trading system built with Docker, Redis Streams, and Python microservices.

## 🏗️ Architecture Overview

```
Signal/API/Bot
    ↓
Order Core (Fast Path) → Redis Stream → Consumers
    ↓
Order State | Risk Engine | TP/SL Engine | Notify Service
```

- **Order Core**: Fast order processing (O(1) pre-risk)
- **Order State**: Order persistence and exchange sync
- **Risk Engine**: Post-risk checks, kill switches
- **TP/SL Engine**: Take-profit and stop-loss automation
- **Notify Service**: Alerts via Telegram/Discord/Email
- **Redis Streams**: Asynchronous inter-service communication

---

## 📁 Project Structure

```
trading-system/
├── .github/
│   └── workflows/
│       ├── docker-release.yml     (Build & push images)
│       └── deploy.yml             (Deploy to VPS)
│
├── order-core/
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app/
│
├── order-state/
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app/
│
├── risk-engine/
├── tp-sl-engine/
├── notify-service/
│
├── docker-compose.local.yml       (DEV: Build + Redis)
├── docker-compose.yml        (PROD: Pull images only)
├── .env.local                     (Local dev variables)
├── .env                      (Production variables - managed by CI)
├── trading.service                (systemd unit file for VPS)
├── LOCAL_DEVELOPMENT.md           (Dev guide)
├── VPS_DEPLOYMENT.md              (VPS setup guide)
└── README.md
```

---

## 🚀 Quick Start

### Local Development

```bash
# Clone repo
git clone <repo-url>
cd trading-system

# Start all services
docker compose -f docker-compose.local.yml up -d

# View logs
docker compose -f docker-compose.local.yml logs -f

# Stop services
docker compose -f docker-compose.local.yml down
```

See [LOCAL_DEVELOPMENT.md](./LOCAL_DEVELOPMENT.md) for detailed guide.

### VPS Deployment

CI/CD automatically deploys when you push a git tag:

```bash
git tag v1.2.0
git push origin v1.2.0
```

Manual setup: See [VPS_DEPLOYMENT.md](./VPS_DEPLOYMENT.md)

---

## 🔄 Development Workflow

1. **Make changes** → Push to feature branch
2. **Test locally** → `docker compose -f docker-compose.local.yml up`
3. **Tag version** → `git tag v1.2.0`
4. **Push tag** → GitHub Actions builds & deploys
5. **VPS auto-updates** → No manual intervention needed

---

## 🐳 Docker Images

All images pushed to Docker Hub:

- `alexanderloveiris/order-core:latest`
- `alexanderloveiris/order-state:latest`
- `alexanderloveiris/risk-engine:latest`
- `alexanderloveiris/tp-sl-engine:latest`
- `alexanderloveiris/notify-service:latest`

---

## 📊 Redis Streams

Asynchronous communication between services:

| Stream | Consumer Group | Purpose |
|--------|---|---|
| `orders.fast` | `orders.db` | New orders |
| `orders.filled` | `orders.risk`, `orders.tp_sl` | Filled orders |
| `risk.events` | `orders.notify` | Risk alerts |
| `tp_sl.events` | `orders.notify` | TP/SL events |

---

## 📝 Configuration

### Local (`.env.local`)
```env
REDIS_HOST=redis
REDIS_PORT=6379
```

### Production (`.env`)
```env
ORDER_CORE_DIGEST=sha256:xxxxx     # Updated by CI
ORDER_STATE_DIGEST=sha256:xxxxx    # Updated by CI
RISK_ENGINE_DIGEST=sha256:xxxxx    # Updated by CI
TPSL_ENGINE_DIGEST=sha256:xxxxx    # Updated by CI
NOTIFY_DIGEST=sha256:xxxxx         # Updated by CI

REDIS_HOST=managed-redis.example.com
REDIS_PORT=6379
```

---

## 🔑 CI/CD Secrets (GitHub)

Add these to GitHub repository secrets:

```
DOCKER_USERNAME         # Docker Hub username
DOCKER_PASSWORD         # Docker Hub access token
VPS_HOST               # VPS IP or hostname
VPS_USER               # VPS username
VPS_SSH_KEY            # Private SSH key for VPS
```

---

## 📚 Documentation

- [Local Development Guide](./LOCAL_DEVELOPMENT.md) - How to develop locally
- [VPS Deployment Guide](./VPS_DEPLOYMENT.md) - How to set up production
- [System Architecture](./docs/trading_system_architecture.md) - Technical details
- [Deployment Blueprint](./docs/local_vs_vps_docker_blueprint.md) - Design decisions

---

## 🛠️ Development

### Add a new service

1. Create folder: `mkdir my-service/app`
2. Create: `my-service/Dockerfile`
3. Create: `my-service/requirements.txt`
4. Create: `my-service/app/main.py`
5. Add to `docker-compose.local.yml`

### View Redis Streams (local)

```bash
docker compose -f docker-compose.local.yml exec redis redis-cli
XINFO STREAM orders.fast
XREAD COUNT 10 STREAMS orders.fast 0
```

---

## 🐛 Troubleshooting

### Services won't start
```bash
docker compose -f docker-compose.local.yml logs
```

### Redis connection error
```bash
docker compose -f docker-compose.local.yml ps redis
```

### Need to reset everything
```bash
docker compose -f docker-compose.local.yml down -v
docker compose -f docker-compose.local.yml up -d
```

---

## 📄 License

MIT

---

## 👤 Author

Alexander Love Iris

---

**KEY PRINCIPLE**: 
> LOCAL = code + build + version  
> VPS = pull image + run  
> CI/CD = bridge between them
