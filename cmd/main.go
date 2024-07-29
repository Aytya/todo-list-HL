package main

import (
	_ "github.com/Aytya/todo-list-HL/docs"
	"github.com/Aytya/todo-list-HL/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title           Service for compiling task lists
// @version         1.0
// @description     RESTful API for the TodoList microservice.

// @host      localhost:8080
// @BasePath  /api/todo-list
func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/api/todo-list", func(r chi.Router) {
		r.Post("/tasks", handler.AddTodo)
		r.Put("/tasks/{ID}", handler.UpdateTodo)
		r.Get("/tasks", handler.GetTodo)
		r.Delete("/tasks/{ID}", handler.DeleteTodo)
		r.Put("/tasks/{ID}/done", handler.Checked)
	})

	log.Println("Starting HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
