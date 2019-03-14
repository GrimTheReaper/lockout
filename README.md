# Lockout

Lockout is a basic micro service for checking against GEO-IP, using [GeoLite2](https://dev.maxmind.com/geoip/geoip2/geolite2/).

# How To Use

## Make
| Command | Description |
| ------- | ------------|
|| Run `test` `build` `package` |
| build | Builds the golang application |
| buildtiny | If you have `upx` installed, you may run this to make the resulting binary ~2MB. Useful for a tiny micro-service. Does other things to slim down the binary |
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

# License
```
Copyright 2019 GrimTheReaper

Licensed under the Apache License, Version 2.0 (the "License");
you may not use these files except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
