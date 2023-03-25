package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rotisserie/eris"
	sqlc "github.com/swuecho/chatgpt_backend/sqlc_queries"
)

type UserActiveChatSessionHandler struct {
	service *UserActiveChatSessionService
}

func NewUserActiveChatSessionHandler(service *UserActiveChatSessionService) *UserActiveChatSessionHandler {
	return &UserActiveChatSessionHandler{
		service: service,
	}
}

func (h *UserActiveChatSessionHandler) Register(router *mux.Router) {
	router.HandleFunc("/uuid/user_active_chat_session", h.GetUserActiveChatSessionHandler).Methods(http.MethodGet)
	router.HandleFunc("/uuid/user_active_chat_session", h.CreateUserActiveChatSessionHandler).Methods(http.MethodPost)
	router.HandleFunc("/uuid/user_active_chat_session", h.CreateOrUpdateUserActiveChatSessionHandler).Methods(http.MethodPut)
}

// CreateUserActiveChatSessionHandler handles POST requests to create a new session.
func (h *UserActiveChatSessionHandler) CreateUserActiveChatSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := ctx.Value(userContextKey).(string)
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Error: '"+userIDStr+"' is not a valid user ID. Please enter a valid user ID.", http.StatusBadRequest)
		return
	}

	var sessionParams sqlc.CreateUserActiveChatSessionParams
	if err := json.NewDecoder(r.Body).Decode(&sessionParams); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	sessionParams.UserID = int32(userIDInt)

	session, err := h.service.CreateUserActiveChatSession(r.Context(), sessionParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// CreateOrUpdateUserActiveChatSessionHandler handles POST requests to create a new session.
func (h *UserActiveChatSessionHandler) CreateOrUpdateUserActiveChatSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := ctx.Value(userContextKey).(string)
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Error: '"+userIDStr+"' is not a valid user ID. Please enter a valid user ID.", http.StatusBadRequest)
		return
	}

	var sessionParams sqlc.CreateOrUpdateUserActiveChatSessionParams
	if err := json.NewDecoder(r.Body).Decode(&sessionParams); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// use the user_id from token
	sessionParams.UserID = int32(userIDInt)
	session, err := h.service.CreateOrUpdateUserActiveChatSession(r.Context(), sessionParams)
	if err != nil {
		http.Error(w, eris.Wrap(err, "fail to update or create action user session record, ").Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// GetUserActiveChatSessionHandler handles GET requests to get a session by user_id.
func (h *UserActiveChatSessionHandler) GetUserActiveChatSessionHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDStr := ctx.Value(userContextKey).(string)
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Error: '"+userIDStr+"' is not a valid user ID. Please enter a valid user ID.", http.StatusBadRequest)
		return
	}

	session, err := h.service.GetUserActiveChatSession(r.Context(), int32(userIDInt))
	if err != nil {
		http.Error(w, fmt.Errorf("error: %v", err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// UpdateUserActiveSessionHandler handles PUT requests to update an existing session.
func (h *UserActiveChatSessionHandler) UpdateUserActiveSessionHandler(w http.ResponseWriter, r *http.Request) {
	// Get the path variables using gorilla/mux
	ctx := r.Context()
	userIDStr := ctx.Value(userContextKey).(string)
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Error: '"+userIDStr+"' is not a valid user ID. Please enter a valid user ID.", http.StatusBadRequest)
		return
	}
	// Get the query parameters using r.URL.Query() method
	queryParams := r.URL.Query()
	chatSessionUUID := queryParams.Get("chatSessionUuid")

	session, err := h.service.UpdateUserActiveChatSession(r.Context(), int32(userIDInt), chatSessionUUID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}
