from confluent_kafka import Consumer, KafkaException, Producer
import json

# NOTE: This code is in a very early stage and has been placed here only
# to complete the data flow up to the 'tank' module. If int(price) is even, a notification to 'tank' should be sent.
# It should be refactored and moved to the appropriate code structure in the near future,
# following best practices for project organization and architecture.

def run():
  global producer
  producer_config = {
    "bootstrap.servers": "nabucodonosor:9092"
  }

  producer = Producer(producer_config)

  config = {
    "bootstrap.servers": "nabucodonosor:9092",  # your broker
    "group.id": "price-analyzer",              # consumer group name
    "auto.offset.reset": "earliest"            # start from the beginning if there's no previous offset
  }

  consumer = Consumer(config)
  consumer.subscribe(["price.db.new"])

  print("Listening for messages on 'price.db.new'...")
  try:
    while True:
      msg = consumer.poll(1.0)
      if msg is None:
        continue

      if msg.error():
        raise KafkaException(msg.error())

      raw = msg.value().decode("utf-8")
      obj = json.loads(raw)
      process_message(obj)

  except KeyboardInterrupt:
    print("\n Stopping consumer...")

  finally:
    consumer.close()


def process_message(price):
  symbol = price["symbol"]
  price_value = price["price"]
  currency = price["currency"]
  timestamp = price["timestamp"]
  print(f"Processing {symbol} {price_value} {currency} at {timestamp}")
  try:
    price_int = int(price_value)
    if price_int % 2 == 0: # Notify if int(price_value) is even (Testing purposes)
      print(f"Notifying tank: {symbol} price {price_int} is even.")
      payload = {
          "type": "new_price",
          "symbol": symbol,
          "price": price,
          "currency": currency,
          "timestamp": timestamp
      }

      publicate_message(payload)

  except (ValueError, TypeError):
    print(f"Warning: Could not convert price_value '{price_value}' to int for even check.")


def publicate_message(payload):
  try:
    producer.produce("notification.new", value=json.dumps(payload).encode("utf-8"))
    producer.flush()
    print(f"Notification sent to 'notification.new': {json.dumps(payload)}")

  except Exception as e:
    print(f"Error sending notification to Kafka: {e}")

