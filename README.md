# Go - Service for compiling task lists

### RESTful API for the Todo List microservice.

## Getting Started

### Prerequisites:
    - Go 1.22 or later
    - Docker
### Installation:
1. Clone the repository:
   ```bash
   https://github.com/Aytya/todo-list-HL
   ```
2. Navigate into the project directory:
   ```bash
    cd todo-list-HL
   ```
3. Install dependencies:
   ```bash
    go get -u "github.com/swaggo/http-swagger"
    go get -u "github.com/go-chi/chi/v5/middleware"
    go get -u "github.com/go-chi/chi/v5"
    go get -u "github.com/go-chi/chi/render"
   ```

##  Build and Run Locally:
### Build the application:
   ```bash
   make build
   ```
### Run the application:
   ```bash
   make run
   ```
### Stop the application:
   ```bash
   make down
   ```

### Generate Swagger Documentation:
   ```bash
   swag init -g cmd/main.go
   ```
## API Endpoints:
### Create a New Task:
   - URL: http://localhost:8080/api/todo-list/tasks
   - Method: POST
   - Request Body: 
 ```bash
    {
       "title": "Купить квартиру",
       "activeAt": "2023-02-05"
    }
 ```
### Update an Existing Task:
   - URL: http://localhost:8080/api/todo-list/tasks/{id}
   - Method: PUT
   - Request Body:
 ```bash
    { 
        "title": "Купить квартиру - завтра",
        "activeAt": "2023-08-05" 
    }
 ```
### Delete an Existing Task:
- URL: http://localhost:8080/api/todo-list/tasks/{id}
- Method: DELETE

### Change Task Status From Active To Done:
- URL: http://localhost:8080/api/todo-list/tasks/{id}/done
- Method: PUT 

### Change Task Status From Active To Done:
- URL: http://localhost:8080/api/todo-list/tasks
- Method: GET

### Swagger Documentation
- URL: /swagger/
