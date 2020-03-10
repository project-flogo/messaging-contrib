<!--
title: Kafka
weight: 4701
-->
# Kafka Trigger

This trigger subscribes to a topic on Kafka cluster and listens for the messages.

### Flogo CLI
```bash
flogo install github.com/project-flogo/messaging-contrib/kafka/trigger
```

## Configuration

### Setting :

| Name       | Type   | Description
|:---        | :---   | :---     
| connection | any    | The connection object which is use to connect to Kafka - ***REQUIRED*** [Connection](../connection/README.md)

### HandlerSettings:

| Name       | Type   | Description
|:---        | :---   | :---   
| topic      | string | The Kafka topic on which to listen for messages
| partitions | string | The specific partitions to consume messages from
| offset     | int64  | The offset to use when starting to consume messages

### Output:

| Name         | Type     | Description
|:---          | :---     | :---   
| message      | string   | The message that was consumed


## Examples

```json
{
  "triggers": [
    {
      "id": "flogo-kafka",
      "ref": "github.com/project-flogo/contrib/trigger/kafka",
      "settings": {
        "brokerUrls" : "localhost:9092",
        "trustStore" : "" 
      },
      "handlers": [
        {
          "settings": {
            "topic": "syslog",
          },
          "action": {
            "ref": "github.com/project-flogo/flow",
            "settings": {
              "flowURI": "res://flow:my_flow"
            }
          }
        }
      ]
    }
  ]
}
```
 
## Development

### Testing

To run tests first set up the kafka broker using the docker-compose file given below:

```yaml
version: '2'
  
services:

  zookeeper:
    image: wurstmeister/zookeeper:3.4.6
    expose:
    - "2181"

  kafka:
    image: wurstmeister/kafka:2.11-2.0.0
    depends_on:
    - zookeeper
    ports:
    - "9092:9092"
    environment:
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
```

Then run the following command: 

```bash
go test 
```
