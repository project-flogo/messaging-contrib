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


For Example:

```json
"connections": {
    "a7730ae0-0199-11ea-9e1b-1b6d6afda999": {
      "ref": "github.com/project-flogo/messaging-contrib/pulsar/connection",
      "settings": {
        "description": "Plain text connection",
        "name": "plain",
        "url": "pulsar+ssl://flex-linux-gazelle.na.tibco.com:16651",
        "auth":"None"
      }
    }
  }

"connections": {
    "a7730ae0-0199-11ea-9e1b-1b6d6afda999": {
      "ref": "github.com/project-flogo/messaging-contrib/pulsar/connection",
      "settings": {
        "description": "Tls encrypted connection with JWT auth",
        "name": "jwt",
        "url": "pulsar+ssl://flex-linux-gazelle.na.tibco.com:26651",
        "auth":"JWT",
        "caCert":"/gil/local/pub/certs/pulsar_tls_jwt/certs/my-ca/certs/ca.cert.pem",
        "jwt":"eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ0aWJjby11c2VyIn0.KUx6eIwNKaa5wb_SqjFQViXX5wAF5gAJGvTCbi6q1Hw",
        "allowInsecure":true
      }
    }
  }

  "connections": {
    "a7730ae0-0199-11ea-9e1b-1b6d6afda999": {
      "ref": "github.com/project-flogo/messaging-contrib/pulsar/connection",
      "settings": {
        "description": "Encrypted connection using TLS for auth",
        "name": "tlsauth",
        "url": "pulsar+ssl://flex-linux-gazelle.na.tibco.com:36651",
        "auth":"TLS",
        "caCert":"/gil/local/pub/certs/pulsar_tls_jwt/certs/my-ca/certs/ca.cert.pem",
        "certFile":"/gil/local/pub/certs/pulsar_tls_jwt/certs/my-ca/client_cert.pem",
        "keyFile":"/gil/local/pub/certs/pulsar_tls_jwt/certs/my-ca/client_key.pem",
        "allowInsecure":true
      }
    }
  }


```