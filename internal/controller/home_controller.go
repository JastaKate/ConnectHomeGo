package controller

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"connectHome/internal/request"
	"connectHome/internal/service"
)

type HomeController struct {
	service *service.HomeService
}

func NewHomeController(homeService *service.HomeService) *HomeController {
	return &HomeController{service: homeService}
}

func (h *HomeController) Create(w http.ResponseWriter, r *http.Request) {
	var input request.HomeRequest

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		writeError(w, http.StatusBadRequest, "Ожидается один JSON-объект")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	home, err := h.service.Create(ctx, input)
	if errors.Is(err, service.ErrInvalidHome) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		log.Printf("Create home: %v", err)
		writeError(w, http.StatusInternalServerError, "Не удалось создать дом")
		return
	}
	writeJSON(w, http.StatusCreated, home)
}

func (h *HomeController) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	homes, err := h.service.List(ctx)
	if err != nil {
		log.Printf("List homes: %v", err)
		writeError(w, http.StatusInternalServerError, "Не удалось получить дома")
		return
	}
	writeJSON(w, http.StatusOK, homes)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("Write JSON response: %v", err)
	}
}
