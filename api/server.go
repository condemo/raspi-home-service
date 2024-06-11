package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/condemo/raspi-home-service/api/handlers"
	"github.com/condemo/raspi-home-service/api/middlewares"
	"github.com/condemo/raspi-home-service/store"
)

type ApiServer struct {
	store store.Store
	addr  string
}

func NewAPIServer(addr string, db store.Store) *ApiServer {
	return &ApiServer{addr: addr, store: db}
}

func (s ApiServer) Run() {
	auth := http.NewServeMux()
	router := http.NewServeMux()
	ws := http.NewServeMux()
	view := http.NewServeMux()
	fs := http.FileServer(http.Dir("public/static"))

	basicMiddStack := middlewares.MiddlewareStack(
		middlewares.RequireAuth,
		middlewares.SimpleLogger,
	)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", auth))
	router.Handle("/ws/", http.StripPrefix("/ws", basicMiddStack(ws)))
	router.Handle("/", view)
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers
	userHandler := handlers.NewUserHandler(s.store)
	wsHander := handlers.NewWSHandler(s.store)
	viewHander := handlers.NewViewHanlder(s.store)

	// Routes Load
	userHandler.RegisterRoutes(auth)
	wsHander.RegisterRoutes(ws)
	viewHander.RegisterRoutes(view)

	server := http.Server{
		Addr:         s.addr,
		Handler:      router,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 4 * time.Second,
	}

	// Running the server in a separate goroutine is necessary
	// so that it does not block the execution of Shutdown
	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt)
	signal.Notify(sigC, syscall.SIGTERM)

	<-sigC

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// server.Shutdown ends the execution of the program
	// after waiting for all active connections to finish or 30 seconds to pass
	server.Shutdown(ctx)
}
