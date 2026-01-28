
# Trading System Architecture – Redis Stream Blueprint

## 1. Overall Architecture

Signal/API/Bot → Order Core (Fast Path) → Redis Stream → Consumers

Order Core publishes events, all other services consume asynchronously.

---

## 2. Services Overview

### Order Core (Service A)
- Parse signal
- Pre-risk (O(1))
- Send order to exchange
- Publish ORDER_SENT

No DB, no TP/SL, no heavy logic.

### Order State (Service B)
- Persist orders
- Sync exchange state
- Emit ORDER_FILLED / ORDER_CLOSED

### Risk Engine
- Post-risk checks (exposure, leverage, loss limit)
- Kill switch
- Reduce or force close positions

### TP/SL Engine
- Create TP/SL after order filled
- Trailing stop
- Multi TP ladder

### Notify Service
- Telegram / Discord / Email
- Logging & alerts

---

## 3. Redis Stream Design

Streams:
- orders.fast
- orders.filled
- risk.events
- tp_sl.events

Consumer Groups:
- orders.db
- orders.risk
- orders.tp_sl
- orders.notify

---

## 4. Folder Structure

### order-core
```
order-core/
├── app/
│   ├── main.py
│   ├── services/order_service.py
│   ├── domain/order.py
│   ├── adapters/
│   │   ├── exchange/
│   │   └── stream/
│   ├── risk/pre_risk.py
│   └── dto/order_event.py
```

### order-state
```
order-state/
├── app/
│   ├── consumers/
│   ├── handlers/
│   ├── repositories/
│   ├── domain/
│   └── infra/
│       ├── db.py
│       └── redis_consumer.py
```

### risk-engine
```
risk-engine/
├── app/
│   ├── consumers/
│   ├── handlers/
│   │   ├── exposure_rule.py
│   │   ├── leverage_rule.py
│   │   └── kill_switch.py
│   ├── domain/
│   └── infra/
```

### tp-sl-engine
```
tp-sl-engine/
├── app/
│   ├── consumers/
│   ├── handlers/
│   │   ├── fixed_rr.py
│   │   ├── trailing_sl.py
│   │   └── ladder_tp.py
│   ├── domain/
│   └── adapters/exchange_client.py
```

### notify-service
```
notify-service/
├── app/
│   ├── consumers/
│   ├── handlers/
│   └── adapters/telegram.py
```

---

## 5. Patterns Used

- Event-Driven Architecture
- Fast Path / Slow Path
- Hexagonal (Clean Architecture)
- Strategy Pattern (TP/SL)
- Chain of Responsibility (Risk)
- Idempotent Consumer
- Consumer Groups (Redis Stream)

---

## 6. Event Schema Example

```
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

---

## 7. Production Checklist

- Order Core has no DB
- Redis Stream with consumer groups
- ACK after processing
- Idempotency by order_id
- Kill switch via Redis key
- Stream MAXLEN
- Deploy near exchange

---

## 8. Summary

Order Core is fast and dumb.
Redis Stream is the backbone.
Risk and TP/SL are async side effects.

This architecture scales, is low-latency, and safe for trading systems.
