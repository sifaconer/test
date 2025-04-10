package main

import (
	"api-test/cmd/api"
	"api-test/cmd/banner"
	"api-test/cmd/database"
	"api-test/src/common"
	"api-test/src/config"
	"api-test/src/database/postgres"
	"api-test/src/modules/admin/usecase"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := config.NewConfig()
	log := common.NewLogger()
	if err := conf.Load(); err != nil {
		log.Error(context.Background(), "Error loading config")
		for _, e := range err {
			log.Error(context.Background(), e.Error())
		}
		os.Exit(1)
	}
	tenantManager := common.NewTenantConnectionManager(conf)
	log = common.NewLoggerWithTenantManager(tenantManager)

	banner.Banner(conf)
	database := database.NewDatabase(
		conf,
		log,
		tenantManager,
		postgres.Admins,
		postgres.Tenants,
		postgres.Common)
	if err := database.Run(); err != nil {
		log.Error(context.Background(), "Error starting database", "error", err)
		os.Exit(1)
	}

	migrations := usecase.NewTenantMigrations(log, conf, tenantManager,
		postgres.Admins,
		postgres.Tenants,
		postgres.Common)

	if conf.IsDev() {
		log.Warn(context.Background(), "Starting API in development mode")
	}
	go func() {
		api.NewRest(conf, log, tenantManager, database.PSQL(), migrations).Run()
	}()

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := database.Stop(); err != nil {
		log.Error(ctx, "Error stopping database", "error", err)
	}
	log.Info(ctx, "Shutting down")
}
