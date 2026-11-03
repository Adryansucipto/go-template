package config

import "strings"

type KafkaConfig struct {
	Enabled         bool   `yaml:"enabled"`
	Brokers         string `yaml:"brokers"`
	GroupID         string `yaml:"group_id"`
	MaxRetries      int    `yaml:"max_retries"`
	RetryBackoffMs  int    `yaml:"retry_backoff_ms"`
	DLQTopicSuffix  string `yaml:"dlq_topic_suffix"`
	OutboxBatchSize int    `yaml:"outbox_batch_size"`
}

func (k *KafkaConfig) BrokerList() []string {
	if strings.TrimSpace(k.Brokers) == "" {
		return []string{"localhost:9092"}
	}
	parts := strings.Split(k.Brokers, ",")
	var brokers []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}
	return brokers
}
