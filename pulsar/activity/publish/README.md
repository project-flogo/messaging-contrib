
# Apache Pulsar Activity

This activity publishes messages on a topic in a Apache Pulsar.

### Flogo CLI
```bash
flogo install github.com/project-flogo/messaging-contrib/pulsar/activity/publish
```

## Configuration

### Settings: 
| Name              | Type   | Description
|:---               | :---   | :---   
| connection        | any    | The connection object which is use to connect to pulsar - ***REQUIRED*** [Connection](../connection/README.md)
| topic             | string | The Pulsar topic on which to place the message - ***REQUIRED***
| compressionType   | string | The type of compression to use: "NONE","LZ4","ZLIB","ZSTD" defaults to "NONE"

### Input:

| Name       | Type   | Description
|:---        | :---   | :---  
| payload    | any    | The message to send 


### Output:

| Name       | Type   | Description
|:---        | :---   | :---  
| msgid      | string | The message identifier


### Example:
```json
{
  "name": "PulsarPublish",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "",
  "appModel": "1.1.0",
  "imports": [
    "github.com/project-flogo/contrib/trigger/timer",
    "github.com/project-flogo/flow",
    "github.com/project-flogo/contrib/activity/log",
    "github.com/project-flogo/messaging-contrib/pulsar/activity/publish",
    "github.com/project-flogo/messaging-contrib/pulsar/connection",
    "github.com/project-flogo/messaging-contrib/pulsar@v0.0.0"
  ],
  "triggers": [
    {
      "id": "timer",
      "ref": "#timer",
      "settings": null,
      "handlers": [
        {
          "settings": null,
          "actions": [
            {
              "ref": "#flow",
              "settings": {
                "flowURI": "res://flow:publish"
              }
            }
          ]
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "flow:publish",
      "data": {
        "name": "Publish",
        "tasks": [
          {
            "id": "activity_2",
            "name": "Publish Pulsar message",
            "description": "Publish a message to a pulsar topic",
            "activity": {
              "ref": "#publish",
              "input": {
                "payload": "mary had a little lamb"
              },
              "settings": {
                "connection": "conn://a7730ae0-0199-11ea-9e1b-1b6d6afda999",
                "topic": "wcnString"
              }
            }
          }
        ]
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