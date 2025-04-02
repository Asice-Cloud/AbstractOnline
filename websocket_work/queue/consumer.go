package queue

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

type consumer struct {
	ctx      context.Context
	duration time.Duration
	channel  chan []string
	handler  handlerfunc
}

func defaultHandler(msg Message) {
	fmt.Println(msg)
}

func NewConsumer(ctx context.Context, handler handlerfunc) *consumer {
	return &consumer{
		ctx:      ctx,
		duration: time.Second,
		channel:  make(chan []string, 1000),
		handler:  handler,
	}
}

func (c *consumer) listen(redisClient *redis.Client, topic string) {
	go func() {
		for {
			select {
			case ret := <-c.channel:
				key := topic + HashSuffix
				result, err := redisClient.HMGet(c.ctx, key, ret...).Result()
				if err != nil {
					log.Println("Failed to HMGet from Redis:", err)
					continue
				}

				if len(result) > 0 {
					_, err := redisClient.HDel(c.ctx, key, ret...).Result()
					if err != nil {
						log.Println("Failed to HDel from Redis:", err)
					}
				}

				for _, v := range result {
					if v == nil {
						continue
					}
					var msg Message
					err := sonic.Unmarshal([]byte(v.(string)), &msg)
					if err != nil {
						log.Println("Failed to unmarshal message:", err)
						continue
					}

					go c.handler(msg)
				}
			}
		}
	}()

	ticker := time.NewTicker(c.duration)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			log.Println("Consumer quit:", c.ctx.Err())
			return
		case <-ticker.C:
			min := strconv.Itoa(0)
			max := strconv.Itoa(int(time.Now().Unix()))
			opt := &redis.ZRangeBy{
				Min: min,
				Max: max,
			}

			key := topic + SetSuffix
			result, err := redisClient.ZRangeByScore(c.ctx, key, opt).Result()
			if err != nil {
				log.Println("Failed to ZRangeByScore from Redis:", err)
				continue
			}

			if len(result) > 0 {
				_, err := redisClient.ZRemRangeByScore(c.ctx, key, min, max).Result()
				if err != nil {
					log.Println("Failed to ZRemRangeByScore from Redis:", err)
					continue
				}

				c.channel <- result
			}
		}
	}
}
