package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Product struct {
	Title    string `json:"title"`
	ID       uint64 `json:"id"`
	Category string `json:"category"`
}

var products []Product
var mutex sync.RWMutex

func createProduct(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, currProduct := range products {
		if currProduct.ID == product.ID {
			http.Error(w, "Item with such ID already exists!", http.StatusBadRequest)
			return
		}
	}
	products = append(products, product)
	_ = json.NewEncoder(w).Encode(products)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for i, product := range products {
		if product.ID == id {
			products = append(products[:i], products[i+1:]...)
			return
		}
	}
	http.Error(w, "Item with given ID not found!", http.StatusBadRequest)
}

func getProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for _, product := range products {
		if product.ID == id {
			if err := json.NewEncoder(w).Encode(product); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			return
		}
	}
	http.Error(w, "Item with given ID not found!", http.StatusBadRequest)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	for i, product := range products {
		if product.ID == id {
			var newProduct Product
			if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			newProduct.ID = product.ID
			products = append(products[:i], products[i+1:]...)
			products = append(products, newProduct)
			_ = json.NewEncoder(w).Encode(product)
			return
		}
	}
	http.Error(w, "Item with given ID not found!", http.StatusBadRequest)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products", createProduct).Methods("POST")
	r.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}
