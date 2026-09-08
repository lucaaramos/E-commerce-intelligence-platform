package sync

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"ecommerce-intelligence-platform/backend/internal/pricehistory"
	"ecommerce-intelligence-platform/backend/internal/product"
	"ecommerce-intelligence-platform/backend/internal/supplier"
	"ecommerce-intelligence-platform/backend/internal/supplierproduct"
)

type Worker struct {
	provider     product.ProductProvider
	interval     time.Duration
	productRepo  product.Repository
	supplierRepo supplier.Repository
	spRepo       supplierproduct.Repository
	priceRepo    pricehistory.Repository
	logger       *slog.Logger
}

func NewWorker(
	provider product.ProductProvider,
	interval time.Duration,
	productRepo product.Repository,
	supplierRepo supplier.Repository,
	spRepo supplierproduct.Repository,
	priceRepo pricehistory.Repository,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		provider:     provider,
		interval:     interval,
		productRepo:  productRepo,
		supplierRepo: supplierRepo,
		spRepo:       spRepo,
		priceRepo:    priceRepo,
		logger:       logger,
	}
}

func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("sync worker starting", "interval", w.interval)

	// Run immediately on start
	w.sync(ctx)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("sync worker stopped")
			return
		case <-ticker.C:
			w.sync(ctx)
		}
	}
}

func (w *Worker) sync(ctx context.Context) {
	w.logger.Info("sync started")

	// 1. Fetch products from provider
	mockProducts, err := w.provider.FetchProducts(ctx)
	if err != nil {
		w.logger.Error("failed to fetch products from provider", "error", err)
		return
	}
	w.logger.Info("fetched products from provider", "count", len(mockProducts))

	// 2. Get existing supplier products
	supplierProducts, err := w.spRepo.List(ctx, supplierproduct.SupplierProductFilter{})
	if err != nil {
		w.logger.Error("failed to list supplier products", "error", err)
		return
	}

	// Build map: mockSKU -> supplierProduct
	skuMap := make(map[string]*supplierproduct.SupplierProduct)
	for _, sp := range supplierProducts {
		if sp.SupplierSKU != "" {
			skuMap[sp.SupplierSKU] = sp
		}
	}

	// 3. Compare and update prices
	for _, mp := range mockProducts {
		sp, exists := skuMap[mp.SKU]
		if !exists {
			w.logger.Info("new product found in provider, skipping (not linked)", "sku", mp.SKU)
			continue
		}

		// Get latest price from history
		latestPrice, err := w.priceRepo.GetLatest(ctx, sp.ID)
		if err != nil {
			w.logger.Error("failed to get latest price", "error", err)
			continue
		}

		// Check if price changed
		priceChanged := false
		if latestPrice == nil {
			priceChanged = true
		} else if latestPrice.Price != mp.Price {
			priceChanged = true
		}

		if priceChanged {
			w.logger.Info("price changed!",
				"sku", mp.SKU,
				"name", mp.Name,
				"old_price", fmt.Sprintf("%.2f", getLatestPrice(latestPrice)),
				"new_price", fmt.Sprintf("%.2f", mp.Price),
			)

			// Record in price_history
			ph := &pricehistory.PriceHistory{
				SupplierProductID: sp.ID,
				ProductID:         sp.ProductID,
				SupplierID:        sp.SupplierID,
				Price:             mp.Price,
				ShippingCost:      sp.ShippingCost,
				Source:            w.provider.ProviderID(),
			}
			if err := w.priceRepo.Create(ctx, ph); err != nil {
				w.logger.Error("failed to create price history", "error", err)
				continue
			}

			// Update supplier product cost
			sp.Cost = mp.Price
			sp.LastSyncedAt = timePtr(time.Now())
			if err := w.spRepo.Update(ctx, sp); err != nil {
				w.logger.Error("failed to update supplier product", "error", err)
				continue
			}

			// Update product margin
			prod, err := w.productRepo.GetByID(ctx, sp.ProductID)
			if err != nil {
				w.logger.Error("failed to get product", "error", err)
				continue
			}
			if err := w.productRepo.Update(ctx, prod); err != nil {
				w.logger.Error("failed to update product", "error", err)
			}

			w.logger.Info("price synced successfully",
				"sku", mp.SKU,
				"new_cost", fmt.Sprintf("%.2f", mp.Price),
				"supplier", sp.SupplierID.Hex(),
			)
		}
	}

	w.logger.Info("sync completed")
}

func getLatestPrice(ph *pricehistory.PriceHistory) float64 {
	if ph == nil {
		return 0
	}
	return ph.Price
}

func timePtr(t time.Time) *time.Time {
	return &t
}
