package repository

import (
	"context"
	"splunk-test/entity"
)

type ProductRepository interface {
	List(ctx context.Context) ([]*entity.Product, error)
}
