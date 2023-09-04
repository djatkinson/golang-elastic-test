package mapper

import (
	"splunk-test/entity"
	"splunk-test/model"
)

func ProductToProductModel(product *entity.Product) *model.Product {
	return &model.Product{
		ID:   product.ID,
		Name: product.Name,
	}
}
