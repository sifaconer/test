package repository

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/domain"
	"context"
)

// repository interface to crud and query operations
type CarritoCompra interface {
	// Create
	Create(ctx context.Context, table domain.TableCarritoCompra) error
	// Read
	Get(ctx context.Context, id int64) (domain.TableCarritoCompra, error)
	// Update
	Update(ctx context.Context, table domain.TableCarritoCompra) error
	// Delete
	Delete(ctx context.Context, id int64) error
	// List
	List(ctx context.Context, filters common.QueryParams) ([]domain.TableCarritoCompra, error)
}
