package httpapi.authz

# Allow service-to-service invocation
allow {
    input.method == "POST"
    startswith(input.path, "/orders")
    input.dapr_app_id == "publisher"
}

# Allow subscriber to receive messages
allow {
    input.method == "POST"
    input.path == "/orders"
    input.dapr_app_id == "subscriber"
}