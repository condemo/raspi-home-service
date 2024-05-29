package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/condemo/raspi-test/api/handlers"
	"github.com/condemo/raspi-test/api/util"
)

type wrapperResponse struct {
	http.ResponseWriter
	status int
}

type Middleware func(next http.Handler) http.HandlerFunc

// func MiddlewareStack(m...Middleware) http.HandlerFunc {}

func RequireAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("Authorization")
		if t == "" {
			handlers.ErrorLog(w, http.StatusBadRequest, "invalid Authorization")
			return
		}

		token := strings.TrimPrefix(t, "Bearer ")
		claims, err := util.ValidateJWT(token)
		if err != nil {
			handlers.ErrorLog(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		id := claims.UserID
		ctx := context.WithValue(r.Context(), util.UserID("userID"), id)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func SimpleLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wr := &wrapperResponse{w, http.StatusOK}
		next.ServeHTTP(wr, r)

		fmt.Printf(
			"[%s] %s [%d] %s - %s", start.Format(time.DateTime), r.Method,
			wr.status, r.URL.Path, time.Since(start),
		)
	}
}
