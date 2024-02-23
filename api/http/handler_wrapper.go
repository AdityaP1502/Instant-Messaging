package httpx

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaP1502/Instant-Messanging/api/http/responseerror"
	"github.com/AdityaP1502/Instant-Messanging/api/jsonutil"
)

type HandlerLogic func(db *sql.DB, conf interface{}, w http.ResponseWriter, r *http.Request) error

type Handler struct {
	DB      *sql.DB
	Config  interface{}
	Handler HandlerLogic
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tmp struct {
		Status      string `json:"status"`
		ErrorType   string `json:"error_type"`
		Description string `json:"description"`
	}

	w.Header().Set("Content-Type", "application/json")

	if err := h.Handler(h.DB, h.Config, w, r); err != nil {
		if errors.As(err, &responseerror.InternalServiceErr) {
			fmt.Println(err.(*responseerror.InternalServiceError).Description)
		}

		if errGetter, ok := err.(responseerror.ResponseErrorGetter); !ok {
			tmp.Status = "fail"
			tmp.ErrorType = "internal_service_error"
			tmp.Description = "Something is wrong. Please try again!"

			w.WriteHeader(500)
		} else {
			requestErr := errGetter.Get()

			tmp.Status = "fail"
			tmp.ErrorType = requestErr.Name
			tmp.Description = requestErr.Message

			w.WriteHeader(requestErr.Code)
		}

		jsonResponse, err := jsonutil.EncodeToJson(&tmp)

		if err != nil {
			http.Error(w, "Something wrong with server!", 500)
		}

		w.Write(jsonResponse)
	}
}
