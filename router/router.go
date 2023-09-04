package router

import (
	"github.com/gofiber/fiber/v2"
	"splunk-test/controller"
)

func RegisterRouter(app *fiber.App, productController controller.ProductController) {
	app.Get("/product", productController.List)
	app.Get("/test-error", productController.TestError)
}
