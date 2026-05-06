package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"Order/internal/models/order"
)

type OrderPublisher struct {
	rabbitMQ *RabbitMQ
}

func NewOrderPublisher(rabbitMQ *RabbitMQ) *OrderPublisher {
	return &OrderPublisher{
		rabbitMQ: rabbitMQ,
	}
}

func (p *OrderPublisher) OrderCreated(ctx context.Context, order order.Order) error {
	const op = "messaging.OrderPublisher.OrderCreated"

	body, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("%s: failed to marshal order: %w", op, err)
	}

	return p.rabbitMQ.Publish(ctx, "order.created", body)
}
