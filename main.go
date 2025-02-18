package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

var products []Product

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		fmt.Fprintf(w, "Product:", product)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)
		products = append(products, product)
	}
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintln(w, "All products fetched successfully")
		fmt.Fprintln(w, products)
	}
	w.WriteHeader(http.StatusOK)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updatedProduct Product
	err := json.NewDecoder(r.Body).Decode(&updatedProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		found := false
		for index, singleProduct := range products {
			if singleProduct.Id == id {
				found = true
				products[index].Name = updatedProduct.Name
				products[index].Description = updatedProduct.Description
				products[index].Price = updatedProduct.Price
			}
		}
		if !found {
			http.NotFound(w, r)
			return
		}

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product with ID %d updated successfully", id)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	found := false
	var itemToDelete Product
	err := json.NewDecoder(r.Body).Decode(&itemToDelete)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		for index, singleProduct := range products {
			if singleProduct.Id == id {
				found = true
				products = append(products[:index], products[index+1:]...)
				break
			}
		}
	}
	if !found {
		http.NotFound(w, r)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
		w.WriteHeader(http.StatusOK)
	}
}

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func handleError(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "An error occurred",
		Status:  http.StatusInternalServerError,
	}
	json.NewEncoder(w).Encode(&response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/products", createProduct).Methods("POST")
	mux.HandleFunc("/products", getProduct).Methods("GET")
	mux.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	mux.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")
	mux.HandleFunc("/error", handleError).Methods("GET")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux.StrictSlash(true),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	} else {
		fmt.Println("Server started successfully on port 8080")
	}
}
