package usecase

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/modules/productos/domain"
)


type ProductosUseCase interface {
	common.UseCase[domain.CreateProductosDTO, domain.ResponseProductosDTO, domain.UpdateProductosDTO, int64]
}

type productosUseCase struct {
	config *config.Config
	log common.Logger
	tenant *common.TenantConnectionManager
	repo common.Repository[domain.ProductosTable, int64]
	common.UseCase[domain.CreateProductosDTO, domain.ResponseProductosDTO, domain.UpdateProductosDTO, int64]
}


func NewProductosUseCase(config *config.Config, log common.Logger, tenant *common.TenantConnectionManager, repo common.Repository[domain.ProductosTable, int64]) ProductosUseCase {

	createToTable := func(dto domain.CreateProductosDTO) domain.ProductosTable {
		return dto.ToTable()
	}

	toDTO := func(table domain.ProductosTable) domain.ResponseProductosDTO {
		return table.ToDTO()
	}

	updateToTable := func(dto domain.UpdateProductosDTO) domain.ProductosTable {
		return dto.ToTable()
	}

	genericUC := common.NewUseCase(config, log, tenant, repo, createToTable, updateToTable, toDTO)

	return &productosUseCase{config, log, tenant, repo, genericUC}
}

var _ ProductosUseCase = (*productosUseCase)(nil)