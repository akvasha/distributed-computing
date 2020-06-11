package main

import (
	"DC-homework-1/skeleton/authClient"
	"DC-homework-1/skeleton/dbClient"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var db dbClient.Client
var ac authClient.AuthClient

type errorResponse struct {
	Error string `json:"error"`
}

func responseError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintln(err)})
}

func dbCreateProduct(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var product dbClient.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	newProduct, err := db.CreateProduct(product)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(newProduct)
}

func dbDeleteProduct(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if err = db.DeleteProduct(id); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
}

func dbGetProducts(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var limit, offset uint64
	var err error
	if _, ok := params["limit"]; ok {
		limit, err = strconv.ParseUint(params["limit"], 10, 64)
		if err != nil {
			responseError(w, err, http.StatusBadRequest)
			return
		}
	}
	if _, ok := params["offset"]; ok {
		offset, err = strconv.ParseUint(params["offset"], 10, 64)
		if err != nil {
			responseError(w, err, http.StatusBadRequest)
			return
		}
	}
	products, err := db.GetProducts(limit, offset)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(products)
}

func dbGetProduct(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	product, err := db.GetProduct(id)
	if err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(product)
}

func dbUpdateProduct(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var product dbClient.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if err := db.UpdateProduct(product); err != nil {
		responseError(w, err, http.StatusInternalServerError)
		return
	}
}

func authorize(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("auth")
	if err := ac.Validate(token); err != nil {
		fmt.Println(err.Error())
		if errResp, ok := err.(*authClient.ErrorRespStatus); ok {
			responseError(w, errResp, errResp.StatusCode)
		} else {
			responseError(w, err, http.StatusInternalServerError)
		}
		return false
	}
	return true
}

func main() {
	var err error
	if db, err = dbClient.InitClient(); err != nil {
		log.Panic(err)
	}
	if ac, err = authClient.InitAuthClient(); err != nil {
		log.Panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/products", dbCreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", dbDeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", dbGetProducts).Methods("GET")
	r.HandleFunc("/products/{id}", dbGetProduct).Methods("GET")
	r.HandleFunc("/products", dbUpdateProduct).Methods("PUT")
	log.Println("Skeleton server started")
	log.Fatal(http.ListenAndServe(":8000", r))
}