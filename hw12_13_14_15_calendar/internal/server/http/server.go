package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	host    string
	port    string
	logger  Logger
	app     Application
	server  *http.Server
	logFile *os.File
}

type Logger interface {
	Info(msg string, a ...any)
	Error(msg string, a ...any) // TODO
}

type Application interface {
	// HandleRequest(w http.ResponseWriter, r *http.Request)
}

// func NewServer(host string, port string, app Application, logger Logger) *Server {

func NewServer(host string, port string, app Application, logger Logger) *Server {
	return &Server{
		host:   host,
		port:   port,
		app:    app,
		logger: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// Создаём маршрутизатор
	router := mux.NewRouter()

	// Регистрация маршрута /hello
	router.HandleFunc("/hello", s.HandleRequest)

	// Настройка логирования
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to open log file: %v", err))
		return err
	}
	s.logFile = logFile
	log.SetOutput(logFile)

	// Обёртываем маршрутизатор для логирования
	handler := s.loggingMiddleware(router)

	// Настраиваем HTTP-сервер
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.host, s.port),
		Handler: handler,
	}

	// Запуск сервера в отдельной горутине
	go func() {
		<-ctx.Done()
		s.Stop(ctx)
	}()
	s.logger.Info(fmt.Sprintf("Starting server on %s:%s...", s.host, s.port))

	// Ожидание завершения через контекст
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	// Закрытие HTTP-сервера
	ctxShutdown, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s.logger.Info("Shutting down server...")
	err := s.server.Shutdown(ctxShutdown)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error during server shutdown: %v", err))
	}
	if s.logFile != nil {
		s.logFile.Close()
	}
	return err
}

func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello, World!")
}
