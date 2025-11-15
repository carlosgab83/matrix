from confluent_kafka import Consumer, KafkaException
import json

def _run():
  print("RUNNING111")


def run():
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

      # decode JSON
      raw = msg.value().decode("utf-8")
      price = json.loads(raw)

      # print("Received message:", price)

      # --- here you perform your calculations and analysis ----
      symbol = price["symbol"]
      price_value = price["price"]
      currency = price["currency"]
      timestamp = price["timestamp"]

      print(f"Processing {symbol} {price_value} {currency} at {timestamp}")
      # ------------------------------------------------

  except KeyboardInterrupt:
    print("\n Stopping consumer...")
  finally:
    consumer.close()
