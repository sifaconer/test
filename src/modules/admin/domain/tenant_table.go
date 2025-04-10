package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TableTenant struct {
	bun.BaseModel `bun:"table:tenants.tenants"`

	ID                uuid.UUID `bun:"id,pk"`
	Name              string    `bun:"name,notnull"`
	DBName            string    `bun:"db_name,notnull"`
	DBHost            string    `bun:"db_host,notnull"`
	DBPort            int64     `bun:"db_port,notnull"`
	DBUser            string    `bun:"db_user,notnull"`
	DBPassword        []byte    `bun:"db_password,notnull"`
	IsActive          bool      `bun:"is_active,notnull,default:true"`
	IV                []byte    `bun:"iv,notnull"`
	Version           string    `bun:"version,notnull"`
	CreationDate      time.Time `bun:"created_at,notnull,default:current_timestamp"`
	PasswordPlaintext string    `bun:"-"`
}

func (table *TableTenant) FromDTO(dto DTOTenant) {
	table.ID = dto.ID
	table.Name = dto.Name
	table.DBName = dto.DBName
	table.IsActive = dto.IsActive
	table.CreationDate = dto.CreationDate
}

func (table *TableTenant) ToDTO() DTOTenant {
	return DTOTenant{
		ID:           table.ID,
		Name:         table.Name,
		DBName:       table.DBName,
		IsActive:     table.IsActive,
		CreationDate: table.CreationDate,
	}
}

type TableUserDirectory struct {
	bun.BaseModel `bun:"table:tenants.users_directory"`

	ID           uuid.UUID `bun:"id,pk,default:uuid_generate_v4()"`
	Email        string    `bun:"email,notnull"`
	Name         string    `bun:"name,notnull"`
	Password     string    `bun:"password,notnull"`
	IsActive     bool      `bun:"is_active,notnull,default:true"`
	CreationDate time.Time `bun:"created_at,notnull,default:current_timestamp"`

	UserTenants []TableUserTenant `bun:"rel:has-many,join:id=user_id"`
}

func (table *TableUserDirectory) FromDTO(dto DTOUserDirectory) {
	table.ID = dto.ID
	table.Email = dto.Email
	table.IsActive = dto.IsActive
	table.CreationDate = dto.CreationDate

	for _, userTenant := range dto.UserTenants {
		table.UserTenants = append(table.UserTenants, TableUserTenant{
			UserID:        userTenant.TenantID,
			TenantID:      userTenant.TenantID,
			DefaultTenant: userTenant.DefaultTenant,
			CreationDate:  userTenant.CreationDate,
		})
	}
}

func (table *TableUserDirectory) ToDTO() DTOUserDirectory {
	var userTenants []DTOUserTenant
	for _, userTenant := range table.UserTenants {
		userTenants = append(userTenants, userTenant.ToDTO())
	}
	return DTOUserDirectory{
		ID:           table.ID,
		Email:        table.Email,
		Name:         table.Name,
		IsActive:     table.IsActive,
		UserTenants:  userTenants,
	}
}

type TableUserTenant struct {
	bun.BaseModel `bun:"table:tenants.user_tenants"`

	ID            uuid.UUID          `bun:"id,pk"`
	UserID        uuid.UUID          `bun:"user_id,notnull"`
	TenantID      uuid.UUID          `bun:"tenant_id,notnull"`
	DefaultTenant bool               `bun:"default_tenant,notnull"`
	CreationDate  time.Time          `bun:"created_at,notnull"`
	UserDirectory TableUserDirectory `bun:"rel:belongs-to,join:user_id=id"`
	Tenant        TableTenant        `bun:"rel:belongs-to,join:tenant_id=id"`
}

func (table *TableUserTenant) FromDTO(dto DTOUserTenant) {
	table.ID = dto.ID
	table.UserID = dto.UserID
	table.TenantID = dto.TenantID
	table.DefaultTenant = dto.DefaultTenant
	table.CreationDate = dto.CreationDate
}

func (table *TableUserTenant) ToDTO() DTOUserTenant {
	return DTOUserTenant{
		TenantID:      table.TenantID,
		DefaultTenant: table.DefaultTenant,
		CreationDate:  table.CreationDate,
	}
}
