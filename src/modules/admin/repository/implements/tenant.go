package implements

import (
	"api-test/src/common"
	"api-test/src/modules/admin/domain"
	"api-test/src/modules/admin/repository"
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type tenantRepository struct {
	log    common.Logger
	tenant *common.TenantConnectionManager
}

// Manual rollback
func (t *tenantRepository) rollback(ctx context.Context, db *bun.DB, tenant domain.TableTenant) {
	// Delete the database
	dbQuery := "DROP DATABASE IF EXISTS ?"
	_, _ = db.ExecContext(ctx, dbQuery, bun.Safe(tenant.DBName))
	// Delete the user
	userQuery := "DROP USER IF EXISTS ?"
	_, _ = db.ExecContext(ctx, userQuery, bun.Safe(tenant.DBUser))
}

// CreateTenant implements repository.TenantRepository.
func (t *tenantRepository) CreateTenant(ctx context.Context, tenant domain.TableTenant) (*domain.TableTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	// 1. Crear el nuevo usuario
	userQuery := "CREATE USER ? WITH PASSWORD ?"
	_, err = db.ExecContext(ctx, userQuery, bun.Safe(tenant.DBUser), tenant.PasswordPlaintext)
	if err != nil {
		t.rollback(ctx, db, tenant)
		return nil, common.CheckDBErrorType(err)
	}

	// 2. Crear la base de datos (fuera de la transacción)
	dbQuery := "CREATE DATABASE ? WITH OWNER ?"
	_, err = db.ExecContext(ctx, dbQuery, bun.Safe(tenant.DBName), bun.Safe(tenant.DBUser))
	if err != nil {
		t.rollback(ctx, db, tenant)
		return nil, common.CheckDBErrorType(err)
	}

	// 3. Otorgar privilegios
	privilegeQuery := "GRANT ALL PRIVILEGES ON DATABASE ? TO ?"
	_, err = db.ExecContext(ctx, privilegeQuery, bun.Safe(tenant.DBName), bun.Safe(tenant.DBUser))
	if err != nil {
		t.rollback(ctx, db, tenant)
		return nil, common.CheckDBErrorType(err)
	}

	// 4. Insertar el nuevo tenant en la base de datos
	_, err = db.NewInsert().Model(&tenant).Exec(ctx)
	if err != nil {
		t.rollback(ctx, db, tenant)
		return nil, common.CheckDBErrorType(err)
	}

	return &tenant, nil
}

// DeleteTenant implements repository.TenantRepository.
func (t *tenantRepository) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return err
	}

	var model domain.TableTenant
	_, err = db.NewDelete().Model(&model).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return common.CheckDBErrorType(err)
	}

	// Delete the database
	t.rollback(ctx, db, model)

	return nil
}

// GetTenantByID implements repository.TenantRepository.
func (t *tenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*domain.TableTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var tenant domain.TableTenant
	err = db.NewSelect().Model(&tenant).Where("tenant_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return &tenant, nil
}

// GetTenantByName implements repository.TenantRepository.
func (t *tenantRepository) GetTenantByName(ctx context.Context, name string) (*domain.TableTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var tenant domain.TableTenant
	err = db.NewSelect().Model(&tenant).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return &tenant, nil
}

// UpdateTenant implements repository.TenantRepository.
func (t *tenantRepository) UpdateTenant(ctx context.Context, tenant domain.TableTenant) (*domain.TableTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	_, err = db.NewUpdate().Model(&tenant).Where("id = ?", tenant.ID).Exec(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return &tenant, nil
}

// GetAllTenants implements repository.TenantRepository.
func (t *tenantRepository) GetAllTenants(ctx context.Context) ([]domain.TableTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var tenants []domain.TableTenant
	err = db.NewSelect().Model(&tenants).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return tenants, nil
}

// CreateUserTenant implements repository.TenantRepository.
func (t *tenantRepository) CreateUserTenant(ctx context.Context, userTenant domain.TableUserTenant) (*domain.TableUserTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	_, err = db.NewInsert().Model(&userTenant).Exec(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return &userTenant, nil
}

// GetTenantsByUser implements repository.TenantRepository.
func (t *tenantRepository) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]domain.TableUserTenant, error) {
	db, err := t.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var tenants []domain.TableUserTenant
	err = db.NewSelect().Model(&tenants).Relation("Tenant").Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}
	return tenants, nil
}

func NewTenantRepository(log common.Logger, tenant *common.TenantConnectionManager) repository.TenantRepository {
	return &tenantRepository{
		log:    log,
		tenant: tenant,
	}
}

var _ repository.TenantRepository = (*tenantRepository)(nil)
