package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type TableCarritoCompra struct {
	bun.BaseModel `bun:"table:carrito_compra"`

	ID            int64     `bun:"id"`
	ClienteID     int64     `bun:"cliente_id"`
	ProductoID    int64     `bun:"producto_id"`
	Cantidad      int64     `bun:"cantidad"`
	FechaAgregado time.Time `bun:"fecha_agregado"`
}

func (table *TableCarritoCompra) FromDTO(dto DTOCarritoCompra) {
	table.ID = dto.ID
	table.ClienteID = dto.ClienteID
	table.ProductoID = dto.ProductoID
	table.Cantidad = dto.Cantidad
	table.FechaAgregado = dto.FechaAgregado
}

func (table *TableCarritoCompra) ToDTO() DTOCarritoCompra {
	return DTOCarritoCompra{
		ID:            table.ID,
		ClienteID:     table.ClienteID,
		ProductoID:    table.ProductoID,
		Cantidad:      table.Cantidad,
		FechaAgregado: table.FechaAgregado,
	}
}
