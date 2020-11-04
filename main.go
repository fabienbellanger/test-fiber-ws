package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Static("/", ".")

	hub := newHub()
	go hub.run()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("home.html", fiber.Map{})
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		serveWs(hub, c)
	}))

	log.Fatal(app.Listen(":3000"))
}
