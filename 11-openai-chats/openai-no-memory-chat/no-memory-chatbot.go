package openai_no_memory_chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
)

type NoMemoryChatbot struct {
	OpenAPIClient *openai.Client
	MaxTokens     int64
	Temperature   float64
}

func (service *NoMemoryChatbot) RunNoMemoryChatbot() {

	fmt.Println("Hello with no Memory Chatbot")
	systemMessage := "You are good personal assistant. Never response in more tha 100 words"
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("You: ")
		userMessage, _ := reader.ReadString('\n')
		userMessage = strings.TrimSpace(userMessage)

		if userMessage == "" {
			fmt.Println("Please type your message")
			continue
		}

		if userMessage == "exit" || userMessage == "quit" || userMessage == "bye" {
			fmt.Println("Bye. Thanks for chatting with me.")
			break
		}

		resp, err := service.OpenAPIClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(userMessage),
				openai.SystemMessage(systemMessage),
			},
			Model:       openai.ChatModelGPT4_1,
			MaxTokens:   openai.Int(service.MaxTokens),
			Temperature: openai.Float(service.Temperature),
		})

		if err != nil {
			fmt.Println("Error: ", err)
			break
		}

		// Safely print the first text part if the SDK returns structured content
		if len(resp.Choices) > 0 && len(resp.Choices[0].Message.Content) > 0 {
			fmt.Printf("Chatbot: %s\n", resp.Choices[0].Message.Content)
			continue
		}

	}

}
