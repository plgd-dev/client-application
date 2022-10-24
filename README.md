[![CI](https://github.com/plgd-dev/client-application/workflows/Test/badge.svg)](hhttps://github.com/plgd-dev/client-application/actions/workflows/test.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=plgd-dev_client-application&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=plgd-dev_client-application)
[![codecov](https://codecov.io/gh/plgd-dev/client-application/branch/main/graph/badge.svg)](https://codecov.io/gh/plgd-dev/client-application)

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

* `--config`: path to the config file
* `--version`: print the version of the client application

## Build

The build process uses goreleaser, so you will need to commit all changes and create a tag on the local machine.

```sh
git commit -a -m "my changes"
git tag -f v0.0.1-myversion
make build
```

## YAML Configuration

A configuration template is available on [config.yaml](./config.yaml).

### Logging

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `log.dumpBody` | bool | `Set to true if you would like to dump raw messages.` | `false` |
| `log.level` | string | `Logging enabled from level.` | `"info"` |
| `log.encoding` | string | `Logging format. The supported values are: "json", "console"` | `"console"` |
| `log.stacktrace.enabled` | bool | `Log stacktrace.` | `"false` |
| `log.stacktrace.level` | string | `Stacktrace from level.` | `"warn` |
| `log.encoderConfig.timeEncoder` | string | `Time format for logs. The supported values are: "rfc3339nano", "rfc3339".` | `"rfc3339nano` |

### HTTP API

HTTP API of the client application service as defined [swagger](./pb/service.swagger.json).

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `apis.http.enabled` | bool | `Enable the HTTP API.` | `true` |
| `apis.http.cors.allowedOrigins` | []string | `Sets the allowed origins for CORS requests, as used in the 'Allow-Access-Control-Origin' HTTP header. Passing in a "*" will allow any domain.` | `"*"` |
| `apis.http.cors.allowedHeaders` | []string | `Adds the provided headers to the list of allowed headers in a CORS request. This is an append operation so the headers Accept, Accept-Language, and Content-Language are always allowed. Content-Type must be explicitly declared if accepting Content-Types other than application/x-www-form-urlencoded, multipart/form-data, or text/plain.` | `"Accept","Accept-Language","Accept-Encoding","Content-Type","Content-Language","Content-Length","Origin","X-CSRF-Token","Authorization"` |
| `apis.http.cors.allowedMethods` | []string | `Explicitly set allowed methods in the Access-Control-Allow-Methods header. This is a replacement operation so you must also pass GET, HEAD, and POST if you wish to support those methods.` | `"GET","PATCH","HEAD","POST","PUT","OPTIONS","DELETE"` |
| `apis.http.cors.allowCredentials` | bool | `User agent may pass authentication details along with the request.` | `false` |
| `apis.http.address` | string | `Listen specification <host>:<port> for http client connection.` | `"0.0.0.0:8080"` |
| `apis.http.readTimeout` | string | `The maximum duration for reading the entire request, including the body by the server. A zero or negative value means there will be no timeout.` | `8s` |
| `apis.http.readHeaderTimeout` | string | `The amount of time allowed to read request headers by the server. If readHeaderTimeout is zero, the value of readTimeout is used. If both are zero, there is no timeout.` | `4s` |
| `apis.http.writeTimeout` | string | `The maximum duration before the server times out writing of the response. A zero or negative value means there will be no timeout.` | `16s` |
| `apis.http.idleTimeout` | string | `The maximum amount of time the server waits for the next request when keep-alives are enabled. If idleTimeout is zero, the value of readTimeout is used. If both are zero, there is no timeout.` | `30s` |
| `apis.http.ui.enabled` | bool | `Set to true if you would like to run the web UI.` | `false` |
| `apis.http.ui.directory` | string | `A path to the directory with web UI files. When it is not present, it creates <client_application_binary>/www with default ui.` | `""` |
| `apis.http.tls.enabled` | bool | `Enable HTTPS.` | `false` |
| `apis.http.tls.caPool` | string | `File path to the root certificate in PEM format which might contain multiple certificates in a single file.` |  `""` |
| `apis.http.tls.keyFile` | string | `File path to private key in PEM format.` | `""` |
| `apis.http.tls.certFile` | string | `File path to certificate in PEM format.` | `""` |
| `apis.http.tls.clientCertificateRequired` | bool | `If true, require client certificate.` | `true` |

### gRPC API

gRPC API of the client application service as defined [service](./pb/service.proto).

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `apis.grpc.enabled` | bool | `Enable the GRPC API.` | `true` |
| `apis.grpc.address` | string | `Listen specification <host>:<port> for grpc client connection.` | `"0.0.0.0:8081"` |
| `apis.grpc.enforcementPolicy.minTime` | string | `The minimum amount of time a client should wait before sending a keepalive ping. Otherwise the server close connection.` | `5s`|
| `apis.grpc.enforcementPolicy.permitWithoutStream` | bool |  `If true, server allows keepalive pings even when there are no active streams(RPCs). Otherwise the server close connection.`  | `false` |
| `apis.grpc.keepAlive.maxConnectionIdle` | string | `A duration for the amount of time after which an idle connection would be closed by sending a GoAway. 0s means infinity.` | `0s` |
| `apis.grpc.keepAlive.maxConnectionAge` | string | `A duration for the maximum amount of time a connection may exist before it will be closed by sending a GoAway. 0s means infinity.` | `0s` |
| `apis.grpc.keepAlive.maxConnectionAgeGrace` | string | `An additive period after MaxConnectionAge after which the connection will be forcibly closed. 0s means infinity.` | `0s` |
| `apis.grpc.keepAlive.time` | string | `After a duration of this time if the server doesn't see any activity it pings the client to see if the transport is still alive.` | `2h` |
| `apis.grpc.keepAlive.timeout` | string | `After having pinged for keepalive check, the client waits for a duration of Timeout and if no activity is seen even after that the connection is closed.` | `20s` |
| `apis.grpc.tls.caPool` | string | `File path to the root certificate in PEM format which might contain multiple certificates in a single file.` |  `""` |
| `apis.grpc.tls.enabled` | bool | `Enable TLS for grpc.` | `false` |
| `apis.grpc.tls.keyFile` | string | `File path to private key in PEM format.` | `""` |
| `apis.grpc.tls.certFile` | string | `File path to certificate in PEM format.` | `""` |
| `apis.grpc.tls.clientCertificateRequired` | bool | `If true, require client certificate.` | `true` |

### Device client

The configuration sets up access to the devices via COAP protocol.

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `apis.coap.maxMessageSize` | int | `Max message size which can be sent/received via coap. i.e. 256*1024 = 262144 bytes.` | `262144` |
| `apis.coap.inactivityMonitor.timeout` | string | `Time limit to close inactive connection.` | `20s` |
| `apis.coap.blockwiseTransfer.enabled` | bool | `If true, enable blockwise transfer of coap messages.` | `false` |
| `apis.coap.blockwiseTransfer.blockSize` | int | `Size of blockwise transfer block.` | `1024` |
| `apis.coap.ownershipTransfer.methods` | []string | `Allowed ownership transfer methods. The supported values are: "justWorks", "manufacturerCertificate".` | `"justWorks"` |
| `apis.coap.ownershipTransfer.manufacturerCertificate.tls.caPool` | string | `File path to the root certificate certificate in PEM format which might contain multiple certificates in a single file.` |  `""` |
| `apis.coap.ownershipTransfer.manufacturerCertificate.tls.keyFile` | string | `File path to certificate client application private key in PEM format.` | `""` |
| `apis.coap.ownershipTransfer.manufacturerCertificate.tls.certFile` | string | `File path to certificate client application certificate in PEM format.` | `""` |
| `apis.coap.tls.subjectUuid` | UUID | `Provides an identifier for client applications for establishing TLS connections or for devices that are set as owner devices` | `""` |
| `apis.coap.tls.preSharedKeyUuid` | UUID | `Pre-shared key used in conjunction with subjectUUID to enable TLS connection.` | `""` |

### Remote provisioning

The configuration sets up ownership and authorization of devices via the [remote provisioning mode](https://docs.plgd.dev/docs/device-to-device-client/client-initialization).

Supported modes:

* [Devices and single client](https://docs.plgd.dev/docs/device-to-device-client/client-initialization/#devices-and-single-client) configure `remoteProvisioning.mode` to `""`
* [User agent mediates CSR from the client](https://docs.plgd.dev/docs/device-to-device-client/client-initialization/#devices-and-single-client) configure `remoteProvisioning.mode` to `"userAgent"` and all other properties.

| Property | Type | Description | Default |
| ---------- | -------- | -------------- | ------- |
| `remoteProvisioning.mode` | string | `Provides remote provisioning mode. In "userAgent" mode all signing certificates goes through the user agent (browser,cli). "" means none. The supported values are: "", "userAgent"` | `""` |
| `remoteProvisioning.userAgent.certificateAuthorityAddress` | string | `Certificate authority server address in format {DNS}:{PORT}` | `""` |
| `remoteProvisioning.userAgent.csrChallengeStateExpiration` | string | `Defines how long is valid csr challenge.` | `"1m"` |
| `remoteProvisioning.authorization.authority` | string | `Authority is the address of the token-issuing authentication server.` | `""` |
| `remoteProvisioning.authorization.clientId` | string | `Client ID to exchange an authorization code for an access token.` | `""` |
| `remoteProvisioning.authorization.audience` | string | `Identifier of the API configured in your OAuth provider.` | `""` |
| `remoteProvisioning.authorization.scopes` | []string | `List of required scopes.` | `[]` |
| `remoteProvisioning.authorization.ownerClaim` | string | `Claim used to identify owner of the device.` | `"sub"` |

> Note that the string type related to time (i.e. timeout, idleConnTimeout, expirationTime) is decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "1.5h" or "2h45m". Valid time units are "ns", "us", "ms", "s", "m", "h".
