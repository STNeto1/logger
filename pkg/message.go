package pkg

import "encoding/json"

type Message struct {
	Topic string          `json:"topic"`
	Body  json.RawMessage `json:"body"`
}
