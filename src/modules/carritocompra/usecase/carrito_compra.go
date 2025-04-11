package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/carritocompra/domain"
)

type CarritoCompra interface {
	common.UseCase[domain.DTOCarritoCompra, int64]
}

type carritoCompra struct {
	log  common.Logger
	repo common.Repository[domain.TableCarritoCompra, int64]
	common.UseCase[domain.DTOCarritoCompra, int64]
}

func NewCarritoCompra(config *config.Config, log common.Logger, tenant *common.TenantConnectionManager, repo common.Repository[domain.TableCarritoCompra, int64]) CarritoCompra {

	toTable := func(dto domain.DTOCarritoCompra) domain.TableCarritoCompra {
		return dto.ToTable()
	}

	toDTO := func(table domain.TableCarritoCompra) domain.DTOCarritoCompra {
		return table.ToDTO()
	}

	genericUC := common.NewUseCase(config, log, tenant, repo, toTable, toDTO)

	return &carritoCompra{
		log:  log,
		repo: repo,
		UseCase: genericUC,
	}
}

var _ CarritoCompra = (*carritoCompra)(nil)
