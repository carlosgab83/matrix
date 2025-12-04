package reception

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
	"github.com/carlosgab83/matrix/go/internal/tank/domain"
)

type KafkaReceptor struct {
	Config         domain.Config
	Logger         logging.Logger
	Topic          string
	Receive        chan domain.NotificationPayload
	SaramaConsumer sarama.ConsumerGroup
	Ctx            context.Context
}

const (
	TopicName string = "notification.new"
	GroupID   string = "price-analyzer"
)

func NewKafkaReceptor(ctx context.Context, cfg domain.Config, logger logging.Logger) (*KafkaReceptor, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup([]string{cfg.KafkaProducerAddress}, GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("error creating kafka consumer group: %v", err)
	}

	return &KafkaReceptor{
		Config:         cfg,
		Logger:         logger,
		Topic:          TopicName,
		Receive:        make(chan domain.NotificationPayload),
		SaramaConsumer: consumer,
		Ctx:            ctx,
	}, nil
}

func (kr *KafkaReceptor) Setup(s sarama.ConsumerGroupSession) error {
	kr.Logger.Info("Kafka ConsumerGroupSession Setup done")
	return nil
}

func (kr *KafkaReceptor) Cleanup(s sarama.ConsumerGroupSession) error {
	// TODO
	kr.Logger.Info("Kafka ConsumerGroupSession Cleanup done")
	return nil
}

func (kr *KafkaReceptor) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		kr.Logger.Info("Message received",
			"topic", message.Topic,
			"partition", message.Partition,
			"offset", message.Offset,
			"value", string(message.Value))

		var notificationPayload domain.NotificationPayload
		if err := json.Unmarshal(message.Value, &notificationPayload); err != nil {
			kr.Logger.Error("error ConsumeClaim Unmarshaling: %v", err)
			return err
		}

		select {
		case kr.Receive <- notificationPayload:
			// NOTE: Assuming the message always is received
			session.MarkMessage(message, "")
		case <-kr.Ctx.Done():
			return kr.Ctx.Err()
		}
	}

	return nil
}

func (kr *KafkaReceptor) BeginConsumption() error {
	for {
		// The Consume method should be called in a loop because it returns
		// if a recoverable error occurs or when a rebalance completes.

		err := kr.SaramaConsumer.Consume(kr.Ctx, []string{TopicName}, kr)
		if err != nil {
			// The Sarama.ErrClosedConsumerGroup error occurs if the context is canceled
			if err != sarama.ErrClosedConsumerGroup {
				kr.Logger.Error("Error consuming", "error", err)
			}
		}

		// If the context was canceled (for example by a signal), exit.
		if kr.Ctx.Err() != nil {
			return kr.Ctx.Err()
		}
	}
}

func (kr *KafkaReceptor) ReceiveCh() chan domain.NotificationPayload {
	return kr.Receive
}

func (kr *KafkaReceptor) Close() error {
	return kr.SaramaConsumer.Close()
}
