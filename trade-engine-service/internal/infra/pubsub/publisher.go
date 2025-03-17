package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/francisco-alonso/trade-engine-service/internal/domain"
)

type Publisher struct {
	client  *pubsub.Client
	topic   *pubsub.Topic
}

func NewPublisher(projectId, topicId string) (*Publisher, error)  {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client: %w", err)
	}
	
	topic := client.Topic(topicId)
	
	return &Publisher{client: client, topic: topic}, nil
}

func (p *Publisher) Publish(ctx context.Context, order domain.Order) error {
	data, err := json.Marshal(order)

	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}
	
	
	result := p.topic.Publish(ctx, &pubsub.Message{Data: data})	
	
	_, err = result.Get(ctx)
		
	return err
}
