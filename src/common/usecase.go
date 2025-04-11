package common

import (
	"api-test/src/config"
	"context"

	"github.com/uptrace/bun"
)

type UseCase[DTO any, ID any] interface {
	// CRUD
	Create(ctx context.Context, dto DTO) (*DTO, error)
	GetById(ctx context.Context, id ID) (*DTO, error)
	Update(ctx context.Context, id ID, dto DTO) (*DTO, error)
	Delete(ctx context.Context, id ID) error
	// Search
	Search(ctx context.Context, filters *QueryParams) ([]DTO, error)
	// Bulk
	CreateMany(ctx context.Context, dtos []DTO) ([]DTO, error)
	UpdateMany(ctx context.Context, dtos []DTO) ([]DTO, error)
	DeleteMany(ctx context.Context, ids []ID) error
	// Transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
	CreateTx(ctx context.Context, tx bun.Tx, dto DTO) (*DTO, error)
	UpdateTx(ctx context.Context, tx bun.Tx, id ID, dto DTO) (*DTO, error)
	DeleteTx(ctx context.Context, tx bun.Tx, id ID) error
	// Bulk Transaction
	CreateManyTx(ctx context.Context, tx bun.Tx, dtos []DTO) ([]DTO, error)
	UpdateManyTx(ctx context.Context, tx bun.Tx, dtos []DTO) ([]DTO, error)
	DeleteManyTx(ctx context.Context, tx bun.Tx, ids []ID) error
}

type usecase[DTO any, Table any, ID any] struct {
	config *config.Config
	log Logger
	tenant *TenantConnectionManager
	repo Repository[Table, ID]
	toTable func(DTO) Table
	toDTO func(Table) DTO
}

func (u *usecase[DTO, Table, ID]) Create(ctx context.Context, dto DTO) (*DTO, error) {
	table, err := u.repo.Create(ctx, u.toTable(dto))
	if err != nil {
		return nil, err
	}

	result := u.toDTO(*table)
	return &result, nil
}

func (u *usecase[DTO, Table, ID]) GetById(ctx context.Context, id ID) (*DTO, error) {
	table, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	result := u.toDTO(*table)
	return &result, nil
}
	
func (u *usecase[DTO, Table, ID]) Update(ctx context.Context, id ID, dto DTO) (*DTO, error) {
	table, err := u.repo.Update(ctx, id, u.toTable(dto))
	if err != nil {
		return nil, err
	}
	result := u.toDTO(*table)
	return &result, nil
}
	
func (u *usecase[DTO, Table, ID]) Delete(ctx context.Context, id ID) error {
	return u.repo.Delete(ctx, id)
}
	
func (u *usecase[DTO, Table, ID]) Search(ctx context.Context, filters *QueryParams) ([]DTO, error) {
	tables, err := u.repo.Search(ctx, filters)
	if err != nil {
		return nil, err
	}
	dtos := make([]DTO, len(tables))
	for i, table := range tables {
		dtos[i] = u.toDTO(table)
	}
	return dtos, nil
}
	
func (u *usecase[DTO, Table, ID]) CreateMany(ctx context.Context, dtos []DTO) ([]DTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.toTable(dto)
	}
	tablesResult, err := u.repo.CreateMany(ctx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]DTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toDTO(table)
	}
	return result, nil
}
	
func (u *usecase[DTO, Table, ID]) UpdateMany(ctx context.Context, dtos []DTO) ([]DTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.toTable(dto)
	}
	tablesResult, err := u.repo.UpdateMany(ctx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]DTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toDTO(table)
	}
	return result, nil
}

func (u *usecase[DTO, Table, ID]) DeleteMany(ctx context.Context, ids []ID) error {
	return u.repo.DeleteMany(ctx, ids)
}

func (u *usecase[DTO, Table, ID]) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return u.repo.WithTransaction(ctx, fn)
}

func (u *usecase[DTO, Table, ID]) CreateTx(ctx context.Context, tx bun.Tx, dto DTO) (*DTO, error) {
	table, err := u.repo.CreateTx(ctx, tx, u.toTable(dto))
	if err != nil {
		return nil, err
	}
	result := u.toDTO(*table)
	return &result, nil
}
	
func (u *usecase[DTO, Table, ID]) UpdateTx(ctx context.Context, tx bun.Tx, id ID, dto DTO) (*DTO, error) {
	table, err := u.repo.UpdateTx(ctx, tx, id, u.toTable(dto))
	if err != nil {
		return nil, err
	}
	result := u.toDTO(*table)
	return &result, nil
}
	
func (u *usecase[DTO, Table, ID]) DeleteTx(ctx context.Context, tx bun.Tx, id ID) error {
	return u.repo.DeleteTx(ctx, tx, id)
}
	
func (u *usecase[DTO, Table, ID]) CreateManyTx(ctx context.Context, tx bun.Tx, dtos []DTO) ([]DTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.toTable(dto)
	}
	tablesResult, err := u.repo.CreateManyTx(ctx, tx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]DTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toDTO(table)
	}
	return result, nil
}
	
func (u *usecase[DTO, Table, ID]) UpdateManyTx(ctx context.Context, tx bun.Tx, dtos []DTO) ([]DTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.toTable(dto)
	}
	tablesResult, err := u.repo.UpdateManyTx(ctx, tx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]DTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toDTO(table)
	}
	return result, nil
}
	
func (u *usecase[DTO, Table, ID]) DeleteManyTx(ctx context.Context, tx bun.Tx, ids []ID) error {
	return u.repo.DeleteManyTx(ctx, tx, ids)
}

func NewUseCase[DTO any, Table any, ID any](config *config.Config, log Logger, tenant *TenantConnectionManager, repo Repository[Table, ID], toTable func(DTO) Table, toDTO func(Table) DTO) UseCase[DTO, ID] {
	return &usecase[DTO, Table, ID]{
		config: config,
		log: log,
		tenant: tenant,
		repo: repo,
		toTable: toTable,
		toDTO: toDTO,
	}
}