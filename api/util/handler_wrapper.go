package util

import (
	"database/sql"
	"net/http"
)

type HandlerLogic func(db *sql.DB, config *Config, w http.ResponseWriter, r *http.Request) error

type Handler struct {
	DB      *sql.DB
	Config  *Config
	Handler HandlerLogic
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.Handler(h.DB, h.Config, w, r); err != nil {
		// handle erorr here
	}
}
