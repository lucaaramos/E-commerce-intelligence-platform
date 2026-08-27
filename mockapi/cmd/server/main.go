package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Product struct {
	ID          string  `json:"id"`
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	UpdatedAt   string  `json:"updatedAt"`
}

var (
	products = map[string]*Product{
		"1": {ID: "1", SKU: "AUD-001", Name: "Wireless Headphones", Price: 32.00, Stock: 150},
		"2": {ID: "2", SKU: "SPK-002", Name: "Bluetooth Speaker", Price: 24.50, Stock: 200},
		"3": {ID: "3", SKU: "WCH-003", Name: "Smart Watch", Price: 89.00, Stock: 75},
		"4": {ID: "4", SKU: "HUB-004", Name: "USB-C Hub", Price: 18.00, Stock: 300},
		"5": {ID: "5", SKU: "SHO-005", Name: "Running Shoes", Price: 45.00, Stock: 120},
	}
	mu sync.RWMutex
)

func main() {
	http.HandleFunc("GET /products", handleList)
	http.HandleFunc("GET /products/", handleGet)
	http.HandleFunc("PUT /products/", handleUpdate)
	http.HandleFunc("POST /products/", handleUpdatePrice)

	fmt.Println("Mock Supplier API running on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func handleList(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	list := make([]*Product, 0, len(products))
	for _, p := range products {
		list = append(list, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	id := r.URL.Path[len("/products/"):]
	p, ok := products[id]
	if !ok {
		http.Error(w, "not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	id := r.URL.Path[len("/products/"):]
	p, ok := products[id]
	if !ok {
		http.Error(w, "not found", 404)
		return
	}

	var input struct {
		Price *float64 `json:"price"`
		Stock *int     `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if input.Price != nil {
		old := p.Price
		p.Price = *input.Price
		log.Printf("Price changed: %s %s $%.2f -> $%.2f", p.SKU, p.Name, old, *input.Price)
	}
	if input.Stock != nil {
		p.Stock = *input.Stock
	}
	p.UpdatedAt = time.Now().Format(time.RFC3339)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func handleUpdatePrice(w http.ResponseWriter, r *http.Request) {
	handleUpdate(w, r)
}
