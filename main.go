package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/thedevsaddam/renderer"
	"log"
	"net/http"
	"sort"
	"time"
)

var rnd *renderer.Render

type (
	TodoModel struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		ActiveAt string `json:"activeAt"`
		Status   string `json:"status"`
	}
)

func init() {
	fmt.Println("init function running")

	rnd = renderer.New()
}

var todos []TodoModel

func isUnique(title, activeAt string) bool {
	for _, t := range todos {
		if t.Title == title && t.ActiveAt == activeAt {
			return false
		}
	}
	return true
}

func findById(id string) (*TodoModel, int) {
	for i, todo := range todos {
		if todo.ID == id {
			return &todo, i
		}
	}
	return nil, -1
}

func sortTodos(todos []TodoModel) {
	sort.Slice(todos, func(i, j int) bool {
		activeDateI, err := time.Parse("2006-01-02", todos[i].ActiveAt)
		if err != nil {
			log.Fatal(err)
		}
		activeDateJ, err := time.Parse("2006-01-02", todos[j].ActiveAt)
		if err != nil {
			log.Fatal(err)
		}
		return activeDateI.Before(activeDateJ)
	})
}

func addWeekendPrefix(todos []TodoModel) {
	for i, todo := range todos {
		activeDate, err := time.Parse("2006-01-02", todo.ActiveAt)
		if err != nil {
			log.Fatal(err)
		}

		if activeDate.Weekday() == time.Saturday || activeDate.Weekday() == time.Sunday {
			todos[i].Title = "ВЫХОДНОЙ - " + todo.Title
		}
	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTodo TodoModel
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if len(newTodo.Title) == 0 || len(newTodo.Title) > 200 {
		http.Error(w, "Title length must be between 1 and 200", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("2006-01-02", newTodo.ActiveAt)
	if err != nil {
		http.Error(w, "Active at must be a valid date", http.StatusBadRequest)
	}

	if !isUnique(newTodo.Title, newTodo.ActiveAt) {
		http.Error(w, "Title must be unique", http.StatusBadRequest)
		return
	}

	newTodo.ID = uuid.New().String()
	newTodo.Status = "active"
	todos = append(todos, newTodo)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["ID"]

	var updatedTodo TodoModel
	err := json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(updatedTodo.Title) == 0 || len(updatedTodo.Title) > 200 {
		http.Error(w, "Title length must be between 1 and 200", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("2006-01-02", updatedTodo.ActiveAt)
	if err != nil {
		http.Error(w, "Active at must be a valid date", http.StatusBadRequest)
	}

	_, found := findById(id)
	if found == -1 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos[found].Title = updatedTodo.Title
	todos[found].ActiveAt = updatedTodo.ActiveAt

	w.WriteHeader(http.StatusNoContent)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["ID"]

	_, index := findById(id)
	if index == -1 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos = append(todos[:index], todos[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

func checked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["ID"]

	var checkedTodo TodoModel
	err := json.NewDecoder(r.Body).Decode(&checkedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, found := findById(id)
	if todo == nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos[found].Status = "done"
	w.WriteHeader(http.StatusNoContent)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	var filteredTodos []TodoModel

	now := time.Now()
	for _, todo := range todos {
		activeDate, err := time.Parse("2006-01-02", todo.ActiveAt)
		if err != nil {
			http.Error(w, "Active at must be a valid date", http.StatusBadRequest)
			return
		}

		switch status {
		case "active":
			if todo.Status != "done" && activeDate.Before(now) || activeDate.Equal(now) {
				filteredTodos = append(filteredTodos, todo)
			}
		case "done":
			if todo.Status == "done" {
				filteredTodos = append(filteredTodos, todo)
			}
		default:
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

	}

	sortTodos(filteredTodos)

	addWeekendPrefix(filteredTodos)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredTodos)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/todo-list/tasks", addTodo).Methods("POST")
	router.HandleFunc("/api/todo-list/tasks/{ID}", updateTodo).Methods("PUT")
	router.HandleFunc("/api/todo-list/tasks/{ID}", deleteTodo).Methods("DELETE")
	router.HandleFunc("/api/todo-list/tasks/{ID}/done", checked).Methods("PUT")
	router.HandleFunc("/api/todo-list/tasks/{ID}", getTodo).Methods("GET")

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Listening on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
