package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	openai_memory_chat "github.com/pavitra93/11-openai-chats/openai-memory-chat"
	//openai_no_memory_chat "github.com/pavitra93/11-openai-chats/openai-no-memory-chat"
	singleton_openai_client "github.com/pavitra93/11-openai-chats/singleton-openai-client"
)

func main() {
	_ = godotenv.Load()
	// Initialize openai client
	openapiKey := os.Getenv("OPENAI_API_KEY")
	if openapiKey == "" {
		panic("OPENAI_API_KEY not found")
	}
	openAIClient := singleton_openai_client.GetInstance(openapiKey)

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
		OpenAPIClient: openAIClient.OpenaiClient,
		MaxTokens:     100,
		Temperature:   0.7,
	}
	MemoryChatService.RunMemoryChatbot()

}
