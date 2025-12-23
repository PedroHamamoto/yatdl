package main

import (
	"log"
	"net/http"
	"yatdl/internal/config"
	"yatdl/internal/user"
)

func main() {
	db := config.ConnectToDatabase()
	defer db.Close()

	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.Create)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
