# Apache Pulsar Connection

This connection connects to apache pulsar.

### Flogo CLI
```bash
flogo install github.com/project-flogo/messaging-contrib/pulsar/connection
```

## Configuration

### Settings: 
| Name          | Type   | Description
|:---           | :---   | :---   
| url           | string | The url used to connect to pulsar - ***REQUIRED***
| auth          | string | The type of authentication used: None, TLS, JWT, Athenz
| allowInsecure | bool   | Allow self signed certs or not
| athenzAuth    | params | The params used for Athenz Authentication
| jwt           | string | The JWT authentication token
| caCert        | string | The location of the ca cert file used in TLS.
| certFile      | string | The location of the certificate file used in TLS.
| keyFile       | string | The location of the key file used in TLS.