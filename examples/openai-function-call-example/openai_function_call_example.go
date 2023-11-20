package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

func main() {
	baseUrl := "https://apiagent.kaopuai.com/v1"
	llm, err := openai.NewChat(openai.WithModel("gpt-4"), openai.WithBaseURL(baseUrl))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	completion, err := llm.Call(ctx, []schema.ChatMessage{
		schema.HumanChatMessage{Content: "湖北税局状态?"},
	}, llms.WithFunctions(functions))
	if err != nil {
		log.Fatal(err)
	}

	if completion.FunctionCall != nil {
		fmt.Printf("Function call: %v\n", completion.FunctionCall)
	}
}

func getCurrentWeather(location string, unit string) (string, error) {
	weatherInfo := map[string]interface{}{
		"location":    location,
		"temperature": "72",
		"unit":        unit,
		"forecast":    []string{"sunny", "windy"},
	}
	b, err := json.Marshal(weatherInfo)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var functions = []llms.FunctionDefinition{
	{
		Name:        "get_shuiju_status", // name 必须为英文 ^[a-zA-Z0-9_-]{1,64}$
		Description: "根据输入的地区，获取各地区税务局的状态",
		Parameters:  json.RawMessage(`{"type": "object", "properties": {"location": {"type": "string", "description": "地区，税务局的所属地区，例如：北京市、湖北省"}, "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}}, "required": ["location"]}`),
	},
}
