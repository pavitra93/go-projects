package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	singletons "github.com/pavitra93/11-openai-chats/external"
	"github.com/pavitra93/11-openai-chats/internal/service"
)

func main() {
	_ = godotenv.Load()
	// Initialize openai client
	openapiKey := os.Getenv("OPENAI_API_KEY")
	// Initialize tokenizer client
	maxTokens, _ := strconv.ParseInt(os.Getenv("MAX_TOKENS"), 10, 64)
	temperature, _ := strconv.ParseFloat(os.Getenv("TEMPERATURE"), 64)
	systemMessage := os.Getenv("SYSTEM_MESSAGE")
	if openapiKey == "" {
		panic("OPENAI_API_KEY not found")
	}
	openAIServiceClient := singletons.GetOpenAIClientInstance(openapiKey)

	fmt.Println("========Chatbot with Memory=========")

	// Initialize No Memory chatbot service
	chatbotServiceBasic := &service.NoMemoryChatbotService{
		OpenAPIClient: openAIServiceClient.OpenAIClient,
		SystemMessage: systemMessage,
		MaxTokens:     maxTokens,
		Temperature:   temperature,
	}

	chatbotServiceBasic.RunNoMemoryChatbot()

	fmt.Println("========Chatbot with No Memory=========")

	// Initialize Memory chatbot service
	//chatbotServicePremium := &service.MemoryChatbotService{
	//	OpenAPIClient: openAIServiceClient.OpenAIClient,
	//	MaxTokens:     maxTokens,
	//	Temperature:   temperature,
	//	SystemMessage: systemMessage,
	//}
	//
	//worker := service.WorkerService{chatbotServicePremium}
	//chatbotServicePremium.RunMemoryChatbot(worker)

}
