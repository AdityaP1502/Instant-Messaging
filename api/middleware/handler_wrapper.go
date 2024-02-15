package middleware

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/AdityaP1502/Instant-Messaging/api/api/model"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	requesterror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error"
	internalserviceerror "github.com/AdityaP1502/Instant-Messaging/api/api/util/request_error/internal_service_error"
)

type HandlerLogic func(db *sql.DB, config *util.Config, w http.ResponseWriter, r *http.Request) error

type Handler struct {
	DB      *sql.DB
	Config  *util.Config
	Handler HandlerLogic
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := h.Handler(h.DB, h.Config, w, r); err != nil {
		var resp *model.ErrorResponse

		if errors.As(err, &internalserviceerror.InternalServiceErr) {
			fmt.Println(err.(*internalserviceerror.InternalServiceError).Description)
		}

		if errGetter, ok := err.(requesterror.CustomErrorGetter); !ok {
			resp = &model.ErrorResponse{
				Status:      "fail",
				ErrorType:   "internal_service_error",
				Description: "Something is wrong. Please try again!",
			}
			w.WriteHeader(500)
		} else {
			requestErr := errGetter.Get()
			resp = &model.ErrorResponse{
				Status:      "fail",
				ErrorType:   requestErr.Name,
				Description: requestErr.Message,
			}

			w.WriteHeader(requestErr.Code)
		}

		jsonResponse, err := util.CreateJSONResponse(resp)

		if err != nil {
			http.Error(w, "Something wrong with server!", 500)
		}

		w.Write(jsonResponse)
	}
}
