package main

import (
	"database/sql"
	"encoding/json"
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

	jwt := auth.NewJwt(os.Getenv("JWT_SECRET"))
	authService := auth.NewService(jwt, userStore)
	authHandler := auth.NewHandler(authService)

	mux := http.NewServeMux()
	// Public
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("POST /login", authHandler.Login)

	// Protected
	authMiddleware := middleware.Auth(jwt)
	protectedGreeting := authMiddleware(
		middleware.JSONContentType(http.HandlerFunc(Greeting)),
	)

	mux.Handle("GET /greeting", protectedGreeting)

	handler := middleware.JSONContentType(mux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}

func Greeting(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode("Hello World!")
	if err != nil {
		return
	}
}
