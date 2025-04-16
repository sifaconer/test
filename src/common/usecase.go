package common

import (
	"api-test/src/config"
	"context"

	"github.com/uptrace/bun"
)

type UseCase[CreateDTO, ResponseDTO, UpdateDTO, ID any] interface {
	// CRUD
	Create(ctx context.Context, dto CreateDTO) (*ResponseDTO, error)
	GetById(ctx context.Context, id ID) (*ResponseDTO, error)
	Update(ctx context.Context, id ID, dto UpdateDTO) (*ResponseDTO, error)
	Delete(ctx context.Context, id ID) error
	// Search
	Search(ctx context.Context, filters *QueryParams) ([]ResponseDTO, error)
	// Bulk
	CreateMany(ctx context.Context, dtos []CreateDTO) ([]ResponseDTO, error)
	UpdateMany(ctx context.Context, dtos []UpdateDTO) ([]ResponseDTO, error)
	DeleteMany(ctx context.Context, ids []ID) error
	// Transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error
	CreateTx(ctx context.Context, tx bun.Tx, dto CreateDTO) (*ResponseDTO, error)
	UpdateTx(ctx context.Context, tx bun.Tx, id ID, dto UpdateDTO) (*ResponseDTO, error)
	DeleteTx(ctx context.Context, tx bun.Tx, id ID) error
	// Bulk Transaction
	CreateManyTx(ctx context.Context, tx bun.Tx, dtos []CreateDTO) ([]ResponseDTO, error)
	UpdateManyTx(ctx context.Context, tx bun.Tx, dtos []UpdateDTO) ([]ResponseDTO, error)
	DeleteManyTx(ctx context.Context, tx bun.Tx, ids []ID) error
}

type usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID any] struct {
	config *config.Config
	log Logger
	tenant *TenantConnectionManager
	repo Repository[Table, ID]
	createToTable func(CreateDTO) Table
	updateToTable func(UpdateDTO) Table
	toResponseDTO func(Table) ResponseDTO
}

func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) Create(ctx context.Context, dto CreateDTO) (*ResponseDTO, error) {
	table, err := u.repo.Create(ctx, u.createToTable(dto))
	if err != nil {
		return nil, err
	}

	result := u.toResponseDTO(*table)
	return &result, nil
}

func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) GetById(ctx context.Context, id ID) (*ResponseDTO, error) {
	table, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	result := u.toResponseDTO(*table)
	return &result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) Update(ctx context.Context, id ID, dto UpdateDTO) (*ResponseDTO, error) {
	table, err := u.repo.Update(ctx, id, u.updateToTable(dto))
	if err != nil {
		return nil, err
	}
	result := u.toResponseDTO(*table)
	return &result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) Delete(ctx context.Context, id ID) error {
	return u.repo.Delete(ctx, id)
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) Search(ctx context.Context, filters *QueryParams) ([]ResponseDTO, error) {
	tables, err := u.repo.Search(ctx, filters)
	if err != nil {
		return nil, err
	}
	dtos := make([]ResponseDTO, len(tables))
	for i, table := range tables {
		dtos[i] = u.toResponseDTO(table)
	}
	return dtos, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) CreateMany(ctx context.Context, dtos []CreateDTO) ([]ResponseDTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.createToTable(dto)
	}
	tablesResult, err := u.repo.CreateMany(ctx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]ResponseDTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toResponseDTO(table)
	}
	return result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) UpdateMany(ctx context.Context, dtos []UpdateDTO) ([]ResponseDTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.updateToTable(dto)
	}
	tablesResult, err := u.repo.UpdateMany(ctx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]ResponseDTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toResponseDTO(table)
	}
	return result, nil
}

func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) DeleteMany(ctx context.Context, ids []ID) error {
	return u.repo.DeleteMany(ctx, ids)
}

func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx bun.Tx) error) error {
	return u.repo.WithTransaction(ctx, fn)
}

func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) CreateTx(ctx context.Context, tx bun.Tx, dto CreateDTO) (*ResponseDTO, error) {
	table, err := u.repo.CreateTx(ctx, tx, u.createToTable(dto))
	if err != nil {
		return nil, err
	}
	result := u.toResponseDTO(*table)
	return &result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) UpdateTx(ctx context.Context, tx bun.Tx, id ID, dto UpdateDTO) (*ResponseDTO, error) {
	table, err := u.repo.UpdateTx(ctx, tx, id, u.updateToTable(dto))
	if err != nil {
		return nil, err
	}
	result := u.toResponseDTO(*table)
	return &result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) DeleteTx(ctx context.Context, tx bun.Tx, id ID) error {
	return u.repo.DeleteTx(ctx, tx, id)
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) CreateManyTx(ctx context.Context, tx bun.Tx, dtos []CreateDTO) ([]ResponseDTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.createToTable(dto)
	}
	tablesResult, err := u.repo.CreateManyTx(ctx, tx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]ResponseDTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toResponseDTO(table)
	}
	return result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) UpdateManyTx(ctx context.Context, tx bun.Tx, dtos []UpdateDTO) ([]ResponseDTO, error) {
	tables := make([]Table, len(dtos))
	for i, dto := range dtos {
		tables[i] = u.updateToTable(dto)
	}
	tablesResult, err := u.repo.UpdateManyTx(ctx, tx, tables)
	if err != nil {
		return nil, err
	}
	result := make([]ResponseDTO, len(tablesResult))
	for i, table := range tablesResult {
		result[i] = u.toResponseDTO(table)
	}
	return result, nil
}
	
func (u *usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]) DeleteManyTx(ctx context.Context, tx bun.Tx, ids []ID) error {
	return u.repo.DeleteManyTx(ctx, tx, ids)
}

func NewUseCase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID any](config *config.Config, log Logger, tenant *TenantConnectionManager, repo Repository[Table, ID], createToTable func(CreateDTO) Table, updateToTable func(UpdateDTO) Table, toResponseDTO func(Table) ResponseDTO) UseCase[CreateDTO, ResponseDTO, UpdateDTO, ID] {
	return &usecase[CreateDTO, ResponseDTO, UpdateDTO, Table, ID]{
		config: config,
		log: log,
		tenant: tenant,
		repo: repo,
		createToTable: createToTable,
		updateToTable: updateToTable,
		toResponseDTO: toResponseDTO,
	}
}