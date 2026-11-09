package event

import (
	"encoding/json"
	"time"
)

type CloudEvent struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

func NewCloudEvent(id, eventType, source string, payload any) (CloudEvent, []byte, error) {
	rawData, err := json.Marshal(payload)
	if err != nil {
		return CloudEvent{}, nil, err
	}

	evt := CloudEvent{
		ID:        id,
		Type:      eventType,
		Source:    source,
		Timestamp: time.Now().UTC(),
		Data:      rawData,
	}

	bytes, err := json.Marshal(evt)
	if err != nil {
		return CloudEvent{}, nil, err
	}

	return evt, bytes, nil
}
