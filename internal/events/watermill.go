package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	appConfig "github.com/shapnesht/ecommerce-platform-go/internal/config"
	"github.com/shapnesht/ecommerce-platform-go/internal/providers"
)

type EventPublisher struct {
	publisher message.Publisher
	queueName string
}

func (ep *EventPublisher) Publish(eventName string, payload interface{}, metadata map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)
	msg.Metadata.Set("event_name", eventName)
	for key, value := range metadata {
		msg.Metadata.Set(key, value)
	}

	return ep.publisher.Publish(ep.queueName, msg)
}

func (ep *EventPublisher) Close() error {
	return ep.publisher.Close()
}

func NewEventPublisher(ctx context.Context, cfg *appConfig.AWSConfig) (*EventPublisher, error) {
	logger := watermill.NewStdLogger(false, false)

	awsConfig, err := providers.CreateAWSConfig(ctx, cfg.S3Endpoint, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("unable to create aws config: %w", err)
	}

	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &EventPublisher{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}, nil
}
