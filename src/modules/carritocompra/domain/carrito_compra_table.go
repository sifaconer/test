package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type TableCarritoCompra struct {
	bun.BaseModel `bun:"table:carrito_compra"`

	ID            int64     `bun:"id,pk,autoincrement"`
	ClienteID     int64     `bun:"cliente_id"`
	ProductoID    int64     `bun:"producto_id"`
	Cantidad      int64     `bun:"cantidad"`
	FechaAgregado time.Time `bun:"fecha_agregado"`
}


func (table *TableCarritoCompra) ToDTO() ResponseCarritoCompraDTO {
	return ResponseCarritoCompraDTO{
		ID:            table.ID,
		ClienteID:     table.ClienteID,
		ProductoID:    table.ProductoID,
		Cantidad:      table.Cantidad,
		FechaAgregado: table.FechaAgregado,
	}
}
