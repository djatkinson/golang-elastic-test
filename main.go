package main

import (
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/module/apmfiber/v2"
	"splunk-test/config"
	impl3 "splunk-test/controller/impl"
	"splunk-test/database"
	"splunk-test/middleware"
	"splunk-test/repository/impl"
	"splunk-test/router"
	impl2 "splunk-test/service/impl"
)

func main() {
	//log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	config.Load()
	//
	//sdk, err := distro.Run()
	//if err != nil {
	//	panic(err)
	//}
	//// Flush all spans before the application exits
	//defer func() {
	//	if err := sdk.Shutdown(context.Background()); err != nil {
	//		panic(err)
	//	}
	//}()
	//
	db, err := database.Init()
	if err != nil {
		panic(err)
	}

	productRepository := impl.NewProductRepository(db)

	productService := impl2.NewProductService(productRepository)

	productController := impl3.NewProductController(productService)

	//elkLogger()
	//
	app := fiber.New()
	//
	//app.Use(cors.New())
	app.Use(apmfiber.Middleware())
	app.Use(middleware.Config())
	router.RegisterRouter(app, productController)
	app.Listen(":4000")
}
