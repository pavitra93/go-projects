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
	doneChan := make(chan struct{})

	// start goroutine to send & recieve messages from OpenAI
	go service.SendMessagestoOpenAI(JobMessages, RecieveMessages)
	go service.RecieveMessagesfromOpenAI(RecieveMessages, doneChan)

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
		fmt.Print("🧔🏻‍♂️ You: ")
		userMessage, _ := reader.ReadString('\n')
		userMessage = strings.TrimSpace(userMessage)

		if userMessage == "" {
			fmt.Println("Please type your message")
			continue
		}

		// handle exit, quit and bye
		switch userMessage {
		case "exit":
		case "quit":
		case "bye":
			fmt.Println("Bye. Thanks for chatting with me.")
			doneChan <- struct{}{}
			break
		default:
			// make history window and append user message
			service.History.Messages = MakeHistoryWindow(service.History.Messages, userMessage, service.HistorySize)
			service.History.Messages = append(service.History.Messages, openai.UserMessage(userMessage))
			fmt.Println("Bot is thinking...💭")
			JobMessages <- service.History.Messages
		}

		<-doneChan
		continue

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

func (service *MemoryChatbot) SendMessagestoOpenAI(messages <-chan []openai.ChatCompletionMessageParamUnion, reciever chan<- string) {
	for {
		// get messages from channel
		InputMessages := <-messages

		// send messages to OpenAI
		param := openai.ChatCompletionNewParams{
			Messages:    InputMessages,
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

func (service *MemoryChatbot) RecieveMessagesfromOpenAI(messages <-chan string, done chan<- struct{}) {
	for {
		select {
		case msg := <-messages:
			fmt.Printf("🤖 Chatbot: %s\n", msg)
			done <- struct{}{}
		}
	}
}
