package main

import (
	"context"
	"errors"
	"fmt"
	"minning/internal/app/usecases"
	"minning/internal/http/handlers"
	"minning/internal/http/middlewares"
	"minning/internal/http/routes"
	"minning/internal/infra/config"
	"minning/internal/infra/database"
	"minning/internal/infra/database/postgres"
	"minning/internal/infra/database/postgres/repositories"
	"minning/internal/infra/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "API:", err)
		os.Exit(1)
	}

}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load(config.DefaultFile)
	if err != nil {
		return fmt.Errorf("Config: %w", err)
	}

	fmt.Println("[DEBUG] CONFIG LOADED - Env: %s, DB Host: %s, DB Port: %d, HTTP Port: %d",
		cfg.Env,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.HTTP.Port,
	)

	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("Config: %w", err)
	}

	appLogger, err := logger.New(cfg.Log.Format, cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("Logger: %w", err)
	}

	appLogger.Info("Starting Mining Equipment API")

	dbConn, err := postgres.NewConnection(&cfg.DB)
	if err != nil {
		return fmt.Errorf("Failed to create Tx Manager: %w", err)
	}
	defer dbConn.Close()

	pgxPool := dbConn.GetPool()

	txManager, err := database.NewTransactionManager(pgxPool)
	if err != nil {
		return fmt.Errorf("Failed to create Tx Manager: %w", err)
	}

	appLogger.Info("Database connection established")

	seeder, err := postgres.NewSeederService(appLogger, txManager, &cfg.Seeder)
	if err != nil {
		return fmt.Errorf("Failed to execute seeder service: %w", err)
	}
	if err = seeder.ApplyAll(ctx); err != nil {
		return fmt.Errorf("Failed to execute all migrations with seeder service: %w", err)
	}

	appLogger.Info(fmt.Sprintf("Server log level: %v", cfg.Log.Level))

	equipmentTypeRepository := repositories.NewEquipmentTypePostgreRepository(txManager)
	equipmentRepository := repositories.NewEquipmentPostgreRepository(txManager)
	maintenanceRepository := repositories.NewMaintenancePostgreRepository(txManager)

	equipmentTypeUseCase := usecases.NewEquipmentTypeUseCase(appLogger, txManager, equipmentTypeRepository)
	equipmentUseCase := usecases.NewEquipmentUseCase(appLogger, txManager, equipmentRepository)
	maintenanceUseCase := usecases.NewMaintenanceUseCase(appLogger, txManager, maintenanceRepository)

	equipmentTypeHandler := handlers.NewnEquipmentTypeHandler(appLogger, equipmentTypeUseCase)
	equipmentHandler := handlers.NewEquipmentHandler(appLogger, equipmentUseCase)
	maintenanceHanlder := handlers.NewnMaintenanceHandler(appLogger, maintenanceUseCase)

	router := chi.NewRouter()
	router.Use(middlewares.RequestLogger(appLogger))
	router.Use(middlewares.CORS(&cfg.HTTP))
	router.Use(middlewares.Recover(appLogger))

	routes.SetupRoutes(appLogger, router,
		equipmentTypeHandler, equipmentHandler, maintenanceHanlder)

	server := &http.Server{
		Addr:         cfg.HTTP.GetAddress(),
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	shutdownErr := make(chan error, 1)
	go func() {
		appLogger.Info("API Listening", "addr", server.Addr, "env", cfg.Env)
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			shutdownErr <- err
			return
		}
		shutdownErr <- nil
	}()

	select {
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		serverErr := server.Shutdown(timeout)
		if serverErr != nil {
			return serverErr
		}
	case <-shutdownErr:
		return err
	}

	return nil
}
