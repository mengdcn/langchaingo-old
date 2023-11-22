package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ernie"
	"github.com/tmc/langchaingo/schema"
	"log"
	"os"
	"time"
)

type mcache struct {
	rdb *redis.Client
}

func (m mcache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	cmd := m.rdb.Set(ctx, key, value, expiration)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (m mcache) Get(ctx context.Context, key string) (string, error) {
	return m.rdb.Get(ctx, key).Result()
}

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.5.89:6379",
		Password: "",
		DB:       0,
	})

	cache := mcache{
		rdb: rdb,
	}
	k := os.Getenv("ERNIE_API_KEY")
	v := os.Getenv("ERNIE_SECRET_KEY")

	llm, err := ernie.NewChatWithCallback(callbacks.LogHandler{}, ernie.WithAKSK(k, v), ernie.WithModelName(ernie.ModelNameERNIEBot4), ernie.WithCache(cache))
	// note:
	// You would include ernie.WithAKSK(apiKey,secretKey) to use specific auth info.
	// You would include ernie.WithModelName(ernie.ModelNameERNIEBot) to use the ERNIE-Bot model.
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	messages := []schema.ChatMessage{
		schema.HumanChatMessage{Content: "介绍一下你自己"},
	}
	completion, err := llm.Call(ctx, messages,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			log.Println(string(chunk))
			return nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=================")
	fmt.Printf("%#v", llm.GetUsage())
	fmt.Printf("%#v", completion)

}
