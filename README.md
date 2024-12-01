# Dapr OPA Authorization Example

This example demonstrates how to use Open Policy Agent (OPA) with Dapr to implement authorization for service invocation.

## Prerequisites

- Docker 
- Dapr CLI
- Go 1.19 or later

## Project Structure

```
dapr-opa-example/
├── README.md
├── config/
│   └── config.yaml             # Dapr configuration referencing the OPA middleware
├── components/
│   └── middleware.yaml         # OPA middleware component definition
├── publisher/
│   ├── main.go                 # Service that exposes endpoints to be authorized
│   └── go.mod
└── subscriber/
    ├── main.go                 # Service that consumes messages
    └── go.mod
```

## Setup Steps

1. Start OPA container:

```bash
docker run -d --name opa --platform linux/arm64 -p 8181:8181 openpolicyagent/opa:latest run --server --addr :8181
```

2. Initialize Go modules:

```bash
cd publisher && go mod init publisher && go mod tidy
cd ../subscriber && go mod init subscriber && go mod tidy
```

3. Start the publisher service:

```bash
cd publisher
dapr run --app-id publisher \
         --app-port 3000 \
         --dapr-http-port 3500 \
         --resources-path ../components \
         --config ../config/config.yaml \
         go run main.go
```

4. Start the subscriber service in another terminal:

```bash
cd subscriber
dapr run --app-id subscriber \
         --app-port 3001 \
         --dapr-http-port 3501 \
         --resources-path ../components \
         --config ../config/config.yaml \
         go run main.go
```

## Configuration Files

### OPA Middleware Component (components/middleware.yaml)

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: opa-middleware
spec:
  type: middleware.http.opa
  version: v1
  metadata:
    - name: readBody
      value: "true"
    - name: rego
      value: |
        package http

        default allow = false

        # Allow service invocation to orders endpoint
        allow {
            input.request.method == "POST"
            input.request.path == "/v1.0/invoke/publisher/method/orders"
            trace(sprintf("Allowing access to orders endpoint: %v", [input.request.path]))
        }

        # Allow health checks
        allow {
            input.request.method == "GET"
            input.request.path == "/healthz"
        }
```

### Dapr Configuration (config/config.yaml)

```yaml
apiVersion: dapr.io/v1alpha1
kind: Configuration
metadata:
  name: dapr-config
spec:
  httpPipeline:
    handlers:
    - name: opa-middleware
      type: middleware.http.opa
```

## Testing

1. Test allowed endpoint (should succeed):

```bash
curl -X POST http://localhost:3500/v1.0/invoke/publisher/method/orders \
  -H "Content-Type: application/json" \
  -d '{"orderId": "123", "customer": "example"}'
```

2. Test unauthorized endpoint (should fail):

```bash
curl -X POST http://localhost:3500/v1.0/invoke/publisher/method/internal \
  -H "Content-Type: application/json" \
  -d '{"message": "secret"}'
```

## Understanding the Policy

The OPA policy in this example:

- Allows POST requests to the `/v1.0/invoke/publisher/method/orders` endpoint
- Allows GET requests to `/healthz` for health checks
- Denies all other requests by default

## Clean Up

To stop and clean up:

```bash
# Stop the Dapr applications using Ctrl+C
# Remove the OPA container
docker rm -f opa
```