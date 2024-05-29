package api

import (
	"net/http"

	"github.com/condemo/raspi-test/api/handlers"
	"github.com/condemo/raspi-test/api/middlewares"
	"github.com/condemo/raspi-test/store"
)

type ApiServer struct {
	store store.Store
	addr  string
}

func NewAPIServer(addr string, db store.Store) *ApiServer {
	return &ApiServer{addr: addr, store: db}
}

func (s ApiServer) Run() error {
	auth := http.NewServeMux()
	router := http.NewServeMux()
	info := http.NewServeMux()

	basicMiddlewareStack := middlewares.MiddlewareStack(
		middlewares.RequireAuth,
		middlewares.SimpleLogger,
	)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", auth))
	router.Handle("/api/v1/info/", http.StripPrefix("/api/v1/info",
		basicMiddlewareStack(info)),
	)

	// Handlers
	userHandler := handlers.NewUserHandler(s.store)
	infoHander := handlers.NewInfoHandler(s.store)

	// Routes Load
	userHandler.RegisterRoutes(auth)
	infoHander.RegisterRoutes(info)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	return server.ListenAndServe()
}
