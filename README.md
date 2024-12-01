# Dapr with Standalone OPA Example

This example demonstrates using Dapr with a standalone OPA instance for policy enforcement.

## Prerequisites

- Dapr CLI installed
- Docker installed
- Go 1.19 or later

## Setup Steps

1. Start OPA container:
```bash
docker run -d --name opa -p 8181:8181 openpolicyagent/opa:latest run --server --addr :8181
```

2. Load policies into OPA:
```bash
curl -X PUT localhost:8181/v1/policies/pubsub --data-binary @policies/pubsub.rego
curl -X PUT localhost:8181/v1/policies/service --data-binary @policies/service.rego
```

3. Start the publisher service:
```bash
cd publisher
dapr run --app-id publisher --app-port 3000 --dapr-http-port 3500 --config ../config/config.yaml go run main.go
```

4. Start the subscriber service:
```bash
cd subscriber
dapr run --app-id subscriber --app-port 3001 --dapr-http-port 3501 --config ../config/config.yaml go run main.go
```

## Testing

1. Publish a message (allowed):
```bash
curl -X POST http://localhost:3500/v1.0/publish/pubsub/orders \
  -H "Content-Type: application/json" \
  -d '{"orderId": "123", "customer": "example"}'
```

2. Try publishing to restricted topic (denied):
```bash
curl -X POST http://localhost:3500/v1.0/publish/pubsub/internal \
  -H "Content-Type: application/json" \
  -d '{"data": "restricted"}'
```
