package minimax

import (
	"errors"
	"fmt"
	"github.com/tmc/langchaingo/llms/minimax/internal/minimaxclient"
	"net/http"
	"os"
)

var (
	ErrEmptyResponse            = errors.New("no response")
	ErrMissingToken             = errors.New("缺少GROUP ID 或者 API KEY") //nolint:lll
	ErrUnexpectedResponseLength = errors.New("unexpected length of response")
)

func newClient(opts ...Option) (*minimaxclient.Client, error) {
	options := &options{
		groupId:        os.Getenv(groupIdEnvVarName),
		apiKey:         os.Getenv(apiKeyEnvVarName),
		baseUrl:        os.Getenv(baseURLEnvVarName),
		httpClient:     http.DefaultClient,
		embeddingModel: "",
		model:          "",
	}

	fmt.Println(options.groupId)

	for _, opt := range opts {
		opt(options)
	}

	if options.model == "" {
		options.model = defaultModel
	}

	if options.embeddingModel == "" {
		options.embeddingModel = defaultEmbeddingModel
	}

	return minimaxclient.NewClient(minimaxclient.WithGroupId(options.groupId),
		minimaxclient.WithApiKey(options.apiKey),
		minimaxclient.WithBaseUrl(options.baseUrl),
		minimaxclient.WithHttpClient(options.httpClient),
		minimaxclient.WithModel(options.model),
		minimaxclient.WithEmbeddingsModel(options.embeddingModel),
	)
}
