package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce-intelligence-platform/backend/internal/product"
	"ecommerce-intelligence-platform/backend/internal/supplier"
	"ecommerce-intelligence-platform/backend/pkg/config"
	"ecommerce-intelligence-platform/backend/pkg/httpx"
	"ecommerce-intelligence-platform/backend/pkg/logging"
	"ecommerce-intelligence-platform/backend/pkg/mongodb"
	"ecommerce-intelligence-platform/backend/pkg/rabbitmq"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()
	logger := logging.New(cfg.Log)

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	bootCtx, cancel := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancel()

	db, err := mongodb.Connect(bootCtx, cfg.MongoDB.URI, cfg.MongoDB.DB)
	if err != nil {
		return fmt.Errorf("connect mongodb: %w", err)
	}
	defer db.Close(context.Background())
	logger.Info("mongodb connected", "db", cfg.MongoDB.DB)

	mq, err := rabbitmq.Connect(bootCtx, cfg.RabbitMQ.URI)
	if err != nil {
		return fmt.Errorf("connect rabbitmq: %w", err)
	}
	defer mq.Close()
	logger.Info("rabbitmq connected")

	engine, api := httpx.New(&cfg, logger)
	httpx.RegisterReadiness(engine, db, mq)

	// Repositories
	productRepo := product.NewRepository(db.DB())
	supplierRepo := supplier.NewRepository(db.DB())

	// Ensure indexes
	if err := productRepo.EnsureIndexes(bootCtx); err != nil {
		return fmt.Errorf("product indexes: %w", err)
	}
	if err := supplierRepo.EnsureIndexes(bootCtx); err != nil {
		return fmt.Errorf("supplier indexes: %w", err)
	}

	// Services
	productSvc := product.NewService(productRepo)
	supplierSvc := supplier.NewService(supplierRepo)

	// Handlers
	productHandler := product.NewHandler(productSvc)
	supplierHandler := supplier.NewHandler(supplierSvc)

	// Register routes
	productHandler.RegisterRoutes(api)
	supplierHandler.RegisterRoutes(api)

	server := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           engine,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("http server starting", "addr", cfg.Addr())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-rootCtx.Done():
		logger.Info("shutdown signal received")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}
	logger.Info("http server stopped gracefully")
	return nil
}
