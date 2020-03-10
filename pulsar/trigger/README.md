
# Apache Pulsar Trigger
This trigger provides your flogo application the ability to listen from the Apache Pulsar.

## Installation

```bash
flogo install github.com/project-flogo/contrib/trigger/pulsar
```

## Configuration

### Settings:
| Name      | Type   | Description
|:---       | :---   | :---       
| connection| any    | The connection object which is use to connect to pulsar - ***REQUIRED*** [Connection](../connection/README.md)

### Handler Settings:
| Name         | Type   | Description
|:---          | :---   | :---          
| topic        | string | The Pulsar topic from which to get the message - ***REQUIRED***
| subscription | string | The subscription name - **REQUIRED**

### Output:
| Name        | Type   | Description
|:---         | :---   | :---        
| message     | string | The message from the Pulsar.

