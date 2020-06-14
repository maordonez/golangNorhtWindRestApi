package main

import (
	"GolangNorhtWindRestApi/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var databaseConnection *sql.DB

// Product descripcion
type Product struct {
	ID          int    `json:"id"`
	ProductCode string `json:"productCode"`
	Description string `json:"description"`
}

// TestConnection prueba de endpoint
func testConnection(w http.ResponseWriter, r *http.Request) {

	con := database.InitDB()
	defer con.Close()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "con: %v\n", con)
}

// AllProducts lista todos los productos de la bd
func AllProducts(w http.ResponseWriter, r *http.Request) {
	const sql = `SELECT id, product_code, COALESCE(description, '')
               FROM products`
	results, err := databaseConnection.Query(sql)
	catch(err)

	var products []*Product
	for results.Next() {
		product := &Product{}
		err = results.Scan(&product.ID, &product.ProductCode, &product.Description)

		catch(err)
		products = append(products, product)
	}
	responseWithJSON(w, r, http.StatusOK, products)
}

// UpdateProduct actualizar un producto
func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	vars := mux.Vars(r)
	id := vars["productId"]
	query, err := databaseConnection.Prepare(`UPDATE product SET product_code = ?, description = ?, where id = ?`)
	catch(err)
	_, er := query.Exec(product.ProductCode, product.Description, id)
	catch(er)
	defer query.Close()

	responseWithJSON(w, r, http.StatusOK, map[string]string{"message": "update successfully"})
}

// CreateProduct crear un nuevo registro en  la bd
func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	query, err := databaseConnection.Prepare("INSERT INTO products(product_code, description) values (?, ?)")
	catch(err)

	_, er := query.Exec(product.ProductCode, product.Description)
	catch(er)
	defer query.Close()

	responseWithJSON(w, r, http.StatusCreated, map[string]string{"message": "successfully created"})

}

// DeleteProduct eliminar un producto
func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["productId"]

	query, err := databaseConnection.Prepare("DELETE FROM product WHERE id = ?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()
	responseWithJSON(w, r, http.StatusOK, map[string]string{"message": "successfully"})
}

func responseWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	databaseConnection = database.InitDB()
	defer databaseConnection.Close()
	r := mux.NewRouter()
	r.HandleFunc("/test", testConnection).Methods("GET")
	r.HandleFunc("/products", AllProducts).Methods("GET")
	r.HandleFunc("/products", CreateProduct).Methods("POST")
	r.HandleFunc("/products/{productId}", UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{productId}", DeleteProduct).Methods("DELETE")
	http.ListenAndServe(":3000", r)
}
