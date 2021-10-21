//go:generate go run oapi-codegen --package=api --generate types,chi-server,spec -o api.gen.go ../openapi.yaml

package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type TaskHandler struct {
	tasks map[string]NewTask
	lock  sync.Mutex
}

// Make sure we conform to ServerInterface
var _ ServerInterface = (*TaskHandler)(nil)

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		tasks: make(map[string]NewTask),
	}
}

// GetTasks returns all tasks
func (t *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	t.lock.Lock()
	defer t.lock.Unlock()

	tasks := make([]NewTask, 0, len(t.tasks))
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}

	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskById returns a task by name
func (t *TaskHandler) GetTaskByName(w http.ResponseWriter, r *http.Request, name string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	tasks := make([]NewTask, 0, len(t.tasks))
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t.tasks[name])
}

// CreateTask creates a task
func (t *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request, params CreateTaskParams) {
	// We expect a NewPet object in the request body.
	var task NewTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendTaskError(w, http.StatusBadRequest, "Invalid format for NewPet")
		return
	}

	t.tasks[task.Name] = task

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// DeleteTaskByName Deletes a task by name
func (t *TaskHandler) DeleteTaskByName(w http.ResponseWriter, r *http.Request, name string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// delete the task
	delete(t.tasks, name)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendTaskError(w http.ResponseWriter, code int, message string) {
	taskErr := Error{
		Code:    int32(code),
		Message: message,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(taskErr)
}
