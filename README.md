# Lockout

Lockout is a basic micro service for checking against GEO-IP, using [GeoLite2](https://dev.maxmind.com/geoip/geoip2/geolite2/).

# How To Use

## Make
| Command | Description |
| ------- | ------------|
|| Run `test` `build` `package` |
| build | Builds the golang application |
| package | Package the build into a docker container |
| test | Run the unit test |

**NOTES**:
- `package` only works on *nix like systems, as it copies a cert file.

## K8s

I included a `deployment` and `service` in `lockout.k8.yaml`.

# Building Notes

The docker container needs a cert file, so either put one in the root directory of the project, or let the make file copy it into the root directory.

# Flags:
| Command | Type   | Default | Description                                            |
|---------|--------|---------|--------------------------------------------------------|
| api-host    | string | empty   | What host to bind the API to                   |
| api-port    | int    | `8080`  | What port to bind the API to                   |
| grpc-host    | string | empty   | What host to bind the gRPC to                   |
| grpc-port    | int    | `8082`  | What port to bind the gRPC to                   |


# Route(s):
| Route | Input | Output |
| ----- | ----- | ------ |
| `/api/v0/ip/whitelist` | See **Protobuf**:`IPCheckRequest` | See **Protobuf**:`IPCheckResponse` |

# Protobuf:
``` protobuf
syntax = "proto3";

package pb;

message IPCheckRequest {
  string ip = 1;
  repeated string countries = 2;
}

message IPCheckResponse {
  bool whitelisted = 1;
}

service WhitelistChecker {
  rpc CheckIP(IPCheckRequest) returns (IPCheckResponse) {}
}
```
