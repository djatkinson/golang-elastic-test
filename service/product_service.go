package service

import (
	"github.com/gofiber/fiber/v2"
	"splunk-test/model"
)

type ProductService interface {
	List(ctx *fiber.Ctx) ([]*model.Product, error)
}
