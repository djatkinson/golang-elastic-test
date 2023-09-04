package impl

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"splunk-test/logger"
	"splunk-test/mapper"
	"splunk-test/model"
	"splunk-test/repository"
	"splunk-test/service"
)

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) service.ProductService {
	return &productService{productRepository: productRepository}
}

func (p productService) List(ctx *fiber.Ctx) ([]*model.Product, error) {
	products, err := p.productRepository.List(ctx.Context())
	if err != nil {
		logger.LogError(ctx, errors.New("error-test"))
		return nil, err
	}

	var results []*model.Product
	for _, product := range products {
		results = append(results, mapper.ProductToProductModel(product))
	}

	logger.LogInfo(ctx, "test-logging")
	logger.LogError(ctx, errors.New("error-test"))

	return results, nil
}
