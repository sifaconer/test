package usecase

import (
	"api-test/src/common"
	"api-test/src/modules/productos/domain"
)


type ProductosUseCase interface {
	common.Repository[domain.ProductosTable, int64]
}

type productosUseCase struct {
	common.Repository[domain.ProductosTable, int64]
}


func NewProductosUseCase(repo common.Repository[domain.ProductosTable, int64]) ProductosUseCase {
	return &productosUseCase{repo}
}

var _ ProductosUseCase = (*productosUseCase)(nil)