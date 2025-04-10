package domain

import "time"

type DTOCarritoCompra struct {
	ID            int64     `json:"id"`
	ClienteID     int64     `json:"cliente_id"`
	ProductoID    int64     `json:"producto_id"`
	Cantidad      int64     `json:"cantidad"`
	FechaAgregado time.Time `json:"fecha_agregado"`
}

func (dto *DTOCarritoCompra) FromTable(table TableCarritoCompra) {
	dto.ID = table.ID
	dto.ClienteID = table.ClienteID
	dto.ProductoID = table.ProductoID
	dto.Cantidad = table.Cantidad
	dto.FechaAgregado = table.FechaAgregado
}

func (dto *DTOCarritoCompra) ToTable() TableCarritoCompra {
	return TableCarritoCompra{
		ID:            dto.ID,
		ClienteID:     dto.ClienteID,
		ProductoID:    dto.ProductoID,
		Cantidad:      dto.Cantidad,
		FechaAgregado: dto.FechaAgregado,
	}
}
