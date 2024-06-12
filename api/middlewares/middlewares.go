package middlewares

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/condemo/raspi-home-service/api/util"
)

type wrapperResponse struct {
	http.ResponseWriter
	status int
}

func (w *wrapperResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()
}

type Middleware func(next http.Handler) http.HandlerFunc

func MiddlewareStack(m ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}
		return next.ServeHTTP
	}
}

func RequireAuth(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusPermanentRedirect)
			return
		}

		claims, err := util.ValidateJWT(c.Value)
		if err != nil {
			http.Redirect(w, r, "/auth/login", http.StatusUnauthorized)
			return
		}

		id := claims.UserID
		ctx := context.WithValue(r.Context(), util.UserID("userID"), id)

		next.ServeHTTP(w, r.Clone(ctx))
	}
}

func SimpleLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wr := &wrapperResponse{w, http.StatusOK}
		next.ServeHTTP(wr, r)

		fmt.Printf(
			"[%s] %s [%d] %s - %s\n", start.Format(time.DateTime), r.Method,
			wr.status, r.URL.Path, time.Since(start),
		)
	}
}
