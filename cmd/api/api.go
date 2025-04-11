package api

import (
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/database/postgres"
	adminAPI "api-test/src/modules/admin/api"
	"api-test/src/modules/admin/usecase"
	"api-test/src/modules/carritocompra/api"
	"context"
	"fmt"
	"net"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

type Rest struct {
	conf          *config.Config
	log           common.Logger
	tenant        *common.TenantConnectionManager
	psql          postgres.Database
	migrations    usecase.TenantMigrations
	EXCLUDE_PATHS []string
}

func NewRest(
	conf *config.Config,
	log common.Logger,
	tenant *common.TenantConnectionManager,
	psql postgres.Database,
	migrations usecase.TenantMigrations,
) *Rest {
	return &Rest{
		conf:       conf,
		log:        log,
		tenant:     tenant,
		psql:       psql,
		migrations: migrations,
		EXCLUDE_PATHS: []string{
			"/login",
			"/register",
			"/logout",
			"/refresh",
		},
	}
}

func (r *Rest) Run() {
	r.log.Info(context.Background(), "Starting Rest API")
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
	})
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	app.Use(r.ErrorHandler())
	app.Use(helmet.New()) // Helmet middleware to secure default headers
	app.Use(r.RequestLimitMiddleware())
	app.Use(r.CORSMiddleware())
	app.Use(r.LoggerMiddleware())
	app.Use(r.FieldMiddleware())
	app.Use(r.AuthenticationMiddleware())
	app.Use(r.TenantMiddleware())
	// app.Use(r.AuthorizationMiddleware()) // TODO: pendiente definir método de manejo de permisos
	// app.Use(r.FilterMiddleware()) // TODO: pendiente definir método de manejo de filtros que llegan a sql para consultas dinámicas

	admin := adminAPI.NewAdminAPI(r.log, app, r.conf, r.tenant, r.migrations, r.psql)
	if err := admin.RegisterAllTenants(context.Background()); err != nil {
		r.log.Error(context.Background(), "Error registering all tenants", "error", err)
		// os.Exit(1) // No detener la aplicación si falla un tenant, priorizar el inicio de la API
	}
	admin.Register()
	api.NewCarritoCompraAPI(r.log, app, r.conf, r.tenant).Register()

	r.log.Info(context.Background(), "Rest API started")
	host := net.JoinHostPort("0.0.0.0", fmt.Sprintf("%d", r.conf.Port))
	if err := app.Listen(host); err != nil {
		r.log.Error(context.Background(), "Error starting Rest API", "error", err)
		os.Exit(1)
	}
}
