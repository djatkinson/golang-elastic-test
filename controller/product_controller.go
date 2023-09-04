package controller

import (
	"github.com/gofiber/fiber/v2"
)

type ProductController interface {
	List(c *fiber.Ctx) error
	TestError(c *fiber.Ctx) error
}
