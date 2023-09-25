package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms/chatglm"
	"log"
)

func main() {
	llm, err := chatglm.New(chatglm.WithEnableSearch(true), chatglm.WithSearchQuery(""))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	result, err := llm.Call(ctx, "今天天气怎么样")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	fmt.Printf("%#v\n", llm.GetUsage())
	//
	//llmChat, err := chatglm.NewChat()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//messages := []schema.ChatMessage{
	//	schema.HumanChatMessage{Content: "介绍一下你自己"},
	//}
	//completion, err := llmChat.Call(ctx, messages,
	//	llms.WithTemperature(0.8),
	//	llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
	//		log.Println(string(chunk))
	//		return nil
	//	}),
	//)
	//
	//if err != nil {
	//	log.Println("error")
	//	log.Fatal(err)
	//}
	//
	//// 同一个llmChat对象 并发处理回复时，可能会导致GetUsage方法不准确
	//log.Printf("%#v\n", llmChat.GetUsage())
	//
	//fmt.Printf("%v\n", completion)

	// 向量
	//emb, err := emb_chatglm.NewChatglm(emb_chatglm.WithClient(*llm))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//res, err := emb.EmbedDocuments(ctx, []string{"靠谱前程", "高考成绩"})
	//if err != nil {
	//	log.Println(err.Error())
	//	log.Fatal("emb.EmbedDocuments(ctx, []string{\"靠谱前程\", \"高考成绩\"}) error")
	//}
	//fmt.Println(res)
}
