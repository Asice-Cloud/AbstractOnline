package queue

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

const (
	HashSuffix = "_hash"
	SetSuffix  = "_set"
)

var once sync.Once

type Queue struct {
	ctx context.Context

	rc    *redis.Client
	topic string

	producer *producer
	consumer *consumer
}

func NewQueue(ctx context.Context, redis *redis.Client, opt ...Option) *Queue {
	var q *Queue
	once.Do(func() {
		defauleOpt := Options{
			topic:   "asice cloud",
			handler: defaultHandler,
		}

		for _, o := range opt {
			o(&defauleOpt)
		}

		q = &Queue{
			ctx:      ctx,
			rc:       redis,
			topic:    defauleOpt.topic,
			producer: NewProducer(ctx),
			consumer: NewConsumer(ctx, defauleOpt.handler),
		}
	})
	return q
}

func (q *Queue) Start() {
	go q.consumer.listen(q.rc, q.topic)
}

func (q *Queue) Publish(msg *Message) (int64, error) {
	return q.producer.publish(q.rc, q.topic, msg)
}
