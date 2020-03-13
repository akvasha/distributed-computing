package main

import (
	"encoding/json"
	"github.com/akvasha/distributed-computing/dbClient"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var db dbClient.Client

func dbCreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product dbClient.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newProduct, err := db.CreateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(newProduct)
}

func dbDeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = db.DeleteProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func dbGetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var limit, offset uint64
	var err error
	if _, ok := params["limit"]; ok {
		limit, err = strconv.ParseUint(params["limit"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if _, ok := params["offset"]; ok {
		offset, err = strconv.ParseUint(params["offset"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	products, err := db.GetProducts(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(products)
}

func dbGetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	product, err := db.GetProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(product)
}

func dbUpdateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product dbClient.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.UpdateProduct(product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	var err error
	if db, err = dbClient.InitClient(); err != nil {
		log.Panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/products", dbCreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", dbDeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", dbGetProducts).Methods("GET")
	r.HandleFunc("/products/{id}", dbGetProduct).Methods("GET")
	r.HandleFunc("/products", dbUpdateProduct).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", r))
}
