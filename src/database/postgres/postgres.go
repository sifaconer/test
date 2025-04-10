package postgres

import (
	"api-test/src/common"
	"api-test/src/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Database interface {
	Connect(tenantID uuid.UUID, dsn string) (*bun.DB, error)
}

type postgres struct {
	log  common.Logger
	conf *config.Config
}

func (p *postgres) Connect(tenantID uuid.UUID, dsn string) (*bun.DB, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	bunDB := bun.NewDB(sqldb, pgdialect.New())
	bunDB.AddQueryHook(newTenantQueryHook(tenantID, p.log))
	if p.conf.IsDev() {
		bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	err := bunDB.Ping()
	if err != nil {
		p.log.Error(context.Background(), "Failed to ping database", "error", err)
		return nil, err
	}
	return bunDB, nil
}

type tenantQueryHook struct {
	tenantID uuid.UUID
	log      common.Logger
}

func newTenantQueryHook(tenantID uuid.UUID, log common.Logger) *tenantQueryHook {
	return &tenantQueryHook{tenantID: tenantID, log: log}
}

func (h *tenantQueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	if event.Model == nil {
		return ctx
	}
	table := event.Model.(bun.TableModel).Table().Name
	h.log.Info(ctx, fmt.Sprintf("Table:%s", table), "tenant_id", h.tenantID)
	return ctx
}

func (h *tenantQueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	if event.Model == nil {
		return
	}
	table := event.Model.(bun.TableModel).Table().Name
	if event.Err != nil {
		h.log.Error(ctx, fmt.Sprintf("Table:%s", table), "tenant_id", h.tenantID, "error", event.Err)
	}
}

var _ bun.QueryHook = (*tenantQueryHook)(nil)

func NewPostgres(log common.Logger, conf *config.Config) Database {
	return &postgres{
		log:  log,
		conf: conf,
	}
}

var _ Database = (*postgres)(nil)
