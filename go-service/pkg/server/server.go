package server

import (
	"context"
	"fmt"
	"go-service/pkg/types"
	"go-service/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/websocket"
)

type StreamServer struct {
	clients map[string]*Client
	ctx     context.Context
	cancel  context.CancelFunc
	server  *http.Server
}

func NewStreamServer(controller types.IController) *StreamServer {
	ctx, cancel := context.WithCancel(context.Background())

	ss := &StreamServer{
		clients: make(map[string]*Client),
		ctx:     ctx,
		cancel:  cancel,
	}
	ss.server = ss.configServer(controller)

	return ss
}

func (s *StreamServer) handleWebSocket(controller types.IController) func(w http.ResponseWriter, r *http.Request) {

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := NewClient(conn, utils.GenerateUUID())

		controller.HandleDriver(client)

	}

}

func (s *StreamServer) configServer(controller types.IController) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "I'm alive!")
	})
	mux.HandleFunc("/ws", s.handleWebSocket(controller))

	cors := s.corsConfig
	handler := handlers.CORS(cors().AllowedHeaders, cors().AllowedMethods, cors().AllowedHeaders)(mux)
	server := &http.Server{
		Addr:    ":8082",
		Handler: handler,
	}
	return server
}

func (s *StreamServer) Run() {

	log.Println("Server starting on :8082")
	s.loop()

}
func (s *StreamServer) loop() {

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	s.gracefulShutdown(s.server)
}

func (s *StreamServer) gracefulShutdown(server *http.Server) {
	log.Println("Shutting down server...")

	shutdownCtx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	s.cancel()
}

type Cors struct {
	AllowedOrigins handlers.CORSOption
	AllowedMethods handlers.CORSOption
	AllowedHeaders handlers.CORSOption
}

func (s *StreamServer) corsConfig() Cors {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	cors := Cors{
		AllowedOrigins: originsOk,
		AllowedMethods: methodsOk,
		AllowedHeaders: headersOk,
	}
	return cors
}
