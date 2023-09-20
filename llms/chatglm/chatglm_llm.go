package chatglm

import (
	"context"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/chatglm/internal/chatglm_client"
	"github.com/tmc/langchaingo/schema"
)

type LLM struct {
	CallbacksHandler callbacks.Handler
	client           *chatglm_client.Client
}

var (
	_ llms.LLM           = (*LLM)(nil)
	_ llms.LanguageModel = (*LLM)(nil)
)

func New(opts ...Option) (*LLM, error) {
	c, err := newClient(opts...)
	return &LLM{
		client: c,
	}, err
}

// Call requests a completion for the given prompt.
func (o *LLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	r, err := o.Generate(ctx, []string{prompt}, options...)
	if err != nil {
		return "", err
	}
	if len(r) == 0 {
		return "", ErrEmptyResponse
	}
	return r[0].Text, nil
}

func (o *LLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error) {
	if o.CallbacksHandler != nil {
		o.CallbacksHandler.HandleLLMStart(ctx, prompts)
	}

	opts := llms.CallOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	generations := make([]*llms.Generation, 0, len(prompts))
	for _, prompt := range prompts {
		result, err := o.client.CreateCompletion(ctx, &chatglm_client.CompletionRequest{

			Model:       opts.Model,
			Prompt:      prompt,
			Temperature: opts.Temperature,
			TopP:        opts.TopP,
			//RequestId:     opts.RequestId,
			StreamingFunc: opts.StreamingFunc,
		})
		if err != nil {
			return nil, err
		}
		generations = append(generations, &llms.Generation{
			Text: result.Text,
			GenerationInfo: map[string]any{"PromptTokens": result.Usage.PromptTokens,
				"CompletionTokens": result.Usage.CompletionTokens,
				"TotalTokens":      result.Usage.TotalTokens},
		})
	}

	if o.CallbacksHandler != nil {
		o.CallbacksHandler.HandleLLMEnd(ctx, llms.LLMResult{Generations: [][]*llms.Generation{generations}})
	}

	return generations, nil
}

func (o *LLM) GeneratePrompt(ctx context.Context, promptValues []schema.PromptValue, options ...llms.CallOption) (llms.LLMResult, error) { //nolint:lll
	return llms.GeneratePrompt(ctx, o, promptValues, options...)
}

func (o *LLM) GetNumTokens(text string) int {
	return 0
}

// CreateEmbedding creates embeddings for the given input texts.
func (o *LLM) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float64, error) {
	embeddings := make([][]float64, 0, 1)
	for _, input := range inputTexts {
		embedding, err := o.client.CreateEmbedding(ctx, &chatglm_client.EmbeddingRequest{
			Input: input,
		})
		if err != nil {
			return nil, err
		}
		if len(embedding) == 0 {
			return nil, ErrEmptyResponse
		}
		embeddings = append(embeddings, embedding)

	}
	if len(inputTexts) != len(embeddings) {
		return embeddings, ErrUnexpectedResponseLength
	}
	return embeddings, nil
}
