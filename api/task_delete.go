package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/consul-terraform-sync/api/oapigen"
	"github.com/hashicorp/consul-terraform-sync/logging"
)

// DeleteTaskByName deletes an existing task and its events. Does not delete
// if the task is active.
func (h *TaskLifeCycleHandler) DeleteTaskByName(w http.ResponseWriter, r *http.Request, name string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	ctx := r.Context()
	requestID := requestIDFromContext(ctx)
	logger := logging.FromContext(r.Context()).Named(createTaskSubsystemName).With("task_name", name)
	logger.Trace("delete task request", "delete_task_request", name)

	// Check if task exists
	_, err := h.ctrl.Task(ctx, name)
	if err != nil {
		logger.Trace("task not found", "error", err)
		sendError(w, r, http.StatusNotFound, err)
		return
	}

	// TODO: check if task is active

	err = h.ctrl.TaskDelete(ctx, name)
	if err != nil {
		// TODO error types
		switch err.Error() {
		case "active":
			err := fmt.Errorf("task '%s' is currently running and cannot be deleted "+
				"at this time", name)
			logger.Trace("task active", "error", err)
			sendError(w, r, http.StatusConflict, err)
		default:
			sendError(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := oapigen.TaskResponse{RequestId: requestID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error("error encoding json response", "error", err, "response", resp)
	}
}
