package openai_memory_chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go/v2"
)

type MemoryChatbot struct {
	OpenAPIClient *openai.Client
	MaxTokens     int64
	Temperature   float64
}

func (service *MemoryChatbot) RunMemoryChatbot() {

	fmt.Println("Hello from No Memory Chatbot")
	fmt.Println("Chat with OpenAI GPT-4 is ready to talk")
	systemMessage := "You are good personal assistant. Never response in more tha 100 words"
	reader := bufio.NewReader(os.Stdin)

	history := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemMessage),
		},
	}

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

		history.Messages = append(history.Messages, openai.UserMessage(userMessage))
		param := openai.ChatCompletionNewParams{
			Messages:    history.Messages,
			Model:       openai.ChatModelGPT4_1,
			MaxTokens:   openai.Int(service.MaxTokens),
			Temperature: openai.Float(service.Temperature),
		}

		resp, err := service.OpenAPIClient.Chat.Completions.New(context.TODO(), param)

		if err != nil {
			panic(err)
		}

		// Safely print the first text part if the SDK returns structured content
		if len(resp.Choices) > 0 && len(resp.Choices[0].Message.Content) > 0 {
			fmt.Printf("Chatbot: %s\n", resp.Choices[0].Message.Content)
			history.Messages = append(history.Messages, resp.Choices[0].Message.ToParam())
			continue
		}

	}

}
