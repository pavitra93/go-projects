package openai_memory_chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/openai/openai-go/v2"
)

type MemoryChatbot struct {
	OpenAPIClient *openai.Client
	MaxTokens     int64
	Temperature   float64
	History       *openai.ChatCompletionNewParams
	HistorySize   int
	SystemMessage string
}

func (service *MemoryChatbot) RunMemoryChatbot() {

	// start chatbot
	fmt.Println("Hello with Memory Chatbot")

	// send and recieve messages channel
	JobMessages := make(chan []openai.ChatCompletionMessageParamUnion)
	RecieveMessages := make(chan string)

	// create done channel
	doneChan := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create wait group
	wg := &sync.WaitGroup{}

	wg.Add(2)
	// start goroutine to send & recieve messages from OpenAI
	go service.SendMessagestoOpenAI(ctx, JobMessages, RecieveMessages, wg)
	go service.RecieveMessagesfromOpenAI(ctx, RecieveMessages, doneChan, wg)

	// initialize reader
	reader := bufio.NewReader(os.Stdin)

	// initialize history
	service.History = &openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(service.SystemMessage),
		},
	}

	// start chat loop
	for {
		dispatched := false
		fmt.Print("🧔🏻‍♂️ You: ")
		userMessage, _ := reader.ReadString('\n')
		userMessage = strings.TrimSpace(userMessage)
		// handle exit, quit and bye
		switch userMessage {
		case "", " ":
			fmt.Println("Please type your message")
			continue
		case "exit", "quit", "bye":
			fmt.Println("Bye. Thanks for chatting with me.")
			// cancel context
			cancel()

			// stop goroutines
			close(JobMessages)

			// wait for goroutines to finish
			wg.Wait()

			// close channels
			close(RecieveMessages)

			return
		default:
			// make history window and append user message
			service.History.Messages = MakeHistoryWindow(service.History.Messages, userMessage, service.HistorySize)
			service.History.Messages = append(service.History.Messages, openai.UserMessage(userMessage))
			JobMessages <- service.History.Messages
			dispatched = true
			fmt.Println("Bot is thinking...💭")

		}

		if dispatched {
			select {
			case <-doneChan:
				continue
			case <-ctx.Done():
				return
			}
		}

	}

}

func MakeHistoryWindow[T any](archive []T, userMessage string, keepLast int) []T {
	// how many from the tail (excluding the very first element)
	maxTail := len(archive) - 1
	if keepLast > maxTail {
		keepLast = maxTail
	}
	tail := archive[len(archive)-keepLast:]
	out := make([]T, 0, 1+len(tail))
	out = append(out, archive[0])
	out = append(out, tail...)
	return out
}

func (service *MemoryChatbot) SendMessagestoOpenAI(ctx context.Context, messages <-chan []openai.ChatCompletionMessageParamUnion, reciever chan<- string, wg *sync.WaitGroup) {
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
				MaxTokens:   openai.Int(service.MaxTokens),
				Temperature: openai.Float(service.Temperature),
			}

			// Send the request
			resp, err := service.OpenAPIClient.Chat.Completions.New(context.TODO(), param)
			if err != nil {
				reciever <- "Error: " + err.Error()
				return
			}

			// Safely print the first text part if the SDK returns structured content
			if len(resp.Choices) > 0 && len(resp.Choices[0].Message.Content) > 0 {
				service.History.Messages = append(service.History.Messages, resp.Choices[0].Message.ToParam())
			}

			// send messages back to channel
			reciever <- resp.Choices[0].Message.Content
		}

	}

}

func (service *MemoryChatbot) RecieveMessagesfromOpenAI(ctx context.Context, messages <-chan string, done chan<- bool, wg *sync.WaitGroup) {
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
