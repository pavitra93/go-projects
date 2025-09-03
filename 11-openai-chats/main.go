package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai_memory_chat "github.com/pavitra93/11-openai-chats/openai-memory-chat"
	//openai_no_memory_chat "github.com/pavitra93/11-openai-chats/openai-no-memory-chat"
	singletons "github.com/pavitra93/11-openai-chats/singletons"
)

func main() {
	_ = godotenv.Load()
	// Initialize openai client
	openapiKey := os.Getenv("OPENAI_API_KEY")
	if openapiKey == "" {
		panic("OPENAI_API_KEY not found")
	}
	openAIServiceClient := singletons.GetOpenAIClientInstance(openapiKey)

	//fmt.Println("========Chatbot with No Memory=========")
	//
	//noMemoryChatService := &openai_no_memory_chat.NoMemoryChatbot{
	//	OpenAPIClient: openAIClient.OpenaiClient,
	//	MaxTokens:     100,
	//	Temperature:   0.7,
	//}
	//noMemoryChatService.RunNoMemoryChatbot()

	fmt.Println("========Chatbot with Memory=========")

	MemoryChatService := &openai_memory_chat.MemoryChatbot{
		OpenAPIClient: openAIServiceClient.OpenAIClient,
		MaxTokens:     100,
		Temperature:   0.7,
		HistorySize:   5,
		SystemMessage: "You are good personal assistant. Never response in more tha 100 words",
	}

	MemoryChatService.RunMemoryChatbot()

}
