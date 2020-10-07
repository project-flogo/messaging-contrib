
# Apache Pulsar Subscriber
This trigger allows your flogo application to listen on an Apache Pulsar topic.

## Installation

```bash
flogo install github.com/project-flogo/messaging-contrib/pulsar/trigger/subscriber
```

## Configuration

### Settings:
| Name       | Type   | Description
|:---        | :---   | :---       
| connection | any    | The connection object which is used to connect to pulsar - ***REQUIRED*** [Connection](../connection/README.md)

### Handler Settings:
| Name             | Type    | Description
|:---              | :---    | :---          
| topic            | string  | The Pulsar topic from which to get the message - ***REQUIRED***
| subscription     | string  | The subscription name - **REQUIRED**
| subscriptionType | string  | The subscription type: Exclusive, Shared, Failover or KeyShared, defaults to Shared
| initialPosition  | string  | The initial position upon startup: Latest or Earliest, defaults to Latest
| dlqTopic         | string  | If provided, implements dead letter topic processing
| dlqMaxDeliveries | integer | The number of times message processing will be attempted before being relocated to dlqtopic

### Output:
| Name        | Type   | Description
|:---         | :---   | :---        
| payload     | any    | The contents of the message from Pulsar.
| properties  | params | The properties associated with the message
| topic       | string | The topic to which the message was published


### Example:
```json
{
  "name": "testpulsar",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "",
  "appModel": "1.1.0",
  "imports": [
    "github.com/project-flogo/contrib/activity/log",
    "github.com/project-flogo/flow",
    "github.com/project-flogo/messaging-contrib/pulsar/trigger/subscriber",
    "github.com/project-flogo/messaging-contrib/pulsar/connection"
  ],
  "triggers": [
    {
      "id": "receive_pulsar_messages",
      "ref": "#subscriber",
      "settings": {
        "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda999"
      },
      "handlers": [
        {
          "settings": {
            "initialposition": "Latest",
            "subscriptionName": "wcnsub",
            "subscriptionType": "Shared",
            "topic": "wcnString"
          },
          "actions": [
            {
              "ref": "#flow",
              "settings": {
                "flowURI": "res://flow:receive"
              },
              "input": {
                "payload": "=$.payload"
              }
            }
          ]
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "flow:receive",
      "data": {
        "name": "Receive",
        "metadata": {
          "input": [
            {
              "name": "payload",
              "type": "string"
            }
          ],
          "output": [
            {
              "name": "rc",
              "type": "integer"
            }
          ]
        },
        "tasks": [
          {
            "id": "log_2",
            "name": "Log",
            "description": "Logs a message",
            "activity": {
              "ref": "#log",
              "input": {
                "message": "@flow:receive:payload",
                "addDetails": false,
                "usePrint": false
              }
            }
          }
        ],
        "links": []
      }
    }
  ],
  "connections": {
    "a7730ae0-0199-11ea-9e1b-1b6d6afda999": {
      "ref": "github.com/project-flogo/messaging-contrib/pulsar/connection",
      "settings": {
        "description": "",
        "name": "mc2",
        "url": "pulsar://flex-linux-gazelle.na.tibco.com:16650"
      }
    }
  }
}
```