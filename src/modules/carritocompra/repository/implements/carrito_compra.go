package implements

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/domain"
	"api-test/src/modules/carritocompra/repository"
	"context"
)

type carritoCompraRepository struct {
	log    common.Logger
	tenant *common.TenantConnectionManager
}

// Create implements repository.CarritoCompra.
func (c *carritoCompraRepository) Create(ctx context.Context, table domain.TableCarritoCompra) error {
	db, err := c.tenant.GetDBContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().Model(&table).Exec(ctx)
	if err != nil {
		return common.CheckDBErrorType(err)
	}
	return nil
}

// Delete implements repository.CarritoCompra.
func (c *carritoCompraRepository) Delete(ctx context.Context, id int64) error {
	db, err := c.tenant.GetDBContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewDelete().Model(&domain.TableCarritoCompra{ID: id}).Exec(ctx)
	if err != nil {
		return common.CheckDBErrorType(err)
	}
	return nil
}

// Get implements repository.CarritoCompra.
func (c *carritoCompraRepository) Get(ctx context.Context, id int64) (domain.TableCarritoCompra, error) {
	db, err := c.tenant.GetDBContext(ctx)
	if err != nil {
		return domain.TableCarritoCompra{}, err
	}

	var table domain.TableCarritoCompra
	_, err = db.NewSelect().Model(&table).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return domain.TableCarritoCompra{}, common.CheckDBErrorType(err)
	}
	return table, nil
}

// List implements repository.CarritoCompra.
func (c *carritoCompraRepository) List(ctx context.Context, filters common.QueryParams) ([]domain.TableCarritoCompra, error) {
	db, err := c.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	tables := make([]domain.TableCarritoCompra, 0)
	_, err = db.NewSelect().Model(&tables).Exec(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return tables, nil
}

// Update implements repository.CarritoCompra.
func (c *carritoCompraRepository) Update(ctx context.Context, table domain.TableCarritoCompra) error {
	db, err := c.tenant.GetDBContext(ctx)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().Model(&table).Where("id = ?", table.ID).Exec(ctx)
	if err != nil {
		return common.CheckDBErrorType(err)
	}
	return nil
}

func NewCarritoCompraRepository(log common.Logger, tenant *common.TenantConnectionManager) repository.CarritoCompra {
	return &carritoCompraRepository{
		log:    log,
		tenant: tenant,
	}
}

var _ repository.CarritoCompra = (*carritoCompraRepository)(nil)
