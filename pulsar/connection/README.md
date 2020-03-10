# Apache Pulsar Connection

This connection connects to apache pulsar.

### Flogo CLI
```bash
flogo install github.com/project-flogo/messaging-contrib/pulsar/connection
```

## Configuration

### Settings: 
| Name       | Type   | Description
|:---        | :---   | :---   
| url        | string | The url used to connect to pulsar - ***REQUIRED***
| athenzauth | string | The string used for Athenz Authentication
| certFile   | string | The location of the certificate file used in TLS.
| keyFile    | string | The location of the key file used in TLS.