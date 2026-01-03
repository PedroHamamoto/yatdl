package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"yatdl/internal/auth"
	"yatdl/internal/config"
	"yatdl/internal/http/middleware"
	"yatdl/internal/task"
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

	taskStore := task.NewStore(db)
	taskService := task.NewService(taskStore)
	taskHandler := task.NewHandler(taskService)

	mux := http.NewServeMux()
	// Public
	mux.HandleFunc("POST /api/users", userHandler.Create)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	// Protected
	authMiddleware := middleware.Auth(jwt)

	mux.Handle("POST /api/tasks", authMiddleware(http.HandlerFunc(taskHandler.Create)))
	mux.Handle("PATCH /api/tasks/{id}", authMiddleware(http.HandlerFunc(taskHandler.Update)))

	handler := middleware.CORS(middleware.JSONContentType(mux))

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
