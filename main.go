package main

import (
	"database/sql"
	"log"
	"net/http"

	"go-exercise/controller"
	"go-exercise/model"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	model.SetupDatabase(db)

	ltpController := &controller.LTPController{DB: db}
	http.HandleFunc("/api/v1/ltp", ltpController.HandleLTP)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}