package domain

import "github.com/uptrace/bun"


type ProductosTable struct {
	bun.BaseModel `bun:"table:productos,alias:p"`
	ID            int64  `bun:"id,pk,autoincrement"`
	Nombre        string `bun:"nombre,notnull"`
	Precio        float64    `bun:"precio,notnull"`
}

func (p *ProductosTable) ToDTO() ResponseProductosDTO {
	return ResponseProductosDTO{
		Id:   p.ID,
		Nombre: p.Nombre,
		Precio: p.Precio,
	}
}