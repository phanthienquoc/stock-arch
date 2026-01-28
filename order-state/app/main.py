#!/usr/bin/env python3
"""
Order State Service
- Persist orders
- Sync exchange state
- Emit ORDER_FILLED / ORDER_CLOSED events
"""

import os
import json
import logging
from dotenv import load_dotenv
import redis

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

REDIS_HOST = os.getenv('REDIS_HOST', 'redis')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
CONSUMER_GROUP = 'orders.db'
STREAM_NAME = 'orders.fast'

r = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, decode_responses=True)

def setup_consumer_group():
    """Create consumer group if not exists"""
    try:
        r.xgroup_create(STREAM_NAME, CONSUMER_GROUP, id='0', mkstream=True)
        logger.info(f"Consumer group '{CONSUMER_GROUP}' created")
    except redis.ResponseError as e:
        if 'BUSYGROUP' in str(e):
            logger.info(f"Consumer group '{CONSUMER_GROUP}' already exists")
        else:
            logger.error(f"Error creating consumer group: {e}")

def consume_orders():
    """Consume orders from Redis Stream"""
    try:
        setup_consumer_group()
        logger.info("Order State Service started, listening for orders...")
        
        while True:
            messages = r.xreadgroup(CONSUMER_GROUP, 'consumer-1', {STREAM_NAME: '>'}, block=1000)
            if messages:
                for stream, group_messages in messages:
                    for msg_id, data in group_messages:
                        logger.info(f"Processing order: {data}")
                        # TODO: Persist to DB and sync with exchange
                        r.xack(STREAM_NAME, CONSUMER_GROUP, msg_id)
    except Exception as e:
        logger.error(f"Error in consumer loop: {e}")

def main():
    consume_orders()

if __name__ == '__main__':
    main()
