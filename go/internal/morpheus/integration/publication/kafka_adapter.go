package publication

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/carlosgab83/matrix/go/internal/morpheus/domain"
	shared_domain "github.com/carlosgab83/matrix/go/internal/shared/domain"
	"github.com/carlosgab83/matrix/go/internal/shared/integration/logging"
)

type KafkaPublisher struct {
	Producer sarama.SyncProducer
	Logger   logging.Logger
}

func NewKafkaPublisher(cfg domain.Config, logger logging.Logger) (*KafkaPublisher, error) {
	brokers := []string{cfg.KafKaConsumerAddress}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaPublisher{
		Producer: producer,
		Logger:   logger,
	}, nil
}

func (kp *KafkaPublisher) NewDBPrice(ctx context.Context, price shared_domain.Price) error {
	b, err := json.Marshal(price)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: NewDBPriceTopic,
		Value: sarama.StringEncoder(b),
	}

	partition, offset, err := kp.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	kp.Logger.Debug("Message Sent", "partition", partition, "offset", offset)
	return nil
}

func (kp *KafkaPublisher) Close() error {
	return kp.Producer.Close()
}
