package queue

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var RedisQueue *Queue
var Redis *redis.Client

func InitRedisQueue() {
	RedisQueue = NewQueue(context.Background(), Redis,
		WithTopic("send-message"),
		WithHandler(CreateAndSendMessages),
	)
	RedisQueue.Start()
}

func CreateAndSendMessages(msg Message) {
	id := uuid.New().String()
	logMsg := NewMsg(id, time.Now(), map[string]interface{}{"user_id": 123, "action": "login"})
	_, err := RedisQueue.Publish(logMsg)
	if err != nil {
		fmt.Println("send ", err)
	}
}

func Test(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Fatalf("redis connect error: %v", err)
	}
	Redis = client
	InitRedisQueue()

	// Create a test message
	testMsg := NewMsg(uuid.New().String(), time.Now(), map[string]interface{}{"test": "message"})
	_, err = RedisQueue.Publish(testMsg)
	if err != nil {
		t.Fatalf("failed to publish message: %v", err)
	}
	log.Printf("Message published successfully: %v", testMsg)

	// Add a delay to ensure the message is processed
	time.Sleep(2 * time.Second)

	// Verify the message was published
	key := "send-message" + SetSuffix
	ids, err := Redis.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(time.Now().Unix(), 10),
	}).Result()
	if err != nil {
		t.Fatalf("failed to retrieve message IDs: %v", err)
	}

	if len(ids) == 0 {
		t.Fatalf("no messages found in queue")
	}

	log.Printf("Message IDs found in queue: %v", ids)

	// Retrieve and unmarshal the messages
	hashKey := "send-message" + HashSuffix
	messages, err := Redis.HMGet(context.Background(), hashKey, ids...).Result()
	if err != nil {
		t.Fatalf("failed to retrieve messages: %v", err)
	}

	var readableMessages []Message
	for _, msg := range messages {
		if msg == nil {
			continue
		}
		var message Message
		err := sonic.Unmarshal([]byte(msg.(string)), &message)
		if err != nil {
			t.Fatalf("failed to unmarshal message: %v", err)
		}
		readableMessages = append(readableMessages, message)
	}

	// Format the result as JSON
	jsonResult, err := sonic.MarshalIndent(readableMessages, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal messages: %v", err)
	}

	fmt.Println("Test passed, messages found in queue:", string(jsonResult))
}
