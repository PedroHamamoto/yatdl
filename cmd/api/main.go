package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"yatdl/internal/auth"
	"yatdl/internal/config"
	"yatdl/internal/http/middleware"
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

	authService := auth.NewService(os.Getenv("JWT_SECRET"), userStore)
	authHandler := auth.NewHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("POST /login", authHandler.Login)

	handler := middleware.JSONContentType(mux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
