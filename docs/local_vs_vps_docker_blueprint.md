# Deployment Blueprint – Local vs VPS (Docker + Redis Stream)

Tài liệu này chia **RÕ RÀNG**:
- Việc cần làm ở **LOCAL (dev / CI)**
- Việc cần làm ở **VPS / EC2 (prod runtime)**

Mục tiêu:
> VPS chỉ chạy container + pull image  
> CI/CD chịu trách nhiệm build & update version (digest)

---

## 1️⃣ Folder Structure – Chuẩn cho MONO-REPO (LOCAL)

```
trading-system/
├── .github/
│   └── workflows/
│       ├── docker-release.yml
│       └── deploy.yml
│
├── order-core/
│   ├── Dockerfile
│   ├── requirements.txt
│   └── app/
│
├── order-state/
│   ├── Dockerfile
│   └── app/
│
├── risk-engine/
│   ├── Dockerfile
│   └── app/
│
├── tp-sl-engine/
│   ├── Dockerfile
│   └── app/
│
├── notify-service/
│   ├── Dockerfile
│   └── app/
│
├── docker-compose.local.yml
├── docker-compose.yml
└── README.md
```

---

## 2️⃣ docker-compose.local.yml (LOCAL DEV)

👉 Dùng cho:
- Dev
- Test integration
- Debug flow Redis Stream

```yaml
version: "3.9"

services:
  redis:
    image: redis:7
    ports:
      - "6379:6379"

  order-core:
    build: ./order-core
    env_file: .env.local
    depends_on:
      - redis

  order-state:
    build: ./order-state
    env_file: .env.local
    depends_on:
      - redis

  risk-engine:
    build: ./risk-engine
    env_file: .env.local
    depends_on:
      - redis

  tp-sl-engine:
    build: ./tp-sl-engine
    env_file: .env.local
    depends_on:
      - redis

  notify-service:
    build: ./notify-service
    env_file: .env.local
    depends_on:
      - redis
```

---

## 3️⃣ docker-compose.yml (VPS / EC2)

👉 **KHÔNG build**
👉 **KHÔNG chạy Redis / DB**
👉 **CHỈ pull image bằng digest**

```yaml
version: "3.9"

services:
  order-core:
    image: alexanderloveiris/order-core@${ORDER_CORE_DIGEST}
    restart: always

  order-state:
    image: alexanderloveiris/order-state@${ORDER_STATE_DIGEST}
    restart: always

  risk-engine:
    image: alexanderloveiris/risk-engine@${RISK_ENGINE_DIGEST}
    restart: always

  tp-sl-engine:
    image: alexanderloveiris/tp-sl-engine@${TPSL_ENGINE_DIGEST}
    restart: always

  notify-service:
    image: alexanderloveiris/notify-service@${NOTIFY_DIGEST}
    restart: always
```

---

## 4️⃣ Folder Structure trên VPS (CHỈ 1 LẦN)

```
/opt/trading/
├── docker-compose.yml
├── .env
└── logs/            (optional)
```

📌 VPS **KHÔNG chứa source code**

---

## 5️⃣ Việc CẦN LÀM Ở LOCAL (Developer / CI)

### ✅ LOCAL – Checklist
- Viết code
- Viết Dockerfile cho từng service
- Test bằng `docker-compose.local.yml`
- Push code lên GitHub
- Tag version (`git tag v1.2.0`)
- GitHub Actions:
  - Build image
  - Push Docker Hub
  - Lấy digest
  - SSH deploy

📌 LOCAL **KHÔNG deploy trực tiếp**

---

## 6️⃣ Việc CẦN LÀM Ở VPS (CHỈ 1 LẦN)

### ✅ VPS – Initial Setup

```bash
# 1. Install Docker
curl -fsSL https://get.docker.com | sh

# 2. Enable Docker
sudo systemctl enable docker
sudo systemctl start docker

# 3. Allow non-root
sudo usermod -aG docker $USER
logout && login
```

---

### 7️⃣ Chuẩn bị thư mục deploy trên VPS

```bash
sudo mkdir -p /opt/trading
sudo chown -R $USER:$USER /opt/trading
cd /opt/trading
```

Copy vào:
- `docker-compose.yml`
- `.env`

---

## 8️⃣ .env (ban đầu – sẽ bị CI overwrite)

```env
ORDER_CORE_DIGEST=sha256:dummy
ORDER_STATE_DIGEST=sha256:dummy
RISK_ENGINE_DIGEST=sha256:dummy
TPSL_ENGINE_DIGEST=sha256:dummy
NOTIFY_DIGEST=sha256:dummy

REDIS_HOST=managed-redis.xxxx
DB_HOST=managed-db.xxxx
```

---

## 9️⃣ systemd auto-start (VPS)

```ini
[Unit]
Description=Trading System
After=docker.service
Requires=docker.service

[Service]
WorkingDirectory=/opt/trading
ExecStart=/usr/bin/docker compose -f docker-compose.yml up -d
ExecStop=/usr/bin/docker compose -f docker-compose.yml down
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable:
```bash
sudo systemctl daemon-reload
sudo systemctl enable trading
```

---

## 🔟 TÓM TẮT 1 DÒNG (RẤT QUAN TRỌNG)

> LOCAL = code + build + version  
> VPS = pull image + run container  
> CI/CD = cầu nối DUY NHẤT giữa 2 bên

---

## 11️⃣ Khi nào VPS CẦN can thiệp lại?

| Trường hợp | Cần SSH VPS |
|---|---|
| Thay compose | ✅ |
| Thay infra | ✅ |
| Deploy version mới | ❌ |
| Rollback | ❌ |

---

**Tài liệu này là “xương sống vận hành” cho hệ thống trading dùng Docker + Redis Stream.**
