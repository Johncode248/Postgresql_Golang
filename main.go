package main

import (
	"encoding/json"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println(uuid.New())

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/products", createProduct).Methods("POST")
	r.HandleFunc("/api/v1/products", getProducts).Methods("GET")
	r.HandleFunc("/api/v1/products{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/v1/products{id}", updateProduct).Methods("PATCH")
	r.HandleFunc("/api/v1/products{id}", deleteProduct).Methods("DELETE")

	db := connect()
	defer db.Close()

	//product := Product{ID: uuid.New(), Name: "My Product 1", Quantity: 10, Price: 43.99}
	///fmt.Println(product)

	log.Fatal(http.ListenAndServe(":8085", r))
}

// Create product
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	// Creating Product instance
	product := &Product{
		//ID: rand(34),
	}

	// Decoding request
	_ = json.NewDecoder(r.Body).Decode(&product)

	//Inserting into database
	_, err := db.Model(product).Insert()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	// Returning product
	json.NewEncoder(w).Encode(product)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	// Creating Products Slice

	var products []Product
	if err := db.Model(&products).Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning products
	json.NewEncoder(w).Encode(products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	// Get ID
	params := mux.Vars(r)
	productId := params["id"]

	// Creating product instance
	product := &Product{ID: productId}
	if err := db.Model(product).WherePK().Select(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning Product
	json.NewEncoder(w).Encode(product)

}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	// Get ID
	params := mux.Vars(r)
	productId := params["id"]

	// Creating product instance
	product := &Product{ID: productId}

	_ = json.NewDecoder(r.Body).Decode(&product)

	//db.Model(product).WherePK().Set("name = ?, quantity = ?, price = ?, store = ?", product.Name, product.Quantity, product.Price, product.Store).Update()
	_, err := db.Model(product).WherePK().Set("name = ?, quantity = ?, price = ?, store = ?", product.Name, product.Quantity, product.Price, product.Store).Update()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Returning  product
	json.NewEncoder(w).Encode(product)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get connect
	db := connect()
	defer db.Close()

	// Get ID
	params := mux.Vars(r)
	productId := params["id"]

	product := &Product{}
	result, err := db.Model(product).Where("id = ?", productId).Delete()

	//Other way to create product
	//product := &Product(ID: productId)
	//	result, err := db.Model(product).WherePK().Delete())

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Returning json result
	json.NewEncoder(w).Encode(result)

}
