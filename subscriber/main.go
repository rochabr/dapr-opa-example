package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

type Order struct {
	OrderID  string `json:"orderId"`
	Customer string `json:"customer"`
}

func main() {
	s := daprd.NewService(":3001")

	if err := s.AddTopicEventHandler(&common.Subscription{
		PubsubName: "pubsub",
		Topic:      "orders",
		Route:      "/orders",
	}, orderHandler); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start service: %v", err)
	}
}

func orderHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	// Log the raw data for debugging
	log.Printf("Received raw data: %+v", e.Data)

	// Convert the data to JSON
	jsonData, err := json.Marshal(e.Data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return false, err
	}

	var order Order
	if err := json.Unmarshal(jsonData, &order); err != nil {
		log.Printf("Error unmarshaling order: %v, Data: %s", err, string(jsonData))
		return false, err
	}

	log.Printf("Successfully processed order - ID: %s, Customer: %s", order.OrderID, order.Customer)
	return false, nil
}
