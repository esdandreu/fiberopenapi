package router

// Code generated by "genrouter -spec ./specification.json"; DO NOT EDIT.

import (
	"github.com/gofiber/fiber/v2"
)

type Handlers interface {
	GetBoard(c *fiber.Ctx) error
	GetSquare(c *fiber.Ctx) error
	PutSquare(c *fiber.Ctx) error
}

func AddHandlers(app *fiber.App, h Handlers) {
	app.Get("/board", h.GetBoard)
	app.Get("/board/:row/:column", h.GetSquare)
	app.Put("/board/:row/:column", h.PutSquare)
}
