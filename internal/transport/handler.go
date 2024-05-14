package transport

import (
	"ScriptService/internal/model"
	"ScriptService/internal/service"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Handler struct {
	services *service.Services
}

var ContextMap sync.Map

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

	ctx := context.Background()
	c, cancel := context.WithCancel(ctx)
	contextCommand := model.ContextCommand{Ctx: c, Cancel: cancel}

	ContextMap.Store(command.ID, contextCommand)

	go h.services.Command.ExecuteCommand(contextCommand, command)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(command.ID)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(commands)
}

func (h *Handler) GetCommandById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	commandID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	command, err := h.services.Command.GetCommand(commandID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(command)
}

func (h *Handler) CancelCommand(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	commandID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.services.Command.CancelCommand(&ContextMap, commandID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
