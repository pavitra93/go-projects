package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Comments struct {
	Text string `json:"text"`
}

func main() {
	app := fiber.New()
	api := app.Group("/api")
	api.Post("/comments", CreateComment)
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateComment(c *fiber.Ctx) error {
	NewComment := new(Comments)
	if err := c.BodyParser(NewComment); err != nil {
		log.Println(err)
		c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": err.Error(),
		})

		return err
	}

	commentInBytes, err := json.Marshal(NewComment)
	if err != nil {
		log.Println(err)
	}

	PushCommentToQueue("comments", commentInBytes)

	err = c.JSON(fiber.Map{
		"success":  true,
		"message":  "Comment pushed successfully",
		"comments": string(commentInBytes),
	})

	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"message": "Error creating comment",
		})

		return err
	}

	return err
}

func PushCommentToQueue(topic string, commentInBytes []byte) error {
	brokerURL := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokerURL)
	if err != nil {
		return err
	}

	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(commentInBytes),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message successfully pushed to topic %s/partition %d at offset %d\n ", topic, partition, offset)
	return nil
}

func ConnectProducer(URL []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(URL, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
