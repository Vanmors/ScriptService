package transport

import (
	"ScriptService/internal/model"
	"ScriptService/internal/service"
	"encoding/json"
	"net/http"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{services: services}
}

func (h *Handler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var command model.Command
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	command, err := h.services.Command.CreateCommand(command)

	h.services.Command.ExecuteCommandAsync(command)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(command)
}

func (h *Handler) GetAllCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	commands, err := h.services.Command.GetAllCommands()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(commands)
}
