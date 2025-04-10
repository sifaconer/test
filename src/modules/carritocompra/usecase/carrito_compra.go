package usecase

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/domain"
	"api-test/src/modules/carritocompra/repository"
	"context"
)

type CarritoCompra interface {
	Create(ctx context.Context, dto domain.DTOCarritoCompra) error
	Get(ctx context.Context, id int64) (domain.DTOCarritoCompra, error)
	Update(ctx context.Context, dto domain.DTOCarritoCompra) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filters common.QueryParams) ([]domain.DTOCarritoCompra, error)
}

type carritoCompra struct {
	log  common.Logger
	repo repository.CarritoCompra
}

// Create implements CarritoCompra.
func (c *carritoCompra) Create(ctx context.Context, dto domain.DTOCarritoCompra) error {
	return c.repo.Create(ctx, dto.ToTable())
}

// Delete implements CarritoCompra.
func (c *carritoCompra) Delete(ctx context.Context, id int64) error {
	return c.repo.Delete(ctx, id)
}

// Get implements CarritoCompra.
func (c *carritoCompra) Get(ctx context.Context, id int64) (domain.DTOCarritoCompra, error) {
	table, err := c.repo.Get(ctx, id)
	if err != nil {
		return domain.DTOCarritoCompra{}, err
	}
	dto := table.ToDTO()
	return dto, nil
}

// List implements CarritoCompra.
func (c *carritoCompra) List(ctx context.Context, filters common.QueryParams) ([]domain.DTOCarritoCompra, error) {
	tables, err := c.repo.List(ctx, filters)
	if err != nil {
		return nil, err
	}
	dtos := make([]domain.DTOCarritoCompra, len(tables))
	for i, table := range tables {
		dtos[i] = table.ToDTO()
	}
	return dtos, nil
}

// Update implements CarritoCompra.
func (c *carritoCompra) Update(ctx context.Context, dto domain.DTOCarritoCompra) error {
	return c.repo.Update(ctx, dto.ToTable())
}

func NewCarritoCompra(log common.Logger, repo repository.CarritoCompra) CarritoCompra {
	return &carritoCompra{
		log:  log,
		repo: repo,
	}
}

var _ CarritoCompra = (*carritoCompra)(nil)
