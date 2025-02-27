package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
func getTasks(w http.ResponseWriter, r *http.Request) {
    resp, err := json.MarshalIndent(tasks, "", "   ")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.Header().Set("Content-Type", "application/json")
    	w.WriteHeader(http.StatusOK)
    	w.Write(resp)
}	

func postTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if _, ok := tasks[task.ID]; ok {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.WriteHeader(http.StatusCreated)
}
func getTaskID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskID := chi.URLParam(r, "id")
	task, ok := tasks[taskID]
	if !ok {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "400 Bad Request", http.StatusInternalServerError)
		return
	}
}
func deleteTask(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    _, ok := tasks[id]
    if !ok {
        http.Error(w, "Задача не найдена", http.StatusBadRequest)
        return
    }
    
	delete(tasks, id)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTask)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskID)
	r.Delete("/tasks/{id}", deleteTask)


	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
