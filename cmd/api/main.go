package main

import (
	"database/sql"
	"log"
	"net/http"
	"yatdl/internal/config"
	"yatdl/internal/user"
)

func main() {
	db := config.ConnectToDatabase()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	userStore := user.NewStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.Create)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
