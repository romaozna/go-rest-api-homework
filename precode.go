package main

import (
	"bytes"
	"encoding/json"
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

func getTasks(writer http.ResponseWriter, request *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(resp)
}

func postTask(writer http.ResponseWriter, request *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

func getTask(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(writer, "Задача не найдена", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(resp)
}

func deleteTask(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(writer, "Задача не найдена", http.StatusNoContent)
		return
	}

	delete(tasks, task.ID)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
