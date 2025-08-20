# Go Kafka Producer & Consumer (Sarama + Fiber)

A tiny end‑to‑end example that exposes an HTTP API to publish **comments** to Kafka (producer) and a standalone **consumer** that reads from the same topic and logs messages.

* **Producer:** Go + [Fiber](https://github.com/gofiber/fiber) → `POST /api/comments` → Kafka topic `comments`
* **Consumer:** Go + [Sarama](https://github.com/IBM/sarama) → reads `comments` (partition 0) from `localhost:9092`

---

## 🧭 Project Layout (suggested)

```
./
├─ consumer/
│  └─ main.go          # your consumer code (prints and counts messages)
├─ producer/
│  └─ main.go          # your Fiber HTTP API (POST /api/comments)
├─ go.mod
└─ README.md
```

> If you keep files in a single folder, use separate file names like `consumer.go` and `producer.go` and run them individually.

---

## ✅ Prerequisites

* **Go** 1.20+
* **Kafka** available at `localhost:9092`
* **Topic**: `comments`

### Option A: Run Kafka quickly with Docker (KRaft / no ZooKeeper)

Create `docker-compose.yml`:

```yaml
version: '3.8'
services:
  kafka:
    image: bitnami/kafka:latest
    container_name: kafka
    ports:
      - '9092:9092'
    environment:
      - KAFKA_CFG_NODE_ID=1
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@localhost:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
```

Start it:

```bash
docker compose up -d
```

Create the topic:

```bash
docker exec -it kafka kafka-topics.sh \
  --create --topic comments --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```

### Option B: Local install

If you installed Kafka natively, start the broker and run:

```bash
kafka-topics.sh --create --topic comments --bootstrap-server localhost:9092 \
  --partitions 1 --replication-factor 1
```

---

## 📦 Dependencies

From your project root:

```bash
go mod init github.com/<you>/kafka-comments

go get github.com/IBM/sarama@latest
go get github.com/gofiber/fiber/v2@latest
```

---

## 🚀 Running the Producer (HTTP → Kafka)

The producer exposes `POST /api/comments` which accepts JSON like `{ "text": "hello" }` and publishes to Kafka topic **comments**.

```bash
# from the folder where producer main.go lives
go run .   # or: go run producer/main.go
```

**Test with curl:**

```bash
curl -X POST http://localhost:3000/api/comments \
  -H 'Content-Type: application/json' \
  -d '{"text":"first comment"}'
```

Expected response:

```json
{
  "success": true,
  "message": "Comment pushed successfully",
  "comments": "{\"text\":\"first comment\"}"
}
```

> Producer settings are in `ConnectProducer`: broker defaults to `localhost:9092`, sync producer, acks=all, retries=5.

---

## 📥 Running the Consumer (Kafka → Logs)

The consumer connects to `localhost:9092`, reads from **topic `comments`, partition 0**, starting at **oldest** offset, logs and counts messages. It exits cleanly on **Ctrl+C** (SIGINT/SIGTERM).

```bash
# from the folder where consumer main.go lives
go run .   # or: go run consumer/main.go
```

Expected logs:

```
Consumer Started
Received Message: {"text":"first comment"} | count 1
```

> Consumer settings are in `connectConsumer`: `Consumer.Return.Errors=true` and `ConsumePartition(topic, 0, sarama.OffsetOldest)`.

---

## 🔧 Configuration

* **Broker URL:** change `localhost:9092` in `ConnectProducer` / `connectConsumer`.
* **Topic name:** change the `topic := "comments"` and the `PushCommentToQueue("comments", ...)` call.
* **Partition:** consumer is hard-coded to partition `0`. For production, prefer **consumer groups** (see Next Steps).

---

## 🧪 End‑to‑End Test

1. Start Kafka and create topic `comments`.
2. Run **consumer**.
3. Run **producer**.
4. `curl` the producer endpoint with a comment JSON.
5. See the consumer print the message and increment `count`.
6. Hit **Ctrl+C** in the consumer terminal to shut it down gracefully.

---

## 📜 License

MIT (or your choice)

---

## 🙌 Credits

* [Sarama](https://github.com/IBM/sarama)
* [Fiber](https://github.com/gofiber/fiber)

Happy streaming! 💨
