package domain

type ProductosDTO struct {
	Id int64 `json:"id" params:"id"`
	Nombre string `json:"nombre"`
	Precio float64 `json:"precio"`
}

func (dto *ProductosDTO) FromTable(table ProductosTable) {
	dto.Id = table.ID
	dto.Nombre = table.Nombre
	dto.Precio = table.Precio
}

func (dto *ProductosDTO) ToTable() ProductosTable {
	return ProductosTable{
		ID: dto.Id,
		Nombre: dto.Nombre,
		Precio: dto.Precio,
	}
}