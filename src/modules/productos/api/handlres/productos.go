package handlres

import (
	"api-test/src/common"
	"api-test/src/modules/productos/domain"
	"api-test/src/modules/productos/usecase"
)

type ProductosHandler struct {
	log common.Logger
	uc  usecase.ProductosUseCase
	*common.GenericHandler[domain.CreateProductosDTO, domain.ResponseProductosDTO, domain.UpdateProductosDTO, int64]
}


func NewProductosHandler(log common.Logger, uc usecase.ProductosUseCase) *ProductosHandler {
	handlers := common.NewGenericHandler(log, uc)
	return &ProductosHandler{
		log: log,
		uc: uc,
		GenericHandler: handlers,
	}
}