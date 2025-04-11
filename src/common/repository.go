package common

import (
	"api-test/src/config"
	"context"

	"github.com/uptrace/bun"
)

type Repository[Table any, ID any] interface {
	// CRUD
	Create(ctx context.Context, item Table) (*Table, error)
	GetById(ctx context.Context, id ID, relations ...string) (*Table, error)
	Update(ctx context.Context, id ID, item Table) (*Table, error)
	Delete(ctx context.Context, id ID) error

	// Bulk
	CreateMany(ctx context.Context, items []Table) ([]Table, error)
	UpdateMany(ctx context.Context, items []Table) ([]Table, error)
	DeleteMany(ctx context.Context, ids []ID) error

	// Transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
	CreateTx(ctx context.Context, tx bun.Tx, item Table) (*Table, error)
	UpdateTx(ctx context.Context, tx bun.Tx, id ID, item Table) (*Table, error)
	DeleteTx(ctx context.Context, tx bun.Tx, id ID) error
	CreateManyTx(ctx context.Context, tx bun.Tx, items []Table) ([]Table, error)
	UpdateManyTx(ctx context.Context, tx bun.Tx, items []Table) ([]Table, error)
	DeleteManyTx(ctx context.Context, tx bun.Tx, ids []ID) error

	// Search
	Search(ctx context.Context, filters *QueryParams, relations ...string) ([]Table, error)
}

type repository[Table any, ID any] struct {
	config *config.Config
	log Logger
	tenant *TenantConnectionManager

}

func (r *repository[Table, ID]) Create(ctx context.Context, item Table) (*Table, error) {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewInsert().Model(&item).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}

	return &item, nil
}

func (r *repository[Table, ID]) GetById(ctx context.Context, id ID, relations ...string) (*Table, error) {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	var item Table
	q := db.NewSelect().Model(&item)
	for _, relation := range relations {
		q = q.Relation(relation)
	}

	err = q.Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return &item, nil
}

func (r *repository[Table, ID]) Update(ctx context.Context, id ID, item Table) (*Table, error) {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewUpdate().Model(&item).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return &item, nil
}

func (r *repository[Table, ID]) Delete(ctx context.Context, id ID) error {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return err
	}

	var item Table
	_, err = db.NewDelete().Model(&item).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return CheckDBErrorType(err)
	}
	return nil
}

func (r *repository[Table, ID]) Search(ctx context.Context, filters *QueryParams, relations ...string) ([]Table, error) {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	var items []Table
	q := db.NewSelect().Model(&items)
	for _, relation := range relations {
		q = q.Relation(relation)
	}

	if filters != nil {
		q = q.Where(filters.String()) // TODO: Handle filters
	}

	err = q.Scan(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return items, nil
}

// Bulk
func (r *repository[Table, ID]) CreateMany(ctx context.Context, items []Table) ([]Table, error) {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewInsert().Model(&items).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return items, nil
}

func (r *repository[Table, ID]) UpdateMany(ctx context.Context, items []Table) ([]Table, error) {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = db.NewUpdate().Model(&items).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return items, nil
}

func (r *repository[Table, ID]) DeleteMany(ctx context.Context, ids []ID) error {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return err
	}

	var item Table
	_, err = db.NewDelete().Model(&item).Where("id IN (?)", bun.In(ids)).Exec(ctx)
	if err != nil {
		return CheckDBErrorType(err)
	}
	return nil
}

// Transaction
func (r *repository[Table, ID]) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	db, err := r.tenant.GetDBContext(ctx)
	if err != nil {
		return err
	}

	return db.RunInTx(ctx, nil, fn)
}

func (r *repository[Table, ID]) CreateManyTx(ctx context.Context, tx bun.Tx, items []Table) ([]Table, error) {
	_, err := tx.NewInsert().Model(&items).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return items, nil
}

func (r *repository[Table, ID]) UpdateManyTx(ctx context.Context, tx bun.Tx, items []Table) ([]Table, error) {
	_, err := tx.NewUpdate().Model(&items).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return items, nil
}

func (r *repository[Table, ID]) DeleteManyTx(ctx context.Context, tx bun.Tx, ids []ID) error {
	var item Table
	_, err := tx.NewDelete().Model(&item).Where("id IN (?)", bun.In(ids)).Exec(ctx)
	if err != nil {
		return CheckDBErrorType(err)
	}
	return nil
}

func (r *repository[Table, ID]) CreateTx(ctx context.Context, tx bun.Tx, item Table) (*Table, error) {
	_, err := tx.NewInsert().Model(&item).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return &item, nil
}

func (r *repository[Table, ID]) UpdateTx(ctx context.Context, tx bun.Tx, id ID, item Table) (*Table, error) {
	_, err := tx.NewUpdate().Model(&item).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return nil, CheckDBErrorType(err)
	}
	return &item, nil
}

func (r *repository[Table, ID]) DeleteTx(ctx context.Context, tx bun.Tx, id ID) error {
	var item Table
	_, err := tx.NewDelete().Model(&item).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return CheckDBErrorType(err)
	}
	return nil
}

func NewRepository[Table any, ID any](config *config.Config, log Logger, tenant *TenantConnectionManager) Repository[Table, ID] {
	return &repository[Table, ID]{
		config: config,
		log: log,
		tenant: tenant,
	}
}


	
