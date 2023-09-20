package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/chatglm"
	"github.com/tmc/langchaingo/schema"
	"log"
)

func main() {
	llm, err := chatglm.New()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	result, err := llm.Call(ctx, "介绍一下你自己")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	fmt.Printf("%#v\n", llm.GetUsage())

	llmChat, err := chatglm.NewChat()
	if err != nil {
		log.Fatal(err)
	}
	messages := []schema.ChatMessage{
		schema.HumanChatMessage{Content: "介绍一下你自己"},
	}
	completion, err := llmChat.Call(ctx, messages,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			log.Println(string(chunk))
			return nil
		}),
	)

	if err != nil {
		log.Println("error")
		log.Fatal(err)
	}

	// 同一个llmChat对象 并发处理回复时，可能会导致GetUsage方法不准确
	log.Printf("%#v\n", llmChat.GetUsage())

	fmt.Printf("%v\n", completion)

}
