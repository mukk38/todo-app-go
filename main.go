package main

import (
	"log"
	"net/http"
	"todo-api/handlers"

	"github.com/gorilla/mux"
)

func main() {
	handlers.Init() // Uygulama başlarken JSON'dan oku

	r := mux.NewRouter()

	r.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	r.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", handlers.GetTodo).Methods("GET")
	r.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")
	r.HandleFunc("/todos/{id}/toggle", handlers.ToggleTodo).Methods("PATCH")
	r.HandleFunc("/stats", handlers.GetStats).Methods("GET")

	log.Println("Sunucu çalışıyor: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
