package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/openai/openai-go/v2"
	"github.com/pavitra93/11-openai-chats/pkg/utils"
)

type MemoryChatbotService struct {
	OpenAPIClient *openai.Client
	MaxTokens     int64
	Temperature   float64
	History       *openai.ChatCompletionNewParams
	SystemMessage string
	HistorySize   int
}

func (service *MemoryChatbotService) RunMemoryChatbot(worker WorkerService) {

	// start chatbot
	fmt.Println("Hello with Memory Chatbot")

	// send and recieve messages channel
	JobMessages := make(chan []openai.ChatCompletionMessageParamUnion)
	ReceiveMessages := make(chan string)

	// create done channel
	doneChan := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create wait group
	wg := &sync.WaitGroup{}

	wg.Add(2)

	// start goroutine to send & recieve messages from OpenAI
	go worker.SendMessagestoOpenAI(ctx, JobMessages, ReceiveMessages, wg)
	go worker.RecieveMessagesfromOpenAI(ctx, ReceiveMessages, doneChan, wg)

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
			close(ReceiveMessages)

			return
		default:
			// make history window and append user message
			service.History.Messages = utils.MakeHistoryWindow(service.History.Messages, userMessage, service.HistorySize)
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
