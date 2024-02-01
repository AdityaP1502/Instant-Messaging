package middleware

import (
	"database/sql"
	"net/http"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
)

type middleware func(next http.Handler, db *sql.DB, config *util.Config) http.Handler

func UseMiddleware(db *sql.DB, config *util.Config, handler http.Handler, middlewares ...middleware) http.Handler {
	chained := handler

	for i := len(middlewares) - 1; i > -1; i-- {
		chained = middlewares[i](chained, db, config)
	}

	return chained
}

func AuthMiddleware(next http.Handler, db *sql.DB, config *util.Config) http.Handler {
	fn := func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error {
		// Do auth shenanigans here
		return nil
	}

	return &util.Handler{
		DB:      db,
		Config:  config,
		Handler: util.HandlerLogic(fn),
	}
}
