package router

import (
	"encoding/json"
	"net/http"

	"gastoscan-api/internal/expense"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

func New(db *gorm.DB) http.Handler {
	r := chi.NewRouter()

	expenseRepository := expense.NewRepository(db)
	expenseService := expense.NewService(expenseRepository)
	expenseHandler := expense.NewHandler(expenseService)

	r.Get("/expenses", expenseHandler.List)
	r.Get("/expenses/{id}", expenseHandler.Get)
	r.Post("/expenses", expenseHandler.Create)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	return r
}
