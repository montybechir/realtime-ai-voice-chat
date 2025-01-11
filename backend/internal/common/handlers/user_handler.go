package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"interviews-ai/internal/services"
	"net/http"
	"strings"
)

type UserHandler struct {
	Service *services.UserService
}
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserRequest struct {
	Id string `json:"id"`
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (r CreateUserRequest) Validate(ctx context.Context) map[string]string {
	errors := make(map[string]string)
	if r.Email == "" {
		errors["email"] = "email is required"
	}
	if r.Name == "" {
		errors["name"] = "name is required"
	}
	return errors
}

func (r GetUserRequest) Validate(ctx context.Context) map[string]string {
	errors := make(map[string]string)
	if r.Id == "" {
		errors["id"] = "id is required"
	}
	return errors
}

func (h *UserHandler) ListUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users, err := h.Service.GetAllUsers(r.Context())
		if err != nil {
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	})
}

func (h *UserHandler) GetUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		id := path[strings.LastIndex(path, "/")+1:]

		if id == "" {
			encodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "user id is required",
			})
			return
		}
		user, err := h.Service.GetUser(r.Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrNotFound):
				encodeJSON(w, http.StatusNotFound, map[string]string{
					"error": "user not found",
				})
			case errors.Is(err, services.ErrInvalidInput):
				encodeJSON(w, http.StatusBadRequest, map[string]string{
					"error": "invalid user id",
				})
			default:
				encodeJSON(w, http.StatusInternalServerError, map[string]string{
					"error": "internal server error",
				})
			}
			return
		}

		encodeJSON(w, http.StatusOK, user)
	})
}

func (h *UserHandler) CreateUser() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, validationErrors, err := decodeAndValidate[CreateUserRequest](r)

		if err != nil {
			if validationErrors != nil {
				encodeJSON(w, http.StatusUnprocessableEntity, validationErrors)
				return
			}
			encodeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		user, err := h.Service.CreateUser(r.Context(), req.Name, req.Email)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		encodeJSON(w, http.StatusCreated, user)
	})
}
