package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/openai/openai-go/v2"
)

type WorkerService struct {
	MemoryChatbotService *MemoryChatbotService
}

func (w *WorkerService) SendMessagestoOpenAI(ctx context.Context, messages <-chan []openai.ChatCompletionMessageParamUnion, reciever chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-messages:
			if !ok {
				return
			}

			// send messages to OpenAI
			param := openai.ChatCompletionNewParams{
				Messages:    message,
				Model:       openai.ChatModelGPT4_1,
				MaxTokens:   openai.Int(w.MemoryChatbotService.MaxTokens),
				Temperature: openai.Float(w.MemoryChatbotService.Temperature),
			}

			// Send the request
			resp, err := w.MemoryChatbotService.OpenAPIClient.Chat.Completions.New(context.TODO(), param)
			if err != nil {
				reciever <- "Error: " + err.Error()
				return
			}

			// Safely print the first text part if the SDK returns structured content
			if len(resp.Choices) > 0 && len(resp.Choices[0].Message.Content) > 0 {
				w.MemoryChatbotService.History.Messages = append(w.MemoryChatbotService.History.Messages, resp.Choices[0].Message.ToParam())
			}

			// send messages back to channel
			reciever <- resp.Choices[0].Message.Content
		}

	}

}

func (w *WorkerService) RecieveMessagesfromOpenAI(ctx context.Context, messages <-chan string, done chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-messages:
			if !ok {
				return
			}
			fmt.Printf("🤖 Chatbot: %s\n", msg)
			done <- true
		}
	}
}
