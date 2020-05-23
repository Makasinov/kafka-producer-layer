# Kafka-Consumer

Kafka consumer based on forked repo example from https://github.com/confluentinc/confluent-kafka-go

# API
| Path    | Method | Description                                  |
|---------|--------|----------------------------------------------|
| /config | GET    | Get configuration for kafka producer         |
| /:topic | POST   | Send message in body to kafka specific topic |

# Envs
| Name    | Description                                  |
|---------|----------------------------------------------|
| Port    |   Listen and serve                           |
