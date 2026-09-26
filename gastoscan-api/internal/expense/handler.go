package expense

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	expenses, err := handler.service.List()
	if err != nil {
		http.Error(w, "Error al listar los gastos", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, expenses)
}

func (handler *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Id inválido para el gasto", http.StatusBadRequest)
		return
	}

	expense, err := handler.service.Get(uint(id))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Gasto no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Error al obtener el gasto", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, expense)
}

func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var expense Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	expense, err := handler.service.Create(expense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, expense)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
