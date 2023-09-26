package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms/qwen"
	"log"
)

func main() {
	ctx := context.Background()

	// ===================== 单prompt ========================
	llm, err := qwen.New(qwen.WithEnableSearch(true))
	if err != nil {
		log.Fatal(err)
	}
	result, err := llm.Call(ctx, "办公逸创始人是谁")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
	fmt.Printf("%#v\n", llm.GetUsage())

	//// ===================== 多prompt ========================
	//res, err := llm.Generate(ctx, []string{"介绍一下你自己", "介绍一下办公逸"})
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//fmt.Println(res[0].Text)
	//fmt.Println(res[1].Text)
	//fmt.Printf("%#v\n", llm.GetUsage())

	//// ===================== chat ========================
	//llmChat, err := qwen.NewChat()
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

	// ===================== 向量 ========================
	//llmEmb, err := qwen.New()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//emb, err := emb_qwen.NewQwen(emb_qwen.WithClient(*llmEmb))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//resemb, err := emb.EmbedQuery(ctx, "高考成绩")
	//if err != nil {
	//	log.Println(err.Error())
	//	log.Fatal("emb.EmbedDocuments(ctx, []string{\"靠谱前程\", \"高考成绩\"}) error")
	//}
	//fmt.Println(len(resemb))
	//
	//resemb2, err := emb.EmbedDocuments(ctx, []string{"高考成绩", "介绍一下自己"})
	//if err != nil {
	//	log.Fatal("resemb2")
	//}
	//fmt.Println(len(resemb2))
}