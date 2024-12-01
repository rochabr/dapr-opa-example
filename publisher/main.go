package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/client"
	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

type Order struct {
	OrderID  string `json:"orderId"`
	Customer string `json:"customer"`
}

var daprClient client.Client

func main() {
	var err error
	daprClient, err = client.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}
	defer daprClient.Close()

	s := daprd.NewService(":3000")

	// Add handlers for different topics
	if err := s.AddServiceInvocationHandler("/orders", handleOrders); err != nil {
		log.Fatalf("Failed to add orders handler: %v", err)
	}

	if err := s.AddServiceInvocationHandler("/internal", handleInternal); err != nil {
		log.Fatalf("Failed to add internal handler: %v", err)
	}

	if err := s.AddServiceInvocationHandler("/customers", handleCustomers); err != nil {
		log.Fatalf("Failed to add internal handler: %v", err)
	}

	if err := s.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start service: %v", err)
	}
}

func handleOrders(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	log.Printf("Received order request: %s", string(in.Data))

	var order Order
	if err := json.Unmarshal(in.Data, &order); err != nil {
		log.Printf("Error unmarshaling order: %v", err)
		return nil, err
	}

	if err := daprClient.PublishEvent(ctx, "pubsub", "orders", order); err != nil {
		log.Printf("Error publishing order: %v", err)
		return nil, err
	}

	log.Printf("Successfully published to orders topic - ID: %s", order.OrderID)

	return &common.Content{
		Data:        []byte("Order submitted successfully"),
		ContentType: "text/plain",
	}, nil
}

func handleCustomers(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	log.Printf("Received customer request: %s", string(in.Data))

	return &common.Content{
		Data:        []byte("Customer endpoint allowed"),
		ContentType: "text/plain",
	}, nil
}

func handleInternal(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	log.Printf("Attempting to publish to internal topic")

	data := map[string]string{
		"message": "internal data",
	}

	if err := daprClient.PublishEvent(ctx, "pubsub", "internal", data); err != nil {
		log.Printf("Error publishing to internal topic: %v", err)
		return nil, err
	}

	return &common.Content{
		Data:        []byte("Internal message published"),
		ContentType: "text/plain",
	}, nil
}
