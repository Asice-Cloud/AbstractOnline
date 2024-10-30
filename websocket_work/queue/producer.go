package queue

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type producer struct {
	ctx context.Context
}

func NewProducer(ctx context.Context) *producer {
	return &producer{ctx: ctx}
}

func (p *producer) publish(rc *redis.Client, topic string, msg *Message) (int64, error) {
	rz := redis.Z{
		Score:  msg.GetTimeScore(),
		Member: msg.GetID(),
	}

	// set_key and hash_key
	sk := topic + SetSuffix
	hk := topic + HashSuffix

	// Add to sorted set
	n, err := rc.ZAdd(p.ctx, sk, rz).Result()
	if err != nil {
		log.Printf("Failed to add to sorted set: %v", err)
		return n, err
	}

	// Add to hash
	_, err = rc.HSet(p.ctx, hk, msg.GetID(), msg).Result()
	if err != nil {
		log.Printf("Failed to add to hash: %v", err)
		return n, err
	}

	log.Printf("Message published successfully: %v", msg)
	return n, nil
}
