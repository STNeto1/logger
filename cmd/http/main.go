package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/stneto1/logger/pkg"
)

func main() {
	pkg.InitDB()

	log.Println("Http Server")

	msgChan := make(chan pkg.Message, 100_000)

	go handleMessageChannel(msgChan)

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})
	app.Use(logger.New())

	app.Post("/log", func(c *fiber.Ctx) error {
		var msg pkg.Message

		if err := c.BodyParser(&msg); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		msgChan <- msg

		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"message": "success",
		})
	})

	app.Get("/log", func(c *fiber.Ctx) error {
		msgs, err := pkg.DBCon.GetMessages()

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(msgs)
	})

	app.Get("/log/:topic", func(c *fiber.Ctx) error {
		topic := c.Params("topic")

		msgs, err := pkg.DBCon.GetMessagesByTopic(topic)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(http.StatusOK).JSON(msgs)
	})

	app.Listen(":3000")
}

func handleMessageChannel(ch chan pkg.Message) {
	// queue := make([]pkg.Message, 0, 100)
	var queue []pkg.Message

	for {
		select {
		case msg := <-ch:
			if len(queue) == 100 {
				fmt.Println("batch size reached, sending to database")
				if err := pkg.DBCon.CreateMessages(queue); err != nil {
					fmt.Println("error on batch insert:", err)
				} else {
					fmt.Println("batch insert success")
				}

				queue = nil
			}

			queue = append(queue, msg)

		case <-time.After(1 * time.Second):
			fmt.Printf("pool size: %d\n", len(queue))
		}
	}
}
