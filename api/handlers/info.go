package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/condemo/raspi-home-service/store"
	"github.com/condemo/raspi-home-service/tools"
	"github.com/condemo/raspi-home-service/views/components"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	store   store.Store
	sysInfo *tools.SysInfo
	mu      *sync.RWMutex
	conns   map[*websocket.Conn]struct{}
}

func NewWSHandler(s store.Store) *WSHandler {
	return &WSHandler{
		store:   s,
		sysInfo: tools.NewSysInfo(),
		mu:      new(sync.RWMutex),
		conns:   make(map[*websocket.Conn]struct{}),
	}
}

func (h *WSHandler) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("/info", h.getConn)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (h *WSHandler) getConn(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ErrorLog(w, http.StatusBadRequest, "connection err")
	}

	h.handleWs(conn)
}

func (h *WSHandler) handleWs(c *websocket.Conn) {
	fmt.Println("New Connection:", c.RemoteAddr())

	h.mu.Lock()
	h.conns[c] = struct{}{}
	h.mu.Unlock()

	s := make(chan struct{})

	go h.writeLoop(c, s)
	go h.readLoop(c, s)
}

func (h *WSHandler) writeLoop(c *websocket.Conn, s chan struct{}) {
	t := time.NewTicker(5 * time.Second)
	// h.sysInfo.Update()
	for {
		select {
		case <-t.C:
			h.sysInfo.Update()
			tmpl, err := templ.ToGoHTML(context.Background(), components.InfoBar(h.sysInfo))
			if err != nil {
				fmt.Println("error converting component to html:", err)
			}

			c.WriteMessage(websocket.TextMessage, []byte(tmpl))
			// c.WriteJSON(h.sysInfo)

		case <-s:
			h.mu.Lock()
			delete(h.conns, c)
			h.mu.Unlock()
			fmt.Printf("Connection with %s closed\n", c.RemoteAddr())
			return
		}
	}
}

func (h *WSHandler) readLoop(c *websocket.Conn, s chan struct{}) {
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			close(s)
			break
		}
	}
}
