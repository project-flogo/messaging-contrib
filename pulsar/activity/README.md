
# Apache Pulsar Activity

This activity publishes messages on a topic in a Apache Pulsar.

### Flogo CLI
```bash
flogo install github.com/project-flogo/contrib/activity/pulsar
```

## Configuration

### Settings: 
| Name       | Type   | Description
|:---        | :---   | :---   
| connection | any    | The connection object which is use to connect to pulsar - ***REQUIRED***

### Input:

| Name       | Type   | Description
|:---        | :---   | :---  
| payload    | string | The message to send 
