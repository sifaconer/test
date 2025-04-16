package domain

import "time"

type ResponseCarritoCompraDTO struct {
	ID            int64     `json:"id"`
	ClienteID     int64     `json:"cliente_id"`
	ProductoID    int64     `json:"producto_id"`
	Cantidad      int64     `json:"cantidad"`
	FechaAgregado time.Time `json:"fecha_agregado"`
}

func (dto *ResponseCarritoCompraDTO) FromTable(table TableCarritoCompra) {
	dto.ID = table.ID
	dto.ClienteID = table.ClienteID
	dto.ProductoID = table.ProductoID
	dto.Cantidad = table.Cantidad
	dto.FechaAgregado = table.FechaAgregado
}

func (dto *ResponseCarritoCompraDTO) ToTable() TableCarritoCompra {
	return TableCarritoCompra{
		ID:            dto.ID,
		ClienteID:     dto.ClienteID,
		ProductoID:    dto.ProductoID,
		Cantidad:      dto.Cantidad,
		FechaAgregado: dto.FechaAgregado,
	}
}

type CreateCarritoCompraDTO struct {
	ClienteID int64 `json:"cliente_id" validate:"required"`
	ProductoID int64 `json:"producto_id" validate:"required"`
	Cantidad int64 `json:"cantidad" validate:"required"`
}

func (dto *CreateCarritoCompraDTO) ToTable() TableCarritoCompra {
	return TableCarritoCompra{
		ClienteID: dto.ClienteID,
		ProductoID: dto.ProductoID,
		Cantidad: dto.Cantidad,
	}
}


type UpdateCarritoCompraDTO struct {
	ID int64
	ClienteID int64 `json:"cliente_id" validate:"required"`
	ProductoID int64 `json:"producto_id" validate:"required"`
	Cantidad int64 `json:"cantidad" validate:"required"`
}

func (dto *UpdateCarritoCompraDTO) ToTable() TableCarritoCompra {
	return TableCarritoCompra{
		ID: dto.ID,
		ClienteID: dto.ClienteID,
		ProductoID: dto.ProductoID,
		Cantidad: dto.Cantidad,
	}
}
