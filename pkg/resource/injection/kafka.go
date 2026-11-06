package injection

import (
	pkgConfig "go-template/pkg/config"
	"go-template/pkg/resource/model"
	"time"

	"github.com/segmentio/kafka-go"
)

func NewKafka(cfg pkgConfig.Config) (model.Kafka, error) {
	brokers := cfg.Kafka.BrokerList()
	groupID := cfg.Kafka.GroupID
	if groupID == "" {
		groupID = "go-template-consumer-group"
	}

	maxRetries := cfg.Kafka.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	backoffMs := cfg.Kafka.RetryBackoffMs
	if backoffMs <= 0 {
		backoffMs = 500
	}

	dlqSuffix := cfg.Kafka.DLQTopicSuffix
	if dlqSuffix == "" {
		dlqSuffix = ".dlq"
	}

	var writer *kafka.Writer
	if cfg.Kafka.Enabled {
		writer = &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			RequiredAcks:           kafka.RequireOne,
			AllowAutoTopicCreation: true,
			WriteTimeout:           10 * time.Second,
			ReadTimeout:            10 * time.Second,
		}
	}

	return model.Kafka{
		Enabled:        cfg.Kafka.Enabled,
		Brokers:        brokers,
		GroupID:        groupID,
		MaxRetries:     maxRetries,
		RetryBackoffMs: backoffMs,
		DLQTopicSuffix: dlqSuffix,
		Writer:         writer,
	}, nil
}
