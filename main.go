package main

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/phllpmcphrsn/chatwithdatafeed/server"
	log "golang.org/x/exp/slog"
	"golang.org/x/net/websocket"
)

func setLogger() {
	logger := log.New(log.NewJSONHandler(os.Stdout, nil))
	log.SetDefault(logger)
}
func main() {
	// setup logger for global use
	setLogger()
	
	// setup chat server
	s := server.NewServer()
	http.Handle("/ws", websocket.Handler(s.HandleWebSocket))
	log.Info("Gopher Chat Started")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Error("Unable to start server: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<- c
		s.Broadcast([]byte("Server closed"))
		os.Exit(1)
	}()
}
