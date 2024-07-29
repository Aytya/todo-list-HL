package handler

import (
	"encoding/json"
	"github.com/Aytya/todo-list-HL/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sort"
	"time"
)

var todos []domain.TodoModel

func isUnique(title, activeAt string) bool {
	for _, t := range todos {
		if t.Title == title && t.ActiveAt == activeAt {
			return false
		}
	}
	return true
}

func findById(id string) (*domain.TodoModel, int) {
	for i, todo := range todos {
		if todo.ID == id {
			return &todo, i
		}
	}
	return nil, -1
}

func sortTodos(todos []domain.TodoModel) {
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

func addWeekendPrefix(todos []domain.TodoModel) {
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

// AddTodo
// @Summary      add newTodo
// @Tags         todo
// @Accept       json
// @Produce      json
// @Param		 input body domain.TodoModel true "todoModel request"
// @Success      200  {object}  domain.TodoModel
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks [post]
func AddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTodo domain.TodoModel
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

// UpdateTodo
// @Summary      update todo by id
// @Tags         todo
// @Accept       json
// @Produce      json
// @Param 		 id path string true "todo id"
// @Param		 input body domain.TodoModel true "todoModel request"
// @Success      204  "No Content"
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/{id} [put]
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	var updatedTodo domain.TodoModel
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

// DeleteTodo
// @Summary      delete todo by id
// @Tags         todo
// @Param 		 id path string true "todo id"
// @Success      204 "No Content"
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/{id} [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	_, index := findById(id)
	if index == -1 {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	todos = append(todos[:index], todos[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

// Checked
// @Summary      mark todo by id
// @Tags         todo
// @Param 		 id path string true "todo id"
// @Param        input body domain.TodoModel true "Todo Status"
// @Success      204 "No Content"
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/{id}/done [put]
func Checked(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "ID")
	if id == "" {
		http.Error(w, "Invalid or missing ID", http.StatusBadRequest)
		return
	}

	var checkedTodo domain.TodoModel
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

// GetTodo
// @Summary      get todo
// @Tags         todo
// @Produce      json
// @Param        status query string false "Status of the todo" default(active)
// @Success      200  {array}  domain.TodoModel
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks [get]
func GetTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "active"
	}

	var filteredTodos []domain.TodoModel

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
