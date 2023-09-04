package impl

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"reflect"
	"splunk-test/controller"
	"splunk-test/service"
	"unsafe"
)

type productController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) controller.ProductController {
	return &productController{productService: productService}
}

func (p productController) List(c *fiber.Ctx) error {
	products, err := p.productService.List(c)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(`"error":Internal Server error`)
	}

	//var data map[string]interface{}
	//req, _ := http.NewRequest("GET", "http://localhost:8083/test-data", nil)
	//client := apmhttp.WrapClient(http.DefaultClient)
	//resp, err := client.Do(req.WithContext(c.Context()))
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//err = json.NewDecoder(resp.Body).Decode(&data)
	//if err != nil {
	//	panic(err)
	//}
	////fmt.Print(data)
	return c.Status(fiber.StatusOK).JSON(products)
}

func printContextInternals(ctx interface{}, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				printContextInternals(reflectValue.Interface(), true)
			} else {
				fmt.Printf("field name: %+v\n", reflectField.Name)
				fmt.Printf("field name: %+v\n", reflectField.Type)
				fmt.Printf("value: %+v\n", reflectValue.Interface())
			}
		}
	} else {
		fmt.Printf("context is empty (int)\n")
	}
}

func (p productController) TestError(c *fiber.Ctx) error {
	products, err := p.productService.List(c)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(`"error":Internal Server error`)
	}
	return c.Status(fiber.StatusOK).JSON(products)

}
