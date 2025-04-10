package domain

import (
	"time"

	"github.com/google/uuid"
)

type DTOTenant struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	UserID       uuid.UUID `json:"user_id"`
	DBName       string    `json:"db_name"`
	IsActive     bool      `json:"is_active"`
	CreationDate time.Time `json:"creation_date"`
}

func (dto *DTOTenant) FromTable(table TableTenant) {
	dto.ID = table.ID
	dto.Name = table.Name
	dto.DBName = table.DBName
	dto.IsActive = table.IsActive
	dto.CreationDate = table.CreationDate
}

func (dto *DTOTenant) ToTable() TableTenant {
	return TableTenant{
		ID:           dto.ID,
		Name:         dto.Name,
		DBName:       dto.DBName,
		IsActive:     dto.IsActive,
		CreationDate: dto.CreationDate,
	}
}

type DTOUserDirectory struct {
	ID           uuid.UUID       `json:"id"`
	Email        string          `json:"email"`
	Password     string          `json:"password"`
	Name         string          `json:"name"`
	LastLogin    time.Time       `json:"last_login"`
	IsActive     bool            `json:"is_active"`
	CreationDate time.Time       `json:"creation_date"`
	UserTenants  []DTOUserTenant `json:"user_tenants"`
}

func (dto *DTOUserDirectory) FromTable(table TableUserDirectory) {
	dto.ID = table.ID
	dto.Email = table.Email
	dto.IsActive = table.IsActive
	dto.CreationDate = table.CreationDate

	for _, userTenant := range table.UserTenants {
		dto.UserTenants = append(dto.UserTenants, DTOUserTenant{
			ID:            userTenant.ID,
			DefaultTenant: userTenant.DefaultTenant,
			CreationDate:  userTenant.CreationDate,
			Tenant: DTOTenant{
				ID:           userTenant.Tenant.ID,
				Name:         userTenant.Tenant.Name,
				DBName:       userTenant.Tenant.DBName,
				IsActive:     userTenant.Tenant.IsActive,
				CreationDate: userTenant.Tenant.CreationDate,
			},
		})
	}
}

func (dto *DTOUserDirectory) ToTable() TableUserDirectory {
	return TableUserDirectory{
		ID:           dto.ID,
		Email:        dto.Email,
		Password:     dto.Password,
		Name:         dto.Name,
		IsActive:     dto.IsActive,
		CreationDate: dto.CreationDate,
	}
}

type DTOUserTenant struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	TenantID      uuid.UUID `json:"tenant_id"`
	DefaultTenant bool      `json:"default_tenant"`
	CreationDate  time.Time `json:"creation_date"`
	Tenant        DTOTenant `json:"tenant"`
}

func (dto *DTOUserTenant) FromTable(table TableUserTenant) {
	dto.TenantID = table.TenantID
	dto.DefaultTenant = table.DefaultTenant
	dto.CreationDate = table.CreationDate

	dto.Tenant = DTOTenant{
		ID:           table.Tenant.ID,
		Name:         table.Tenant.Name,
		DBName:       table.Tenant.DBName,
		IsActive:     table.Tenant.IsActive,
		CreationDate: table.Tenant.CreationDate,
	}
}

func (dto *DTOUserTenant) ToTable() TableUserTenant {
	return TableUserTenant{
		ID:            dto.ID,
		UserID:        dto.UserID,
		TenantID:      dto.TenantID,
		DefaultTenant: dto.DefaultTenant,
		CreationDate:  dto.CreationDate,
	}
}

type DTOLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type DTORegister struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}

type DTOAuth struct {
	Token        string    `json:"token" validate:"required"`
	RefreshToken string    `json:"refresh_token" validate:"required"`
	ExpiresIn    time.Time `json:"expires_in"`
}
