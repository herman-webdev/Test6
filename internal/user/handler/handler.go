package handler

import (
	"awesomeProject/internal/apperror"
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/user/dto"
	"awesomeProject/internal/user/service"
	"awesomeProject/internal/user/storage"
	"awesomeProject/pkg/api/sort"
	"awesomeProject/pkg/logging"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	usersURL      = "/users"
	userURL       = "/users/:uuid"
	userUpdateURL = "/users/update/:uuid"
	userDeleteURL = "/users/delete/:uuid"
)

type handler struct {
	logger     *logging.Logger
	service    service.UserService
	repository storage.Repository
}

func NewHandler(repository storage.Repository, logger *logging.Logger) handlers.Handler {
	userService := service.NewUserService(repository, logger)
	return &handler{
		logger:     logger,
		service:    userService,
		repository: repository,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, usersURL, sort.Middleware(apperror.Middleware(h.GetList), "created_at", sort.ASC))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetOne))
	router.HandlerFunc(http.MethodPut, userUpdateURL, apperror.Middleware(h.UpdateOne))
	router.HandlerFunc(http.MethodDelete, userDeleteURL, apperror.Middleware(h.DeleteOne))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var createUserDto dto.CreateUserDto

	if err := json.NewDecoder(r.Body).Decode(&createUserDto); err != nil {
		return apperror.BadRequest("Invalid request body", err.Error())
	}

	if err := h.service.CreateUser(r.Context(), createUserDto); err != nil {
		return apperror.BadRequest("Failed to create user", err.Error())
	}

	w.WriteHeader(http.StatusCreated)
	return nil
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	var sortOptions sort.Options
	if options, ok := r.Context().Value(sort.OptionsContextKey).(sort.Options); ok {
		sortOptions = options
	}

	all, err := h.service.GetAll(r.Context(), sortOptions)
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}

func (h *handler) GetOne(w http.ResponseWriter, r *http.Request) error {
	uuid := httprouter.ParamsFromContext(r.Context()).ByName("uuid")
	if uuid == "" {
		return apperror.BadRequest("id has not been provided", "")
	}

	one, err := h.service.GetOne(r.Context(), uuid)
	if err != nil {
		return apperror.NotFound("user not found", err.Error())
	}

	responseBytes, err := json.Marshal(one)
	if err != nil {
		return apperror.InternalServerError("Failed to marshal user data", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)

	return nil
}

func (h *handler) UpdateOne(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")
	if uuid == "" {
		return apperror.BadRequest("ID has not been provided", "")
	}

	var updateUserDto dto.UpdateUserDto
	if err := json.NewDecoder(r.Body).Decode(&updateUserDto); err != nil {
		return apperror.BadRequest("Invalid request body", err.Error())
	}

	err := h.service.UpdateOne(r.Context(), updateUserDto, uuid)
	if err != nil {
		return apperror.InternalServerError("User not found", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *handler) DeleteOne(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	uuid := params.ByName("uuid")
	if uuid == "" {
		return apperror.BadRequest("ID has not been provided", "")
	}

	err := h.service.DeleteOne(r.Context(), uuid)
	if err != nil {
		return apperror.NotFound("User not found", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
