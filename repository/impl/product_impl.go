package impl

import (
	"context"
	"github.com/jmoiron/sqlx"
	"splunk-test/entity"
	"splunk-test/logger"
	"splunk-test/repository"
)

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) repository.ProductRepository {
	return &productRepository{db}
}

func (p *productRepository) List(ctx context.Context) ([]*entity.Product, error) {
	var result []*entity.Product
	query := `SELECT * FROM product`
	l := logger.Log(ctx)

	rows, err := p.db.QueryxContext(ctx, query)
	if err != nil {
		l.Error().Err(err)
		return nil, err
	}
	defer rows.Close()
	if rows == nil {
		l.Error().Err(err)
		return nil, nil
	}

	for rows.Next() {
		product := new(entity.Product)
		err := rows.StructScan(&product)
		if err != nil {
			l.Error().Err(err)
			return nil, err
		}
		result = append(result, product)
	}

	return result, nil
}
