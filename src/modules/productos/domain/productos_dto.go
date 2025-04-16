package domain

type ResponseProductosDTO struct {
	Id int64 `json:"id" params:"id"`
	Nombre string `json:"nombre"`
	Precio float64 `json:"precio"`
}

func (dto *ResponseProductosDTO) FromTable(table ProductosTable) {
	dto.Id = table.ID
	dto.Nombre = table.Nombre
	dto.Precio = table.Precio
}

func (dto *ResponseProductosDTO) ToTable() ProductosTable {
	return ProductosTable{
		ID: dto.Id,
		Nombre: dto.Nombre,
		Precio: dto.Precio,
	}
}

type UpdateProductosDTO struct {
	Id int64 `json:"id" params:"id" validate:"required"`
	Nombre string `json:"nombre"`
	Precio float64 `json:"precio"`
}

func (dto *UpdateProductosDTO) ToTable() ProductosTable {
	return ProductosTable{
		ID: dto.Id,
		Nombre: dto.Nombre,
		Precio: dto.Precio,
	}
}

type CreateProductosDTO struct {
	Nombre string `json:"nombre" validate:"required"`
	Precio float64 `json:"precio" validate:"required"`
}

func (dto *CreateProductosDTO) ToTable() ProductosTable {
	return ProductosTable{
		Nombre: dto.Nombre,
		Precio: dto.Precio,
	}
}
