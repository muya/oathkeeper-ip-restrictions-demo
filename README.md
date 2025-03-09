# Oathkeeper Hydrator Mutator Demo

This project demonstrates an issue with Oathkeeper's hydrator mutator when trying to call endpoints via `host.docker.internal`.

**This issue seems to be affecting Oathkeeper versions 0.40.8 and later (last tested working version is 0.40.6)**

## Issue Description

When Oathkeeper's hydrator mutator is configured to call an API endpoint on the host machine via `host.docker.internal`, it fails with an error similar to:

```
prohibited IP address: 0.250.250.254 is not a permitted destination (denied by: 0.0.0.0/8)
```

However, when the same endpoint is accessed directly within the Docker network (using the service name), it works correctly.

This issue seems to only affect versions 0.40.7 and later

## Project Structure

- `main.go`: A simple Go web server with echo and hydrate endpoints
- `Dockerfile-backend`: Dockerfile for the Go web server
- `docker-compose.yml`: Docker Compose configuration that sets up both services
- `config.yaml`: Oathkeeper configuration
- `rules.yaml`: Oathkeeper access rules
- `bruno-oathkeeper-api-requests`: A [Bruno](https://usebruno.com) API request collection (optional)

This project brings up 4 containers via docker-compose:

1. The simple backend which provides endpoints to test against - exposes port 9090
2. 3 versions of Oathkeeper (0.40.6, 0.40.8, 0.40.9)
    1. 0.40.6 - exposes port 7456
    2. 0.40.8 - exposes port 4456
    3. 0.40.9 - exposes port 8456

## Prerequisites

- Docker
- Docker Compose
- curl (for testing)


## Setup Instructions

1. Clone this repository

```
git clone git@github.com:muya/oathkeeper-ip-restrictions-demo.git
```

2. Bring up the containers via docker-compose

```shell
cd oathkeeper-ip-restrictions-demo

docker-compose up
```

3. Optional: import the bruno API request collection into Bruno. Note: `curl` works too.

(We'll leave them running in the foreground to observe output)

## Reproducing the Issue

To reproduce the issue, we'll be making calls to the decisions API.

We have defined 2 rules that call the hydrator mutator:

1. `test-ip-restriction-intra-docker-api-call`: Calls the `/hydrate` endpoint via intra-docker domain name (i.e. direct to `http://echo-server:9090/hydrate`)
2. `test-ip-restriction-host-network-api-call`: Calls the `/hydrate` endpoint via the host (i.e. `http://host.docker.internal:9090/hydrate`)

For all cases, the expected response is an HTTP 200 OK.

The `curl` commands below should be run in a separate terminal.

### Step 1: Test Version 0.40.6

#### Make call to the host-call endpoint

```shell
curl -i --request POST \
  --url http://127.0.0.1:7456/decisions/host-call
```

Expected: an HTTP 200 OK response

Actual: an HTTP 200 OK response


#### Make call to the intra-docker endpoint

```shell
curl -i --request POST \
  --url http://127.0.0.1:7456/decisions/intra-docker
```

Expected: an HTTP 200 OK response

Actual: an HTTP 200 OK response


### Step 2: Test Version 0.40.8

#### Make call to the host-call endpoint

```shell
curl -i --request POST \
  --url http://127.0.0.1:4456/decisions/host-call
```

Expected: an HTTP 200 OK response

Actual: an HTTP 500 Internal Server Error response, with error message saying: 

```
dial tcp 0.250.250.254:9090: prohibited IP address: 0.250.250.254 is not a permitted destination (denied by: 0.0.0.0/8)
```


#### Make call to the intra-docker endpoint

```shell
curl -i --request POST \
  --url http://127.0.0.1:4456/decisions/intra-docker
```

Expected: an HTTP 200 OK response

Actual: an HTTP 200 OK response


### Step 3: Test Version 0.40.9

#### Make call to the host-call endpoint

```shell
curl -i --request POST \
  --url http://127.0.0.1:8456/decisions/host-call
```

Expected: an HTTP 200 OK response

Actual: an HTTP 500 Internal Server Error response, with error message saying: 

```
dial tcp 0.250.250.254:9090: prohibited IP address: 0.250.250.254 is not a permitted destination (denied by: 0.0.0.0/8)
```


#### Make call to the intra-docker endpoint

```shell
curl -i --request POST \
  --url http://127.0.0.1:8456/decisions/intra-docker
```

Expected: an HTTP 200 OK response

Actual: an HTTP 200 OK response


## Workaround

A potential workaround is to expose the service publicly (e.g., using ngrok) and use the public URL instead of `host.docker.internal`. However, this is not always desirable or practical.
