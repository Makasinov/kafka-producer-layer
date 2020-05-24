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

# Usage
```sh
curl -d '{ "column1_in_clickhouse": "value", "column2_in_clickhouse": "value" }' localhost:8080/my_kafka_topic_name
```
