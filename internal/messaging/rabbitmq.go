package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *slog.Logger
}

func New(amqpURL string, logger *slog.Logger) (*RabbitMQ, error) {
	const op = "messaging.NewRabbitMQ"

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to RabbitMQ: %w", op, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("%s: failed to open channel: %w", op, err)
	}

	// Declare topic exchange
	err = ch.ExchangeDeclare(
		"orders",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("%s: failed to declare exchange: %w", op, err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		logger:  logger.With(slog.String("op", op)),
	}, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}

func (r *RabbitMQ) DeclareQueue(queueName string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQ) BindQueue(queueName, routingKey string) error {
	return r.channel.QueueBind(
		queueName,
		routingKey,
		"orders",
		false,
		nil,
	)
}

func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body []byte) error {
	const op = "messaging.RabbitMQ.Publish"

	err := r.channel.PublishWithContext(
		ctx,
		"orders",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
	if err != nil {
		r.logger.Error("failed to publish message", slog.String("routing_key", routingKey), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	r.logger.Info("message published", slog.String("routing_key", routingKey))
	return nil
}

func (r *RabbitMQ) Consume(queueName string, handler func(ctx context.Context, msg []byte) error) error {
	const op = "messaging.RabbitMQ.Consume"

	msgs, err := r.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%s: failed to start consumer: %w", op, err)
	}

	go func() {
		for d := range msgs {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := handler(ctx, d.Body); err != nil {
				r.logger.Error("message handling failed", slog.String("queue", queueName), slog.Any("error", err))
				d.Nack(false, true)
				continue
			}

			d.Ack(false)
			r.logger.Debug("message processed successfully", slog.String("queue", queueName))
		}
	}()

	r.logger.Info("consumer started", slog.String("queue", queueName))
	return nil
}
