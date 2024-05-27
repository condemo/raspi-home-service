package api

import (
	"net/http"

	"github.com/condemo/raspi-test/api/handlers"
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
	router := http.NewServeMux()
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// Handlers
	userHandler := handlers.NewUserHandler(s.store)
	userHandler.RegisterRoutes(router)

	server := http.Server{
		Addr:    s.addr,
		Handler: v1,
	}

	return server.ListenAndServe()
}
