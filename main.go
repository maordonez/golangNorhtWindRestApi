package main

import (
	"GolangNorhtWindRestApi/database"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// TestConnection prueba de endpoint
func TestConnection(w http.ResponseWriter, r *http.Request) {

	con := database.InitDB()
	defer con.Close()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "con: %v\n", con)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/test", TestConnection)
	http.ListenAndServe(":3000", r)
}
