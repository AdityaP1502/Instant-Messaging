package middleware

import (
	"database/sql"
	"net/http"

	httpx "github.com/AdityaP1502/Instant-Messanging/api/http"
	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
)

func CertMiddleware(db *sql.DB, conf interface{}, next http.Handler) http.Handler {
	fn := func(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error {
		if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
			// TODO: Check the certificate whether it was revoked
			next.ServeHTTP(w, r)
		}

		return responseerror.InvalidTokenErr.Init("Access Denied. Certificate is either empty, invalid, or revoked")
	}

	return &httpx.Handler{
		DB:      db,
		Config:  conf,
		Handler: httpx.HandlerLogic(fn),
	}
}