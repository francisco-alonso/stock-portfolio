package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/francisco-alonso/stock-portfolio/portfolio-service/domain"
)

type Subscriber struct {
    client     *pubsub.Client
    subscription *pubsub.Subscription
    db         *PostgresDB
}

func NewSubscriber(ctx context.Context, projectID, subscriptionID string, db *PostgresDB) (*Subscriber, error) {
    client, err := pubsub.NewClient(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("failed to create Pub/Sub client: %w", err)
    }

    subscription := client.Subscription(subscriptionID)
    return &Subscriber{client: client, subscription: subscription, db: db}, nil
}

func (s *Subscriber) StartListening(ctx context.Context) {
    err := s.subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
        var order domain.Position
        if err := json.Unmarshal(msg.Data, &order); err != nil {
            log.Printf("Error parsing order: %v", err)
            msg.Nack()
            return
        }

        err := s.db.AddPosition(order.UserID, order.Asset, order.Quantity, order.Price)
        if err != nil {
            log.Printf("Error saving position: %v", err)
            msg.Nack()
            return
        }

        msg.Ack()
    })
    if err != nil {
        log.Fatalf("Error receiving messages: %v", err)
    }
}
