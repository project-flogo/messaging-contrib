# Kafka Activity

This connection connects to kafka cluster.

### Flogo CLI
```bash
flogo install github.com/project-flogo/messaging-contrib/kafka/connection
```

## Configuration

### Settings: 
| Name       | Type   | Description
|:---        | :---   | :---   
| brokerUrls | string | The brokers of the Kafka cluster to connect to - ***REQUIRED***
| user       | string | If connecting to a SASL enabled port, the user id to use for authentication
| password   | string | If connecting to a SASL enabled port, the password to use for authentication 
| trustStore | string | If connecting to a TLS secured port, the directory containing the certificates representing the trust chain for the connection. This is usually just the CACert used to sign the server's certificate