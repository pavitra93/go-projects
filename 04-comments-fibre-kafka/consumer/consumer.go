package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	topic := "comments"
	worker, err := connectConsumer([]string{"localhost:9092"})
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer Started")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	msgCount := 0
	doneChan := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				fmt.Printf("Received Message: %s | count %d \n", string(msg.Value), msgCount)
			case <-sigChan:
				fmt.Println("Received SIGINT, shutting down")
				doneChan <- struct{}{}
			}
		}
	}()

	<-doneChan
	fmt.Println("Consumer Stopped")
	fmt.Println("Message Count:", msgCount)
}

func connectConsumer(URL []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer(URL, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
