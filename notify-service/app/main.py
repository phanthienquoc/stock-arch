#!/usr/bin/env python3
"""
Notify Service
- Send notifications to Telegram / Discord / Email
- Logging & alerts
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
CONSUMER_GROUP = 'orders.notify'
STREAM_NAME = 'risk.events'

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

def send_notification(event_data: dict):
    """Send notification via Telegram/Discord/Email"""
    # TODO: Implement notification sending logic
    logger.info(f"Sending notification: {event_data}")
    return True

def consume_events():
    """Consume events from Redis Stream"""
    try:
        setup_consumer_group()
        logger.info("Notify Service started, listening for events...")
        
        while True:
            messages = r.xreadgroup(CONSUMER_GROUP, 'consumer-1', {STREAM_NAME: '>'}, block=1000)
            if messages:
                for stream, group_messages in messages:
                    for msg_id, data in group_messages:
                        logger.info(f"Processing notification event: {data}")
                        send_notification(data)
                        r.xack(STREAM_NAME, CONSUMER_GROUP, msg_id)
    except Exception as e:
        logger.error(f"Error in consumer loop: {e}")

def main():
    consume_events()

if __name__ == '__main__':
    main()
