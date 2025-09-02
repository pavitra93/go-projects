package singleton_openai_client

import (
	"sync"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type openAIServiceClient struct {
	OpenaiClient *openai.Client
}

var openapiInstance *openAIServiceClient
var once sync.Once

func GetInstance(openapiKey string) *openAIServiceClient {
	once.Do(func() {
		client := openai.NewClient(
			option.WithAPIKey(openapiKey), // defaults to os.LookupEnv("OPENAI_API_KEY")
		)
		openapiInstance = &openAIServiceClient{
			OpenaiClient: &client,
		}
	})

	return openapiInstance

}
