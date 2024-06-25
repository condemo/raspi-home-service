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
	return e.Msg
}

func HTTPErrResponse(w http.ResponseWriter, e ApiError, log bool) {
	if log {
		slog.Error("[API]", "msg", e.Error(), "status", strconv.FormatInt(int64(e.Status), 10))
	}
	w.WriteHeader(e.Status)
	json.NewEncoder(w).Encode(e)
}
