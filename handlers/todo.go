package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"todo-api/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var todos []models.Todo

const dataFile = "todos.json"

func authenticate(w http.ResponseWriter, r *http.Request) (string, bool) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Missing token"})
		return "", false
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrNotSupported
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid token"})
		return "", false
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	username := claims["iss"].(string)
	return username, true
}

// JSON'dan verileri yükle
func loadTodos() {
	file, err := os.Open(dataFile)
	if err != nil {
		todos = []models.Todo{} // Dosya yoksa boş başla
		return
	}
	defer file.Close()
	_ = json.NewDecoder(file).Decode(&todos)
}

// Verileri dosyaya yaz
func saveTodos() {
	file, _ := os.Create(dataFile)
	defer file.Close()
	_ = json.NewEncoder(file).Encode(todos)
}

func Init() {
	loadTodos()
}

func GetTodos(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	if search != "" {
		var results []models.Todo
		for _, t := range todos {
			if strings.Contains(strings.ToLower(t.Title), strings.ToLower(search)) {
				results = append(results, t)
			}
		}
		json.NewEncoder(w).Encode(results)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)

	todo.ID = strconv.Itoa(rand.Intn(1000000))
	todo.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	todos = append(todos, todo)
	saveTodos()

	json.NewEncoder(w).Encode(todo)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range todos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo not found"})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedTodo models.Todo
	_ = json.NewDecoder(r.Body).Decode(&updatedTodo)

	for i, t := range todos {
		if t.ID == params["id"] {
			updatedTodo.ID = t.ID // ID değişmesin
			todos[i] = updatedTodo
			saveTodos()
			json.NewEncoder(w).Encode(updatedTodo)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo not found"})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for i, t := range todos {
		if t.ID == params["id"] {
			todos = append(todos[:i], todos[i+1:]...) // slice'tan sil
			saveTodos()
			json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo not found"})
}

func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, t := range todos {
		if t.ID == params["id"] {
			todos[i].Done = !todos[i].Done
			saveTodos()
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo not found"})
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	var total, done, pending int
	total = len(todos)

	for _, t := range todos {
		if t.Done {
			done++
		} else {
			pending++
		}
	}

	stats := map[string]int{
		"total":   total,
		"done":    done,
		"pending": pending,
	}

	json.NewEncoder(w).Encode(stats)
}
