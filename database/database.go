package database

import (
	"database/sql"
	// Register drive mysql
	_ "github.com/go-sql-driver/mysql"
)

// InitDB funcion devuelve un objecto DB
func InitDB() *sql.DB {
	connectionString := "root:12345678@tcp(localhost:3306)/northwind"
	databaseConnection, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err.Error) // Error Handling
	}
	return databaseConnection
}
