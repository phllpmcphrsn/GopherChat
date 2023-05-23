package server

import (
	"io"
	"sync"

	log "golang.org/x/exp/slog"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns 		map[*websocket.Conn]bool
	rwmutex 	*sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		rwmutex: &sync.RWMutex{},
	}
}

func (s *Server) HandleWebSocket(ws *websocket.Conn) {
	log := log.With("client", ws.RemoteAddr())
	log.Info("New incoming connection from client")
	
	// for thread safety
	s.rwmutex.Lock()
	s.conns[ws] = true
	s.rwmutex.Unlock()
	
	// Trying to send a welcome message, but in the way I'm testing
	// it may be difficult since the client isn't listening when the 
	// message is sent. Would probably have to create a phony client
	// that sends a message after some time
	// s.Broadcast([]byte("Welcome to Gopher Chat"))
	s.readLoop(ws)
}

// ReadLoop will perform an infinite loop whilst reading from a client, and 
// will broadcast the bytes (messages) to all clients
func (s *Server) readLoop(ws *websocket.Conn) {
	log := log.With("client", ws.RemoteAddr())

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			// client's connection terminated
			if err == io.EOF {
				break
			}
			// could terminate connection here but we'll allow the user to
			// continue in the case they entered something wrong
			log.Error("Read error", err)
			continue
		}
		msg := buf[:n]
		log.Info(string(msg))

		s.Broadcast(msg)
	}
}

// Broadcast messages to all clients
func (s *Server) Broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			// writes bytes from websocket connection to all websocket connections
			if _, err := ws.Write(b); err != nil {
				log.Error("Write error:", err)
			}
			log.With("client", ws.LocalAddr().String(), "broadcastMsg", string(b)).Debug("Sending message to client")
		}(ws)
	}
}
