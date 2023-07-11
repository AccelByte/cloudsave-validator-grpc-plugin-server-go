# cloudsave-validator-grpc-plugin-server-go

`AccelByte Gaming Services` capabilities can be extended using custom functions implemented in a `gRPC server`. If configured, custom functions in the `gRPC server` will be called by `AccelByte Gaming Services` instead of the default function.

The `gRPC server` and the `gRPC client` can actually communicate directly. However, additional services are necessary to provide **security**, **reliability**, **scalability**, and **observability**. We call these services as `dependency services`. The [grpc-plugin-dependencies](https://github.com/AccelByte/grpc-plugin-dependencies) repository is provided as an example of what these `dependency services` may look like. It
contains a docker compose which consists of these `dependency services`.

> :warning: **grpc-plugin-dependencies is provided as example for local development purpose only:** The dependency services in the actual gRPC server deployment may not be exactly the same.

## Overview

This repository contains a `sample custom cloudsave validator gRPC server app` written in `Go`. It provides a simple custom cloudsave record validation function for cloudsave service in `AccelByte Gaming Services`.

This sample app also shows how this `gRPC server` can be instrumented for better observability.
It is configured by default to send metrics, traces, and logs to the observability `dependency services` in [grpc-plugin-dependencies](https://github.com/AccelByte/grpc-plugin-dependencies).


## Sample use case

### Use case 1: schema validation

Player record with key that has suffix `favourite_weapon` expect follows this schema:
```json
{
  "userId": "string,required",
  "favouriteWeaponType": "enum [SWORD, GUN], required",
  "favouriteWeapon": "string"
}
```

This simple app will demonstrate custom record validation above.
This app contains scenario when client send invalid JSON schema to cloudsave-validator.

When client send create game record with key using suffix `favourite_weapon` with invalid request body:
```json
{
  "foo": "bar"
}
```

cloudsave-validator will response
```json
{
  "isSuccess": false,
  "key": "string, key of record",
  "error": {
    "errorCode": 1,
    "errorMessage": "userid cannot be empty;favouriteWeaponType cannot be empty;favouriteWeapon cannot be empty"
  }
}
```

### Use case 2: custom validation logic

A game record with key that has suffix `daily_message` are expected to have following schema:
```json
{
  "message": "string,required",
  "title": "string,required",
  "availableOn": "time"
}
```

`cloudsave-validator` can be used to validate whether this game record are eligible to accessed by validating time in field `availableOn`.
When client send get game record request with time stamp is before `availableOn`, cloudsave-validator will return following response:
```json
{
  "isSuccess": false,
  "key": "string, key of record",
  "error": {
    "errorCode": 2,
    "errorMessage": "not accessible yet"
  }
}
```

## Prerequisites

Before starting, you will need the following.

1. Windows 10 WSL2 or Linux Ubuntu 20.04 with the following installed.

   a. bash

   b. curl

   c. docker v23.x

   d. docker-compose v2.x

   e. docker loki driver

      ```  
      docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
      ```

   f. make

   g. go v1.19

   h. git

   i. jq

   j. [ngrok](https://ngrok.com/)

2. A local copy of [grpc-plugin-dependencies](https://github.com/AccelByte/grpc-plugin-dependencies) repository.

   ```
   git clone https://github.com/AccelByte/grpc-plugin-dependencies.git
   ```

3. Access to `AccelByte Gaming Services` demo environment.

   a. Base URL: https://demo.accelbyte.io.

   b. [Create a Game Namespace](https://docs.accelbyte.io/esg/uam/namespaces.html#tutorials) if you don't have one yet. Keep the `Namespace ID`.

   c. [Create an OAuth Client](https://docs.accelbyte.io/guides/access/iam-client.html) with `confidential` client type. Keep the `Client ID` and `Client Secret`.

    - NAMESPACE:{namespace}:CLOUDSAVEGRPCSERVICE [READ]

## Setup

To be able to run this sample app, you will need to follow these setup steps.

1. Create a docker compose `.env` file by copying the content of [.env.template](.env.template) file.
2. Fill in the required environment variables in `.env` file as shown below.

   ```
   AB_BASE_URL=https://demo.accelbyte.io      # Base URL of AccelByte Gaming Services demo environment
   AB_CLIENT_ID='xxxxxxxxxx'                  # Use Client ID from the Prerequisites section
   AB_CLIENT_SECRET='xxxxxxxxxx'              # Use Client Secret from the Prerequisites section
   AB_NAMESPACE='xxxxxxxxxx'                  # Use Namespace ID from the Prerequisites section
   PLUGIN_GRPC_SERVER_AUTH_ENABLED=false      # Enable or disable access token and permission verification
   ```

   > :warning: **Keep PLUGIN_GRPC_SERVER_AUTH_ENABLED=false for now**: It is currently not
   supported by `AccelByte Gaming Services`, but it will be enabled later on to improve security. If it is
   enabled, the gRPC server will reject any calls from gRPC clients without proper authorization
   metadata.

## Building

To build this sample app, use the following command.

```
make build
```

## Running

To (build and) run this sample app in a container, use the following command.

```
docker-compose up --build
```

## Testing

### Functional Test in Local Development Environment

The custom function in this sample app can be tested locally using [sample grpc client](https://bitbucket.org/mrth0103/cloudsave-validator-client/src/master/).
