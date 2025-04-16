package handlers

import (
	"api-test/src/common"
	"api-test/src/modules/carritocompra/domain"
	"api-test/src/modules/carritocompra/usecase"
)

type CarritoCompraHandler struct {
	log common.Logger
	uc  usecase.CarritoCompra
	*common.GenericHandler[domain.DTOCarritoCompra, int64]
}


func NewCarritoCompraHandler(log common.Logger, uc usecase.CarritoCompra) *CarritoCompraHandler {
	handlers := common.NewGenericHandler(log, uc)
	return &CarritoCompraHandler{
		log: log,
		uc:  uc,
		GenericHandler: handlers,
	}
}
