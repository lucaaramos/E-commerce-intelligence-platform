package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"ecommerce-intelligence-platform/backend/internal/product"
	"ecommerce-intelligence-platform/backend/internal/supplierproduct"
	"ecommerce-intelligence-platform/backend/pkg/config"
	"ecommerce-intelligence-platform/backend/pkg/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var dummies = []struct {
	sku      string
	name     string
	cost     float64
	price    float64
	shipping float64
}{
	{"BEA-ESS-ESS-001", "Essence Mascara Lash Princess", 5.00, 10.00, 2.00},
	{"FRA-CAL-CAL-006", "Eau de Parfum Classic", 8.00, 15.00, 3.00},
	{"FUR-ANN-ANN-011", "Faux Fur Blanket", 12.00, 25.00, 5.00},
	{"GRO-BRD-APP-016", "Apple Fruit Basket", 8.00, 18.00, 2.50},
	{"GRO-BRD-BEE-017", "Bee Honey Jar", 15.00, 22.00, 3.50},
}

func main() {
	ctx := context.Background()

	cfg := config.Load()

	mongoClient, err := mongodb.Connect(ctx, cfg.MongoDB.URI, cfg.MongoDB.DB)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer mongoClient.Close(ctx)

	db := mongoClient.DB()

	productRepo := product.NewRepository(db)
	supplierProductRepo := supplierproduct.NewRepository(db)

	rand.Seed(time.Now().UnixNano())

	fmt.Println("Seed starting...")

	for _, d := range dummies {
		existing, err := productRepo.GetBySKU(ctx, d.sku)
		if err == nil && existing != nil {
			fmt.Printf("⏭️ Product already exists: %s (ID: %s)\n", d.sku, existing.ID.Hex())

			sp := &supplierproduct.SupplierProduct{
				ProductID:    existing.ID,
				SupplierID:   generateSupplierID(),
				SupplierSKU:  d.sku,
				Cost:         d.cost,
				ShippingCost: d.shipping,
				LeadTimeDays: rand.Intn(30) + 1,
				MinOrderQty:  1,
				IsActive:     true,
				LastSyncedAt: nil,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := supplierProductRepo.Create(ctx, sp); err != nil {
				log.Printf("failed to create supplier product: %v", err)
			} else {
				fmt.Printf("Supplier product created for: %s\n", d.sku)
			}
			continue
		}

		p := &product.Product{
			ID:           bson.NewObjectID(),
			SKU:          d.sku,
			Name:         d.name,
			Status:       product.StatusActive,
			SellingPrice: d.price,
			RealCost:     d.cost,
			Margin:       ((d.price - d.cost) / d.price) * 100,
			Variants:     nil,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := productRepo.Create(ctx, p); err != nil {
			log.Printf("failed to create product %s: %v", d.sku, err)
			continue
		}
		fmt.Printf("🆕 Product created: %s (ID: %s)\n", d.sku, p.ID.Hex())

		sp := &supplierproduct.SupplierProduct{
			ProductID:    p.ID,
			SupplierID:   generateSupplierID(),
			SupplierSKU:  d.sku,
			Cost:         d.cost,
			ShippingCost: d.shipping,
			LeadTimeDays: rand.Intn(30) + 1,
			MinOrderQty:  1,
			IsActive:     true,
			LastSyncedAt: nil,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		if err := supplierProductRepo.Create(ctx, sp); err != nil {
			log.Printf("failed to create supplier product for %s: %v", d.sku, err)
		} else {
			fmt.Printf("Supplier product created for: %s\n", d.sku)
		}
	}

	fmt.Println("\n🌱 Seed completed!")
}

func generateSupplierID() bson.ObjectID {
	return bson.NewObjectID()
}
