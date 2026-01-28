# AI IMPLEMENTATION SPEC – EVENT‑DRIVEN TRADING SYSTEM
*(Redis Stream · Docker · Digest Deploy · Go + Python Hybrid)*

Tài liệu này dùng **LÀM INPUT TRỰC TIẾP CHO AI IMPLEMENT** toàn bộ hệ thống.
Không phải giải thích – mà là **spec + contract + structure**.

---

## 0. GOAL

Xây dựng hệ thống trading:
- Low latency
- Event‑driven
- Redis Stream làm backbone
- Dockerized
- CI/CD dùng digest (deterministic deploy)

---

## 1. ARCHITECTURE OVERVIEW

```
Signal / API
   ↓
Order Core (Go – FAST PATH)
   ↓  XADD
Redis Stream
   ↓
┌────────────┬──────────────┬──────────────┐
│ OrderState │ Risk Engine  │ TP/SL Engine │
│ (Go)       │ (Go+Py)      │ (Python)     │
└────────────┴──────────────┴──────────────┘
                 ↓
           Notify Service
```

RULES:
- Order Core: NO DB, NO TP/SL, NO heavy risk
- All services communicate ONLY via Redis Stream
- No HTTP between services

---

## 2. SERVICES & RESPONSIBILITIES

### 2.1 order-core (Go)
- Parse signal
- Pre-risk O(1)
- Send order to exchange
- Publish ORDER_SENT

### 2.2 order-state (Go)
- Persist orders
- Sync exchange fills
- Emit ORDER_FILLED / ORDER_CLOSED

### 2.3 risk-engine (Go + Python rules)
- Post-risk checks
- Kill-switch
- Reduce / force close

### 2.4 tp-sl-engine (Python)
- Build TP/SL strategy
- Trailing stop
- Ladder TP

### 2.5 notify-service (Python)
- Telegram / alert
- Logging

---

## 3. REDIS STREAM CONTRACT

### Streams
```
orders.fast
orders.filled
risk.events
tp_sl.events
```

### Consumer Groups
```
orders.db
orders.risk
orders.tp_sl
orders.notify
```

---

## 4. EVENT SCHEMA (IMMUTABLE CONTRACT)

### ORDER_SENT
```json
{
  "event": "ORDER_SENT",
  "version": "v1",
  "order_id": "uuid",
  "symbol": "BTCUSDT",
  "side": "BUY",
  "qty": "0.01",
  "ts": 1700000000
}
```

### ORDER_FILLED
```json
{
  "event": "ORDER_FILLED",
  "order_id": "uuid",
  "price": "5080",
  "qty": "0.01",
  "ts": 1700000001
}
```

### RISK_VIOLATION
```json
{
  "event": "RISK_VIOLATION",
  "order_id": "uuid",
  "rule": "MAX_EXPOSURE",
  "action": "REDUCE",
  "ts": 1700000002
}
```

---

## 5. FOLDER STRUCTURE (MONOREPO)

```
trading-system/
├── order-core/        (Go)
├── order-state/       (Go)
├── risk-engine/       (Go + Py rules)
├── tp-sl-engine/      (Python)
├── notify-service/    (Python)
├── docker-compose.local.yml
├── docker-compose.yml
└── .github/workflows/
```

---

## 6. SERVICE INTERNAL STRUCTURE (MANDATORY)

### Go service
```
/app
  /cmd
  /domain
  /service
  /adapter
  /infra
```

### Python service
```
/app
  /consumers
  /handlers
  /domain
  /infra
```

---

## 7. DOCKER RULES

- One service = one container
- Stateless
- No Redis / DB in prod compose
- Prod deploy uses DIGEST ONLY

---

## 8. docker-compose.yml (DIGEST)

```yaml
order-core:
  image: alexanderloveiris/order-core@${ORDER_CORE_DIGEST}
```

All services MUST follow this.

---

## 9. CI/CD CONTRACT

GitHub Actions MUST:
1. Build images
2. Push Docker Hub
3. Extract image digest
4. SSH to VPS
5. Update `.env.prod`
6. `docker compose pull && up -d`

---

## 10. VPS CONTRACT

VPS MUST:
- Have Docker + docker-compose
- Have `/opt/trading`
- Contain only:
  - docker-compose.yml
  - .env.prod

NO SOURCE CODE ON VPS.

---

## 11. KILL SWITCH

- Redis key: `kill_switch`
- If = 1 → Order Core rejects all new orders
- O(1) check

---

## 12. NON‑NEGOTIABLE RULES

- No synchronous cross-service calls
- No business logic in Order Core
- No Pub/Sub (Redis Stream only)
- Every consumer MUST be idempotent
- ACK after success only

---

## 13. SUCCESS CRITERIA

- Order latency < 50ms internal
- Deterministic deploy
- Rollback in < 10 seconds
- Service restart safe
- Horizontal scale via consumer groups

---

## 14. AI IMPLEMENTATION INSTRUCTION

AI SHOULD:
- Generate production-ready code
- Respect folder structure
- Respect event schema
- Respect Redis Stream semantics
- Not introduce extra services
- Not introduce HTTP calls

This document is the SINGLE SOURCE OF TRUTH.
