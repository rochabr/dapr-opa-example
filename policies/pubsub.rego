package httpapi.authz

default allow = false

# Define allowed topics per app
allowed_topics = {
    "publisher": ["orders"],
    "subscriber": ["orders"]
}

# Define denied topics per app
denied_topics = {
    "publisher": ["internal", "sensitive", "admin"],
    "subscriber": ["internal", "sensitive", "admin"]
}

# Helper to check topic permissions
is_allowed_topic(app_id, topic) {
    some i
    topic == allowed_topics[app_id][i]
    
    not topic_denied(app_id, topic)
}

# Helper to check if topic is denied
topic_denied(app_id, topic) {
    some i
    topic == denied_topics[app_id][i]
}

# Allow access to service endpoints
allow {
    input.method == "POST"
    input.path == "/v1.0/invoke/publisher/method/orders"
    input.dapr_app_id == "publisher"
    print("ALLOWED: orders endpoint access")
}

# Allow pub/sub publishing based on topic permissions
allow {
    input.method == "POST"
    startswith(input.path, "/v1.0/publish/pubsub/")
    topic := split(input.path, "/")[5]
    is_allowed_topic(input.dapr_app_id, topic)
    print("ALLOWED: publishing to topic", topic)
}

# Allow pub/sub subscription based on topic permissions
allow {
    input.method == "POST"
    path_parts := split(input.path, "/")
    count(path_parts) > 1
    topic := path_parts[count(path_parts)-1]
    is_allowed_topic(input.dapr_app_id, topic)
    print("ALLOWED: subscribing to topic", topic)
}

# Allow health checks
allow {
    input.method == "GET"
    input.path == "/healthz"
    print("ALLOWED: health check")
}