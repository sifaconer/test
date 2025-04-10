package implements

import (
	"api-test/src/common"
	"api-test/src/modules/admin/domain"
	"api-test/src/modules/admin/repository"
	"context"

	"github.com/google/uuid"
)

type user struct {
	log    common.Logger
	tenant *common.TenantConnectionManager
}

// CreateUserDirectory implements repository.UserDirectoryRepository.
func (u *user) CreateUserDirectory(ctx context.Context, userDirectory domain.TableUserDirectory) (*domain.TableUserDirectory, error) {
	db, err := u.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	_, err = db.NewInsert().Model(&userDirectory).Exec(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}

	return &userDirectory, nil
}

// DeleteUserDirectory implements repository.UserDirectoryRepository.
func (u *user) DeleteUserDirectory(ctx context.Context, id uuid.UUID) error {
	db, err := u.tenant.GetKosviTenantDB()
	if err != nil {
		return err
	}

	_, err = db.NewDelete().Model(&domain.TableUserDirectory{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return common.CheckDBErrorType(err)
	}

	return nil
}

// GetUserDirectoryByEmail implements repository.UserDirectoryRepository.
func (u *user) GetUserDirectoryByEmail(ctx context.Context, email string) (*domain.TableUserDirectory, error) {
	db, err := u.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var userDirectory domain.TableUserDirectory
	err = db.NewSelect().Model(&userDirectory).Relation("UserTenants").Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}

	return &userDirectory, nil
}

// GetUserDirectoryByID implements repository.UserDirectoryRepository.
func (u *user) GetUserDirectoryByID(ctx context.Context, id uuid.UUID) (*domain.TableUserDirectory, error) {
	db, err := u.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var userDirectory domain.TableUserDirectory
	err = db.NewSelect().Model(&userDirectory).Relation("UserTenants").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}

	return &userDirectory, nil
}

// UpdateUserDirectory implements repository.UserDirectoryRepository.
func (u *user) UpdateUserDirectory(ctx context.Context, userDirectory domain.TableUserDirectory) (*domain.TableUserDirectory, error) {
	db, err := u.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	_, err = db.NewUpdate().Model(&userDirectory).Where("id = ?", userDirectory.ID).Exec(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}

	return &userDirectory, nil
}

// GetTenantsByUser implements repository.UserDirectoryRepository.
func (u *user) GetTenantsByUser(ctx context.Context, userID uuid.UUID) ([]domain.TableUserTenant, error) {
	db, err := u.tenant.GetKosviTenantDB()
	if err != nil {
		return nil, err
	}

	var tenants []domain.TableUserTenant
	err = db.NewSelect().Model(&tenants).Where("user_id = ?", userID).Scan(ctx)
	if err != nil {
		return nil, common.CheckDBErrorType(err)
	}

	return tenants, nil
}

func NewUserRepository(log common.Logger, tenant *common.TenantConnectionManager) repository.UserDirectoryRepository {
	return &user{
		log:    log,
		tenant: tenant,
	}
}

var _ repository.UserDirectoryRepository = (*user)(nil)
