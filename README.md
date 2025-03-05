# Oathkeeper Hydrator Mutator Demo

This project demonstrates an issue with Oathkeeper's hydrator mutator when trying to call endpoints via `host.docker.internal`.

## Issue Description

When Oathkeeper's hydrator mutator is configured to call an API endpoint on the host machine via `host.docker.internal`, it fails with an error similar to:

```
prohibited IP address: 0.250.250.254 is not a permitted destination (denied by: 0.0.0.0/8)
```

However, when the same endpoint is accessed directly within the Docker network (using the service name), it works correctly.

## Project Structure

- `main.go`: A simple Go web server with echo and introspection endpoints
- `Dockerfile`: Dockerfile for the Go web server
- `docker-compose.yml`: Docker Compose configuration that sets up both services
- `config/oathkeeper.yaml`: Oathkeeper configuration
- `config/rules.json`: Oathkeeper access rules

## Prerequisites

- Docker
- Docker Compose
- curl (for testing)

## Setup Instructions

1. Create the following directory structure:

```
oathkeeper-demo/
├── config/
│   ├── oathkeeper.yaml
│   └── rules.json
├── main.go
├── Dockerfile
├── docker-compose.yml
└── README.md
```

2. Copy the files from this repository into their respective locations.

3. Start the services:

```bash
docker-compose up -d
```

4. Verify both services are running:

```bash
docker-compose ps
```

## Reproducing the Issue

### Test Case 1: Hydrator Mutator with Direct Service Access (This Should Work)

```bash
curl -X POST http://localhost:4455/direct \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{"test":"data"}'
```

Expected result: The echo server should receive the request and echo back the JSON data.

### Test Case 2: Hydrator Mutator with Host Docker Internal (This Should Fail)

```bash
curl -X POST http://localhost:4455/host-docker-internal \
  -H "Authorization: Bearer token" \
  -H "Content-Type: application/json" \
  -d '{"test":"data"}'
```

Expected result: The request should fail with an error message in the Oathkeeper logs about prohibited IP address.

## Checking Logs

To check the Oathkeeper logs:

```bash
docker-compose logs oathkeeper
```

To check the echo server logs:

```bash
docker-compose logs echo-server
```

## Workaround

A potential workaround is to expose the service publicly (e.g., using ngrok) and use the public URL instead of `host.docker.internal`. However, this is not always desirable or practical.

## Additional Notes

- The echo server has a `/introspect` endpoint that returns a fixed response with `"active": true`. This simulates an OAuth introspection endpoint.
- Oathkeeper is configured to use this endpoint for authentication.
