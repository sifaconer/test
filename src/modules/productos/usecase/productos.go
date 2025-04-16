package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/productos/domain"
)


type ProductosUseCase interface {
	common.UseCase[domain.ProductosDTO, int64]
}

type productosUseCase struct {
	config *config.Config
	log common.Logger
	tenant *common.TenantConnectionManager
	repo common.Repository[domain.ProductosTable, int64]
	common.UseCase[domain.ProductosDTO, int64]
}


func NewProductosUseCase(config *config.Config, log common.Logger, tenant *common.TenantConnectionManager, repo common.Repository[domain.ProductosTable, int64]) ProductosUseCase {

	toTable := func(dto domain.ProductosDTO) domain.ProductosTable {
		return dto.ToTable()
	}

	toDTO := func(table domain.ProductosTable) domain.ProductosDTO {
		return table.ToDTO()
	}

	genericUC := common.NewUseCase(config, log, tenant, repo, toTable, toDTO)

	return &productosUseCase{config, log, tenant, repo, genericUC}
}

var _ ProductosUseCase = (*productosUseCase)(nil)