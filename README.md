# Time Converter

This project start consuming time in epoch format from kafka topic and convert it to RFC3339 format, then produce it to another kafka topic.

**Default Topics are:** 
- Consumer start consuming from **input** topic
- Producer start produce to **output** topic

**This project contains three part:**
- Consumer: Consume epoch time from kafka topic
- Producer: Convert and produce time in RFC3339 format to kafka topic
- Metrics: collect metrics from producer and consumer.

**Note:** there is consumer group and you can scale it with replica on kubernetes

### Config Global

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `LOG_LEVEL` | log level | enum(trace, debug, info, warn, error, fatal, panic) | info

### Config Exporter

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `METRICS_PATH` | metrics path | string | /metrics
| `LISTEN_PORT` | exporter listening port  | number | 8081


### Config Consumer

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `CONSUMER_VERSION` | Kafka Version | string | 2.7.0
| `CONSUMER_GROUP` | Consumer Group Name | string | epochConsumer
| `CONSUMER_TOPIC` | Kafka Topic to Consume | string | input
| `CONSUMER_BROKERS` | Kafka Brokers | comma seperated hosts:port | 127.0.0.1:9092
| `CONSUMER_FETCH_MIN` | Consumer.Fetch.Min | string | 1
| `CONSUMER_FETCH_DEFAULT` | Consumer.Fetch.Default | number | 1048576
| `CONSUMER_RETRY_BACKOFF` | Consumer.Retry.Backoff | number(secound) | 2
| `CONSUMER_MAXWAITTIME` | Consumer.MaxWaitTime | number(milisecound) | 250
| `CONSUMER_MAXPROCESSINGTIME` | Consumer.MaxProcessingTime | number(milisecound) | 100
| `CONSUMER_RETURN_ERROR` | Consumer.Return.Errors | bool | true
| `CONSUMER_OFFSETS_AUTO_COMMIT_ENABLE` | Consumer.Offsets.AutoCommit.Enable | bool | true
| `CONSUMER_OFFSETS_AUTO_COMMIT_INTERVAL` | Consumer.Offsets.AutoCommit.Interval | number(secound) | 1
| `CONSUMER_OFFSETS_INITIAL` | Consumer.Offsets.Initial | enum(oldest , newest) | newest
| `CONSUMER_OFFSETS_RETRY_MAX` | Consumer.Offsets.Retry.Max | number | 3
| `CONSUMER_GROUP_SESSION_TIMEOUT` | Consumer.Group.Session.Timeout | number(secound) | 10
| `CONSUMER_GROUP_HEARTBEAT_INTERVAL` | Consumer.Group.Heartbeat.Interval | number(secound) | 3
| `CONSUMER_GROUP_REBALANCE_STRATEGY` | Consumer.Group.Rebalance.Strategy | enum(range , sticky, rr) | 3
| `CONSUMER_GROUP_REBALANCE_TIMEOUT` | Consumer.Group.Rebalance.Timeout | number(secound) | 60
| `CONSUMER_GROUP_REBALANCE_RETRY_MAX` | Consumer.Group.Rebalance.Retry.Max | number | 4
| `CONSUMER_GROUP_REBALANCE_RETRY_BACKOFF` | Consumer.Group.Rebalance.Retry.Backoff | number(secound) | 2
| `CONSUMER_CLIENTID` | Consumer ClientID  | string | defaultClientID
| `CONSUMER_CHANNELBUFFERSIZE` | ChannelBufferSize | number | 256

Almost all descriptions are map to sarama
[config.go](https://github.com/Shopify/sarama/blob/master/config.go#L441)
values


### Config Producer

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `PRODUCER_VERSION` | Kafka Version | string | 2.7.0
| `PRODUCER_TOPIC` | Kafka Topic to Produce | string | output
| `PRODUCER_BROKERS` | Kafka Brokers | comma seperated hosts:port | 127.0.0.1:9092
| `PRODUCER_MAX_MESSAGE_BYTES` | Producer.MaxMessageBytes | number | 1000000
| `PRODUCER_FLUSH_FREQUENCY` | Producer.Flush.Frequency | number(milisecound) | -
| `PRODUCER_FLUSH_MESSAGE` | Producer.Flush.Messages | number | -
| `PRODUCER_FLUSH_MAX_MESSAGE` | Producer.Flush.MaxMessages | number | -
| `PRODUCER_RETURN_SUCCESS` | Producer.Return.Successes | bool | true
| `PRODUCER_TIMEOUT` | Producer.Timeout | number(secound) | 10
| `PRODUCER_RETRY_MAX` | Producer.Retry.Max | number | 3
| `PRODUCER_RETRY_BACKOFF` | Producer.Retry.Backoff | number(milisecound) | 100
| `PRODUCER_RETURN_ERROR` | Producer.Return.Errors | bool | true
| `PRODUCER_COMPRESSIONLEVEL` | Producer.CompressionLevel | enum(gzip, zstd, snappy, lz4, none) | none
| `PRODUCER_PARTITIONER` | Producer.Partitioner | enum(random, hash, rr) | hash
| `PRODUCER_REQUIRED_ACKS` | Producer.RequiredAcks | enum(0, 1, -1) | 1
| `PRODUCER_CLIENTID` | Producer ClientID  | string | defaultClientID
| `PRODUCER_CHANNELBUFFERSIZE` | ChannelBufferSize | number | 256

**Note:** Required acks enums map to:
- 0: NoResponse
- 1: WaitForLocal
- -1: WaitForAll

Almost all descriptions are map to sarama
[config.go](https://github.com/Shopify/sarama/blob/master/config.go#L441)
values


# Metrics

The following metrics are available:

|Name|Description|
|---|---|
|`produced_message_total`|How many RFC3339 message produced|
|`convert_message_total`|How many epoch converted to RFC3339|

Metrics are counters and might be used with
[`rate()`](https://prometheus.io/docs/prometheus/latest/querying/functions/#rate())
to calculate per-second average rate.

Tests
-----
Run tests are require a reachable kafka broker
```
go test ./...
```

Build
-----
```
docker build .
[...]
Successfully built a7ec197b9502
```

Run container
-------------
```
docker run -p 8081:8081 -e KAFKA_BROKER=localhost:9092 a7ec197b9502
```

Test
----
Run Kafka console producer on topic input
```
kafka-console-producer --bootstrap-server localhost:9092 --topic input
1615035196527
1615035196529
1615035196531
```

Run Kafka console consumer on topic ouput
```
kafka-console-consumer --bootstrap-server localhost:9092 --topic output
```
verify RFC3339 Messages
```
2021-03-06T12:53:16.527Z
2021-03-06T12:53:16.529Z
2021-03-06T12:53:16.531Z
```