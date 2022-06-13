# Client Application

Provides GRPC and HTTP APIs to interact with the OCF devices via D2D communication.

## Run

The client application can be run by executing the following command:

```bash
./client-application
```

* Default `config.yaml` file will be generated in case it's not in the same directory as the client application, or if a different path with the existing config isn't specified.
* Default web UI files will be copied to the `www` folder next to the client application in case it doesn't exist, or if a different path with existing UI files isn't specified in the config file.

### Supported options

- `--config`: path to the config file
- `--version`: print the version of the client application

## YAML Configuration

A configuration template is available on [config.yaml](./config.yaml).

### Logging

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `log.dumpBody` | bool | `Set to true if you would like to dump raw messages.` | `false` |
| `log.level` | string | `Logging enabled from level.` | `"info"` |
| `log.encoding` | string | `Logging format. The supported values are: "json", "console"` | `"json"` |
| `log.stacktrace.enabled` | bool | `Log stacktrace.` | `"false` |
| `log.stacktrace.level` | string | `Stacktrace from level.` | `"warn` |
| `log.encoderConfig.timeEncoder` | string | `Time format for logs. The supported values are: "rfc3339nano", "rfc3339".` | `"rfc3339nano` |

### HTTP API

HTTP API of the client application service as defined [swagger](./pb/service.swagger.json).

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `api.http.enabled` | bool | `Enable the HTTP API.` | `true` |
| `api.http.cors.allowedOrigins` | []string | `Sets the allowed origins for CORS requests, as used in the 'Allow-Access-Control-Origin' HTTP header. Passing in a "*" will allow any domain.` | `"*"` |
| `api.http.cors.allowedHeaders` | []string | `Adds the provided headers to the list of allowed headers in a CORS request. This is an append operation so the headers Accept, Accept-Language, and Content-Language are always allowed. Content-Type must be explicitly declared if accepting Content-Types other than application/x-www-form-urlencoded, multipart/form-data, or text/plain.` | `"Accept","Accept-Language","Accept-Encoding","Content-Type","Content-Language","Content-Length","Origin","X-CSRF-Token","Authorization"` |
| `api.http.cors.allowedMethods` | []string | `Explicitly set allowed methods in the Access-Control-Allow-Methods header. This is a replacement operation so you must also pass GET, HEAD, and POST if you wish to support those methods.` | `"GET","PATCH","HEAD","POST","PUT","OPTIONS","DELETE"` |
| `api.http.cors.allowCredentials` | bool | `User agent may pass authentication details along with the request.` | `false` |
| `api.http.address` | string | `Listen specification <host>:<port> for http client connection.` | `"0.0.0.0:8080"` |
| `api.http.ui.enabled` | bool | `Set to true if you would like to run the web UI.` | `false` |
| `api.http.ui.directory` | string | `A path to the directory with web UI files. When it is not present, it creates <client_application_binary>/www with default ui.` | `""` |
| `api.http.tls.enabled` | bool | `Enable HTTPS.` | `false` |
| `api.http.tls.caPool` | string | `File path to the root certificate in PEM format which might contain multiple certificates in a single file.` |  `""` |
| `api.http.tls.keyFile` | string | `File path to private key in PEM format.` | `""` |
| `api.http.tls.certFile` | string | `File path to certificate in PEM format.` | `""` |
| `api.http.tls.clientCertificateRequired` | bool | `If true, require client certificate.` | `true` |

### gRPC API

gRPC API of the client application service as defined [service](./pb/service.proto).

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `api.grpc.enabled` | bool | `Enable the GRPC API.` | `true` |
| `api.grpc.address` | string | `Listen specification <host>:<port> for grpc client connection.` | `"0.0.0.0:8081"` |
| `api.grpc.enforcementPolicy.minTime` | string | `The minimum amount of time a client should wait before sending a keepalive ping. Otherwise the server close connection.` | `5s`|
| `api.grpc.enforcementPolicy.permitWithoutStream` | bool |  `If true, server allows keepalive pings even when there are no active streams(RPCs). Otherwise the server close connection.`  | `false` |
| `api.grpc.keepAlive.maxConnectionIdle` | string | `A duration for the amount of time after which an idle connection would be closed by sending a GoAway. 0s means infinity.` | `0s` |
| `api.grpc.keepAlive.maxConnectionAge` | string | `A duration for the maximum amount of time a connection may exist before it will be closed by sending a GoAway. 0s means infinity.` | `0s` |
| `api.grpc.keepAlive.maxConnectionAgeGrace` | string | `An additive period after MaxConnectionAge after which the connection will be forcibly closed. 0s means infinity.` | `0s` |
| `api.grpc.keepAlive.time` | string | `After a duration of this time if the server doesn't see any activity it pings the client to see if the transport is still alive.` | `2h` |
| `api.grpc.keepAlive.timeout` | string | `After having pinged for keepalive check, the client waits for a duration of Timeout and if no activity is seen even after that the connection is closed.` | `20s` |
| `api.grpc.tls.caPool` | string | `File path to the root certificate in PEM format which might contain multiple certificates in a single file.` |  `""` |
| `api.grpc.tls.enabled` | bool | `Enable TLS for grpc.` | `false` |
| `api.grpc.tls.keyFile` | string | `File path to private key in PEM format.` | `""` |
| `api.grpc.tls.certFile` | string | `File path to certificate in PEM format.` | `""` |
| `api.grpc.tls.clientCertificateRequired` | bool | `If true, require client certificate.` | `true` |

### Device client

The configuration sets up access to the devices via COAP protocol.

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `api.coap.maxMessageSize` | int | `Max message size which can be sent/received via coap. i.e. 256*1024 = 262144 bytes.` | `262144` |
| `api.coap.inactivityMonitor.timeout` | string | `Time limit to close inactive connection.` | `20s` |
| `api.coap.blockwiseTransfer.enabled` | bool | `If true, enable blockwise transfer of coap messages.` | `false` |
| `api.coap.blockwiseTransfer.blockSize` | int | `Size of blockwise transfer block.` | `1024` |
| `api.coap.tls.subjectUuid` | UUID | `Provides an identifier for client applications for establishing TLS connections or for devices that are set as owner devices` | `""` |
| `api.coap.tls.preSharedKeyUuid` | UUID | `Pre-shared key used in conjunction with subjectUUID to enable TLS connection.` | `""` |

> Note that the string type related to time (i.e. timeout, idleConnTimeout, expirationTime) is decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "1.5h" or "2h45m". Valid time units are "ns", "us", "ms", "s", "m", "h".
