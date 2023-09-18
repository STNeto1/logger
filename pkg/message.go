package pkg

import "encoding/json"

type Message struct {
	ID        int64           `json:"id,omitempty" db:"id"`
	Topic     string          `json:"topic" db:"topic"`
	Body      json.RawMessage `json:"body" db:"data"`
	CreatedAt string          `json:"created_at,omitempty" db:"created_at"`
}
