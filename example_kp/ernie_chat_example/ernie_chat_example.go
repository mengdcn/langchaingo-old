package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/schema"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ernie"
)

type mcache struct {
}

func (m mcache) Set(key string, value string) error {
	return nil
}
func (m mcache) Get(key string) (string, error) {
	return "", nil
}
func (m mcache) Expire(key string, seconds int) error {
	return nil
}

func main() {

	//rdb := redis.NewClient(&redis.Options{
	//	Addr:     "192.168.5.89:6379",
	//	Password: "",
	//	DB:       0,
	//})
	cache := mcache{}
	llm, err := ernie.NewChatWithCallback(callbacks.LogHandler{}, ernie.WithModelName(ernie.ModelNameERNIEBot), ernie.WithCache(cache))
	// note:
	// You would include ernie.WithAKSK(apiKey,secretKey) to use specific auth info.
	// You would include ernie.WithModelName(ernie.ModelNameERNIEBot) to use the ERNIE-Bot model.
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	messages := []schema.ChatMessage{
		schema.HumanChatMessage{Content: "以李白的口吻写一首诗"},
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
