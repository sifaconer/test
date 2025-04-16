package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/carritocompra/domain"
	"context"

	"github.com/uptrace/bun"
)

type CarritoCompra interface {
	common.UseCase[domain.CreateCarritoCompraDTO, domain.ResponseCarritoCompraDTO, domain.UpdateCarritoCompraDTO, int64]
}

type carritoCompra struct {
	log  common.Logger
	repo common.Repository[domain.TableCarritoCompra, int64]
	common.UseCase[domain.CreateCarritoCompraDTO, domain.ResponseCarritoCompraDTO, domain.UpdateCarritoCompraDTO, int64]
}

func (c *carritoCompra) Create(ctx context.Context, dto domain.CreateCarritoCompraDTO) (*domain.ResponseCarritoCompraDTO, error) {
	var result domain.ResponseCarritoCompraDTO
	err := c.repo.WithTransaction(ctx, func(ctx context.Context, tx bun.Tx) error {
		table, err := c.repo.CreateTx(ctx, tx, dto.ToTable())
		if err != nil {
			return err
		}
		result = table.ToDTO()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func NewCarritoCompra(config *config.Config, log common.Logger, tenant *common.TenantConnectionManager, repo common.Repository[domain.TableCarritoCompra, int64]) CarritoCompra {

	createToTable := func(dto domain.CreateCarritoCompraDTO) domain.TableCarritoCompra {
		return dto.ToTable()
	}

	updateToTable := func(dto domain.UpdateCarritoCompraDTO) domain.TableCarritoCompra {
		return dto.ToTable()
	}

	toDTO := func(table domain.TableCarritoCompra) domain.ResponseCarritoCompraDTO {
		return table.ToDTO()
	}

	genericUC := common.NewUseCase(config, log, tenant, repo, createToTable, updateToTable, toDTO)

	return &carritoCompra{
		log:  log,
		repo: repo,
		UseCase: genericUC,
	}
}

var _ CarritoCompra = (*carritoCompra)(nil)
