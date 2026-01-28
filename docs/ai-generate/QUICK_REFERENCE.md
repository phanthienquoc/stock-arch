# Quick Reference - Common Commands

## 🚀 Local Development

### Start all services
```bash
docker compose -f docker-compose.local.yml up -d
```

### Stop all services
```bash
docker compose -f docker-compose.local.yml down
```

### View logs (all)
```bash
docker compose -f docker-compose.local.yml logs -f
```

### View logs (specific service)
```bash
docker compose -f docker-compose.local.yml logs -f order-core
```

### Rebuild a service
```bash
docker compose -f docker-compose.local.yml build order-core
```

### Restart a service
```bash
docker compose -f docker-compose.local.yml restart order-core
```

### Execute command in container
```bash
docker compose -f docker-compose.local.yml exec order-core /bin/bash
```

---

## 📊 Redis Streams

### Connect to Redis
```bash
docker compose -f docker-compose.local.yml exec redis redis-cli
```

### View all streams
```bash
KEYS orders.*
KEYS risk.*
KEYS tp_sl.*
```

### View stream info
```bash
XINFO STREAM orders.fast
XINFO STREAM orders.filled
XINFO STREAM risk.events
XINFO STREAM tp_sl.events
```

### View consumer groups
```bash
XINFO GROUPS orders.fast
```

### Read recent messages
```bash
XREAD COUNT 5 STREAMS orders.fast 0
```

### Monitor stream in real-time
```bash
XREAD BLOCK 0 STREAMS orders.fast $
```

### Clear a stream (careful!)
```bash
DEL orders.fast
```

---

## 📦 Build & Push

### Build local images
```bash
docker compose -f docker-compose.local.yml build
```

### Tag for release (triggers CI/CD)
```bash
git tag v1.2.0
git push origin v1.2.0
```

### Manually build and push specific image
```bash
docker build -t alexanderloveiris/order-core:v1.2.0 ./order-core
docker push alexanderloveiris/order-core:v1.2.0
```

---

## 🔧 Container Maintenance

### List all containers
```bash
docker ps -a
```

### View container stats
```bash
docker stats
```

### Remove unused images
```bash
docker image prune -a
```

### Remove unused volumes
```bash
docker volume prune
```

### View image layers
```bash
docker history alexanderloveiris/order-core:latest
```

---

## 🖥️ VPS Commands

### Check service status
```bash
sudo systemctl status trading
```

### View systemd logs
```bash
sudo journalctl -u trading -f
```

### Restart service
```bash
sudo systemctl restart trading
```

### View running containers
```bash
cd /app/trading
docker compose -f docker-compose.yml ps
```

### Check service logs
```bash
cd /app/trading
docker compose -f docker-compose.yml logs -f order-core
```

### Manual restart (if needed)
```bash
cd /app/trading
docker compose -f docker-compose.yml up -d
```

---

## 🐛 Troubleshooting

### Check why a container exited
```bash
docker compose -f docker-compose.local.yml logs order-core
```

### Inspect container details
```bash
docker inspect <container_id>
```

### Test Redis connectivity
```bash
docker compose -f docker-compose.local.yml exec order-core \
  python -c "import redis; r = redis.Redis(host='redis'); print(r.ping())"
```

### Reset everything (⚠️ dangerous)
```bash
docker compose -f docker-compose.local.yml down -v
docker system prune -a
docker compose -f docker-compose.local.yml up -d
```

---

## 📊 Monitoring

### CPU and Memory usage
```bash
docker stats
```

### Disk usage
```bash
docker system df
```

### Network usage
```bash
docker stats --no-stream
```

### View all events
```bash
docker events --filter type=container
```

---

## 🔐 Security

### Update all base images
```bash
docker compose -f docker-compose.local.yml pull
docker compose -f docker-compose.local.yml up -d
```

### Scan image for vulnerabilities
```bash
docker scan alexanderloveiris/order-core:latest
```

### Run container as non-root
```bash
# Check Dockerfile - add USER directive
USER appuser
```

---

## 📝 Useful Tips

### Save terminal output to file
```bash
docker compose -f docker-compose.local.yml logs > debug.log 2>&1
```

### Follow multiple container logs
```bash
docker compose -f docker-compose.local.yml logs -f --tail=50
```

### Check environment variables in container
```bash
docker compose -f docker-compose.local.yml exec order-core env
```

### Copy file from container
```bash
docker compose -f docker-compose.local.yml cp order-core:/app/data.json ./
```

### Copy file to container
```bash
docker compose -f docker-compose.local.yml cp data.json order-core:/app/
```
