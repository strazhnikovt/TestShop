package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/service"
	"github.com/strazhnikovt/TestShop/pkg/auth"
)

type UserHandler struct {
	service    *service.UserService
	jwtManager *auth.JWTManager
}

func NewUserHandler(s *service.UserService, jm *auth.JWTManager) *UserHandler {
	return &UserHandler{
		service:    s,
		jwtManager: jm,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		FullName:  strings.TrimSpace(req.FirstName + " " + req.LastName),
		Login:     req.Login,
		Age:       req.Age,
		IsMarried: req.IsMarried,
		Password:  req.Password, // hashing in service.Register
		Role:      "user",       // default
	}

	if err := h.service.Register(user); err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	resp := map[string]int{"id": user.ID}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(req.Login, req.Password)
	if err != nil {
		http.Error(w, `{"error":"Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		http.Error(w, `{"error":"Could not generate token"}`, http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"token": token}
	json.NewEncoder(w).Encode(resp)
}
