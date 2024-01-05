package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms/minimax"
)

func main() {

	llmmini, err := minimax.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//result, err := llmmini.Call(context.Background(), "介绍一下你自己")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Printf("completion:%s\n", result)
	//fmt.Println(llmmini.GetUsage())

	emResult, err := llmmini.CreateEmbedding(context.Background(), []string{"介绍一下你自己"})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("emResult", emResult)

	emResult, err = llmmini.CreateDbEmbedding(context.Background(), []string{"介绍一下你自己"})
	//fmt.Println("db", emResult)
	fmt.Println("usage1:", llmmini.GetUsage())

	emResult, err = llmmini.CreateQueryEmbedding(context.Background(), []string{"介绍一下你自己"})
	//fmt.Println("query", emResult)
	fmt.Println("usage2:", llmmini.GetUsage())

	//llmChat ,err := minimax.NewChat()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//resultChat, err := llmChat.
}
