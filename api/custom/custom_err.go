package custom

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type ApiError struct {
	Err    error  `json:"-"`
	Msg    string `json:"error"`
	Status int    `json:"code"`
}

func (e *ApiError) Error() string {
	return e.Err.Error()
}

func NewApiError(err error, msg string, status int) ApiError {
	return ApiError{
		Err:    err,
		Msg:    msg,
		Status: status,
	}
}

func HTTPErrResponse(w http.ResponseWriter, e ApiError, log bool) {
	if log {
		slog.Error("[API]", "err", e.Error(), "status", strconv.FormatInt(int64(e.Status), 10))
	}
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}
