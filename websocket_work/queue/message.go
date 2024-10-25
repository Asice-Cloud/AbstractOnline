package queue

import (
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Message struct {
	ID          string      `json:"id""`
	CreateTime  time.Time   `json:"create_time"`
	ConsumeTime time.Time   `json:"consume_time"`
	Body        interface{} `json:"body"`
}

func NewMsg(ID string, consumeTime time.Time, body interface{}) *Message {
	return &Message{
		ID: func(id string) string {
			if id == "" {
				id = uuid.New().String()
			}
			return id
		}(ID),
		CreateTime:  time.Now(),
		ConsumeTime: consumeTime,
		Body:        body,
	}
}

func (m *Message) GetTimeScore() float64 {
	return float64(m.ConsumeTime.Unix())
}

func (m *Message) GetID() string {
	return m.ID
}

func (m *Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
