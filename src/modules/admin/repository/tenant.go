package repository

import (
	"api-test/src/modules/admin/domain"
	"context"

	"github.com/google/uuid"
)

type TenantRepository interface {
	GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.TableTenant, error)
	GetTenantByName(ctx context.Context, name string) (*domain.TableTenant, error)
	CreateTenant(ctx context.Context, tenant domain.TableTenant) (*domain.TableTenant, error)
	UpdateTenant(ctx context.Context, tenant domain.TableTenant) (*domain.TableTenant, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	GetAllTenants(ctx context.Context) ([]domain.TableTenant, error)
	CreateUserTenant(ctx context.Context, userTenant domain.TableUserTenant) (*domain.TableUserTenant, error)
	GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]domain.TableUserTenant, error)
}

type UserDirectoryRepository interface {
	GetUserDirectoryByID(ctx context.Context, id uuid.UUID) (*domain.TableUserDirectory, error)
	GetUserDirectoryByEmail(ctx context.Context, email string) (*domain.TableUserDirectory, error)
	CreateUserDirectory(ctx context.Context, userDirectory domain.TableUserDirectory) (*domain.TableUserDirectory, error)
	UpdateUserDirectory(ctx context.Context, userDirectory domain.TableUserDirectory) (*domain.TableUserDirectory, error)
	DeleteUserDirectory(ctx context.Context, id uuid.UUID) error
	GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]domain.TableUserTenant, error)
}
