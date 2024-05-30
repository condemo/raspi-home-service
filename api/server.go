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
	view := http.NewServeMux()
	fs := http.FileServer(http.Dir("./public/static"))

	basicMiddlewareStack := middlewares.MiddlewareStack(
		middlewares.RequireAuth,
		middlewares.SimpleLogger,
	)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", auth))
	router.Handle("/api/v1/info/", http.StripPrefix("/api/v1/info",
		basicMiddlewareStack(info)))
	router.Handle("/", view)
	router.Handle("/static/", http.StripPrefix("/static",
		middlewares.ServeStatic(fs)))

	// Handlers
	userHandler := handlers.NewUserHandler(s.store)
	infoHander := handlers.NewInfoHandler(s.store)
	viewHander := handlers.NewViewHanlder(s.store)

	// Routes Load
	userHandler.RegisterRoutes(auth)
	infoHander.RegisterRoutes(info)
	viewHander.RegisterRoutes(view)

	server := http.Server{
		Addr:    s.addr,
		Handler: router,
	}

	return server.ListenAndServe()
}
